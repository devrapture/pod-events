package handler

import (
	"net/http"

	"github.com/devrapture/pod-events/internal/config"
	"github.com/devrapture/pod-events/internal/services"
	"github.com/devrapture/pod-events/pkg/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const stateCookieName = "spotify_oauth_state"

type AuthHandler struct {
	authService services.AuthService
	logger      *zap.Logger
	cfg         *config.Config
}

func NewAuthHandler(authService services.AuthService, logger *zap.Logger, cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		logger:      logger,
		cfg:         cfg,
	}
}

// SpotifyLogin generates a state token, stores it in a cookie, and
// redirects the user to Spotify's authorization page.
//
// GET /auth/spotify/login
func (h *AuthHandler) SpotifyLogin(c *gin.Context) {
	state, err := h.authService.GenerateState()
	if err != nil {
		h.logger.Warn("failed to generate spotify oauth state", zap.Error(err))
		response.ErrorResponse(c, http.StatusInternalServerError, "failed to generate spotify oauth state")
		return
	}
	secureCookie := h.cfg.IsProduction()

	c.SetCookie(
		stateCookieName,
		state,
		300, // 5 minutes
		"/",
		"",
		secureCookie,
		true, // httpOnly
	)

	authURL := h.authService.GetAuthorizationURL(state)
	c.Redirect(http.StatusTemporaryRedirect, authURL)
}

// SpotifyCallback handles the redirect back from Spotify after the user logs in.
// Spotify sends ?code=XXX&state=YYY as query parameters.
//
// GET /auth/spotify/callback?code=XXX&state=YYY
func (h *AuthHandler) SpotifyCallback(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")
	spotifyError := c.Query("error") // Spotify sends "error=access_denied" if user rejects
	if spotifyError != "" {
		response.ErrorResponse(c, http.StatusUnauthorized, "Spotify login rejected")
		return
	}

	if code == "" || state == "" {
		response.ErrorResponse(c, http.StatusBadRequest, "missing code or state parameter")
		return
	}
	cookieState, err := c.Cookie(stateCookieName)
	if err != nil {
		h.logger.Warn("failed to get state cookie", zap.Error(err))
		response.ErrorResponse(c, http.StatusUnauthorized, "invalid or expired oauth state")
		return
	}
	secureCookie := h.cfg.IsProduction()
	// Delete the state cookie (it's been used)
	c.SetCookie(
		stateCookieName,
		"",
		-1,
		"/",
		"",
		secureCookie,
		true, // httpOnly
	)

	user, token, err := h.authService.HandleCallback(c.Request.Context(), code, state, cookieState)
	if err != nil {
		h.logger.Warn("failed to authenticate with spotify", zap.Error(err))
		response.ErrorResponse(c, http.StatusInternalServerError, "authentication failed")
		return
	}

	response.SuccessResponse(c, http.StatusOK, "Successfully authenticated", gin.H{
		"user":  user,
		"token": token,
	}, nil)
}
