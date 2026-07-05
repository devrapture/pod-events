package handler

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/devrapture/pod-events/internal/config"
	"github.com/devrapture/pod-events/internal/repositories"
	"github.com/devrapture/pod-events/internal/services"
	"github.com/devrapture/pod-events/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const stateCookieName = "spotify_oauth_state"

type AuthHandler struct {
	authService services.AuthService
	logger      *zap.Logger
	cfg         *config.Config
	userRepo    repositories.UserRepository
}

func NewAuthHandler(authService services.AuthService, logger *zap.Logger, cfg *config.Config, userRepo repositories.UserRepository) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		logger:      logger,
		cfg:         cfg,
		userRepo:    userRepo,
	}
}

// Me returns the authenticated user's profile.
//
// GET /auth/me
func (h *AuthHandler) Me(c *gin.Context) {
	userIDStr, exists := c.Get("userID")
	if !exists {
		response.ErrorResponse(c, http.StatusUnauthorized, "user not authenticated")
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		h.logger.Warn("invalid userID in context", zap.Error(err))
		response.ErrorResponse(c, http.StatusInternalServerError, "invalid user ID")
		return
	}

	user, err := h.userRepo.GetByID(c.Request.Context(), userID)
	if err != nil {
		h.logger.Error("failed to fetch user", zap.Error(err))
		response.ErrorResponse(c, http.StatusInternalServerError, "failed to fetch user")
		return
	}
	if user == nil {
		response.ErrorResponse(c, http.StatusNotFound, "user not found")
		return
	}

	response.SuccessResponse(c, http.StatusOK, "User fetched successfully", user, nil)
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
	spotifyError := c.Query("error")
	if spotifyError != "" {
		redirectURL := fmt.Sprintf("%s/auth/callback?error=%s", h.cfg.FrontendURL, url.QueryEscape("Spotify login rejected"))
		c.Redirect(http.StatusTemporaryRedirect, redirectURL)
		return
	}

	if code == "" || state == "" {
		redirectURL := fmt.Sprintf("%s/auth/callback?error=%s", h.cfg.FrontendURL, url.QueryEscape("Missing code or state parameter"))
		c.Redirect(http.StatusTemporaryRedirect, redirectURL)
		return
	}
	cookieState, err := c.Cookie(stateCookieName)
	if err != nil {
		h.logger.Warn("failed to get state cookie", zap.Error(err))
		redirectURL := fmt.Sprintf("%s/auth/callback?error=%s", h.cfg.FrontendURL, url.QueryEscape("Invalid or expired OAuth state"))
		c.Redirect(http.StatusTemporaryRedirect, redirectURL)
		return
	}
	secureCookie := h.cfg.IsProduction()
	c.SetCookie(
		stateCookieName,
		"",
		-1,
		"/",
		"",
		secureCookie,
		true,
	)

	_, token, err := h.authService.HandleCallback(c.Request.Context(), code, state, cookieState)
	if err != nil {
		h.logger.Warn("failed to authenticate with spotify", zap.Error(err))
		redirectURL := fmt.Sprintf("%s/auth/callback?error=%s", h.cfg.FrontendURL, url.QueryEscape("Authentication failed"))
		c.Redirect(http.StatusTemporaryRedirect, redirectURL)
		return
	}

	redirectURL := fmt.Sprintf("%s/auth/callback?token=%s", h.cfg.FrontendURL, token)
	c.Redirect(http.StatusTemporaryRedirect, redirectURL)
}
