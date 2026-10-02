package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/devrapture/pod-events/internal/cron"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

type stubEpisodeRunner struct {
	run func(context.Context) (*cron.CheckResult, error)
}

func (s stubEpisodeRunner) Run(ctx context.Context) (*cron.CheckResult, error) {
	return s.run(ctx)
}

func TestCronJobHandlerWaitsForRunAndReturnsResult(t *testing.T) {
	gin.SetMode(gin.TestMode)

	started := make(chan struct{})
	release := make(chan struct{})
	runner := stubEpisodeRunner{run: func(context.Context) (*cron.CheckResult, error) {
		close(started)
		<-release
		return &cron.CheckResult{
			ShowsChecked:      3,
			NewEpisodes:       1,
			NotificationsSent: 2,
			Errors:            []string{},
		}, nil
	}}
	handler := &CronJobHandler{
		logger:         zap.NewNop(),
		episodeChecker: runner,
		timeout:        cronJobTimeout,
	}

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/cron/check-episodes", nil)
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = request
	done := make(chan struct{})

	go func() {
		handler.CheckEpisodes(ctx)
		close(done)
	}()

	<-started
	select {
	case <-done:
		t.Fatal("handler returned before the episode check completed")
	case <-time.After(25 * time.Millisecond):
	}

	close(release)
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("handler did not return after the episode check completed")
	}

	require.Equal(t, http.StatusOK, recorder.Code)
	var body struct {
		Success bool             `json:"success"`
		Message string           `json:"message"`
		Data    cron.CheckResult `json:"data"`
	}
	require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &body))
	assert.True(t, body.Success)
	assert.Equal(t, "Cron job completed", body.Message)
	assert.Equal(t, 3, body.Data.ShowsChecked)
	assert.Equal(t, 1, body.Data.NewEpisodes)
	assert.Equal(t, 2, body.Data.NotificationsSent)
}

func TestCronJobHandlerReturnsFailure(t *testing.T) {
	gin.SetMode(gin.TestMode)

	runner := stubEpisodeRunner{run: func(context.Context) (*cron.CheckResult, error) {
		return nil, errors.New("database unavailable")
	}}
	handler := &CronJobHandler{
		logger:         zap.NewNop(),
		episodeChecker: runner,
		timeout:        cronJobTimeout,
	}

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/cron/check-episodes", nil)

	handler.CheckEpisodes(ctx)

	require.Equal(t, http.StatusInternalServerError, recorder.Code)
	assert.JSONEq(t, `{"success":false,"error":{"message":"Cron job failed"}}`, recorder.Body.String())
}

func TestCronJobHandlerReturnsBeforeSchedulerDeadline(t *testing.T) {
	gin.SetMode(gin.TestMode)

	runner := stubEpisodeRunner{run: func(ctx context.Context) (*cron.CheckResult, error) {
		<-ctx.Done()
		return nil, ctx.Err()
	}}
	handler := &CronJobHandler{
		logger:         zap.NewNop(),
		episodeChecker: runner,
		timeout:        10 * time.Millisecond,
	}

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/cron/check-episodes", nil)

	handler.CheckEpisodes(ctx)

	require.Equal(t, http.StatusGatewayTimeout, recorder.Code)
	assert.JSONEq(t, `{"success":false,"error":{"message":"Cron job timed out"}}`, recorder.Body.String())
}
