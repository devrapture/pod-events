package spotify

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	apperrors "github.com/devrapture/pod-events/internal/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func newTestClient(t *testing.T, handler http.HandlerFunc) *SpotifyClient {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	base, err := url.Parse(server.URL)
	require.NoError(t, err)

	return &SpotifyClient{
		logger: zap.NewNop(),
		httpClient: &http.Client{
			Transport: rewriteHostRoundTripper{base: base, next: http.DefaultTransport},
		},
	}
}

// rewriteHostRoundTripper sends api.spotify.com requests to the httptest server.
type rewriteHostRoundTripper struct {
	base *url.URL
	next http.RoundTripper
}

func (r rewriteHostRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if strings.Contains(req.URL.Host, "api.spotify.com") {
		cloned := req.Clone(req.Context())
		u := *req.URL
		u.Scheme = r.base.Scheme
		u.Host = r.base.Host
		cloned.URL = &u
		cloned.Host = r.base.Host
		req = cloned
	}
	return r.next.RoundTrip(req)
}

func TestGetShowLatestEpisode_SkipsLeadingNulls(t *testing.T) {
	t.Parallel()

	var gotPath string
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.String()
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{
			"items": [
				null,
				null,
				{
					"id": "ep-real",
					"name": "Real Episode",
					"description": "desc",
					"external_urls": {"spotify": "https://open.spotify.com/episode/ep-real"},
					"images": [],
					"duration_ms": 1000,
					"release_date": "2026-08-01",
					"release_date_precision": "day"
				}
			],
			"total": 10
		}`)
	})

	ep, err := client.GetShowLatestEpisode(context.Background(), "token", "show-1")
	require.NoError(t, err)
	require.NotNil(t, ep)
	assert.Equal(t, "ep-real", ep.ID)
	assert.Equal(t, "Real Episode", ep.Name)
	assert.Contains(t, gotPath, "limit=5")
	assert.Contains(t, gotPath, "/shows/show-1/episodes")
}

func TestGetShowLatestEpisode_AllNulls(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"items":[null,null],"total":2}`)
	})

	ep, err := client.GetShowLatestEpisode(context.Background(), "token", "show-1")
	assert.Nil(t, ep)
	assert.ErrorIs(t, err, apperrors.ErrSpotifyResourceNotFound)
}

func TestGetShowLatestEpisode_EmptyItems(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"items":[],"total":0}`)
	})

	ep, err := client.GetShowLatestEpisode(context.Background(), "token", "show-1")
	assert.Nil(t, ep)
	assert.ErrorIs(t, err, apperrors.ErrSpotifyResourceNotFound)
}

func TestGetShowLatestEpisode_FirstItemReal(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"items": []any{
				map[string]any{
					"id":                     "ep-0",
					"name":                   "Latest",
					"description":            "",
					"external_urls":          map[string]string{"spotify": "https://open.spotify.com/episode/ep-0"},
					"images":                 []any{},
					"duration_ms":            1,
					"release_date":           "2026-01-01",
					"release_date_precision": "day",
				},
			},
			"total": 1,
		})
	})

	ep, err := client.GetShowLatestEpisode(context.Background(), "token", "show-1")
	require.NoError(t, err)
	require.NotNil(t, ep)
	assert.Equal(t, "ep-0", ep.ID)
}

func TestGetShowLatestEpisode_SkipsEmptyID(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{
			"items": [
				{"id":"","name":"empty"},
				{"id":"ep-ok","name":"OK","description":"","external_urls":{"spotify":"x"},"images":[],"duration_ms":1,"release_date":"2026-01-02","release_date_precision":"day"}
			],
			"total": 2
		}`)
	})

	ep, err := client.GetShowLatestEpisode(context.Background(), "token", "show-1")
	require.NoError(t, err)
	require.NotNil(t, ep)
	assert.Equal(t, "ep-ok", ep.ID)
}

func TestGetShowLatestEpisode_NotFoundStatus(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = io.WriteString(w, `{"error":{"status":404,"message":"Not found"}}`)
	})

	ep, err := client.GetShowLatestEpisode(context.Background(), "token", "missing")
	assert.Nil(t, ep)
	assert.True(t, errors.Is(err, apperrors.ErrSpotifyResourceNotFound))
}
