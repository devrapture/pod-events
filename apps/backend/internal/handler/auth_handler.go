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

type AuthHandler struct {
	authService services.AuthService
	logger      *zap.Logger
	cfg         *config.Config
	userRepo    repositories.UserRepository
}

type authExchangeRequest struct {
	Code string `json:"code" binding:"required"`
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
	userID, _ := c.Get("userID")

	user, err := h.userRepo.GetByID(c.Request.Context(), userID.(uuid.UUID))
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
	h.authService.RememberOAuthState(state)

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
	user, token, err := h.authService.HandleCallback(c.Request.Context(), code, state)
	if err != nil {
		h.logger.Warn("failed to authenticate with spotify", zap.Error(err))
		redirectURL := fmt.Sprintf("%s/auth/callback?error=%s", h.cfg.FrontendURL, url.QueryEscape("Authentication failed"))
		c.Redirect(http.StatusTemporaryRedirect, redirectURL)
		return
	}

	exchangeCode, err := h.authService.CreateAuthExchangeCode(token, user)
	if err != nil {
		h.logger.Error("failed to create auth exchange code", zap.Error(err))
		redirectURL := fmt.Sprintf("%s/auth/callback?error=%s", h.cfg.FrontendURL, url.QueryEscape("Authentication failed"))
		c.Redirect(http.StatusTemporaryRedirect, redirectURL)
		return
	}

	redirectURL := fmt.Sprintf("%s/auth/callback?code=%s", h.cfg.FrontendURL, url.QueryEscape(exchangeCode))
	c.Redirect(http.StatusTemporaryRedirect, redirectURL)
}

func (h *AuthHandler) ExchangeAuthCode(c *gin.Context) {
	var req authExchangeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "missing auth exchange code")
		return
	}

	exchange, ok := h.authService.ConsumeAuthExchangeCode(req.Code)
	if !ok {
		response.ErrorResponse(c, http.StatusUnauthorized, "invalid or expired auth exchange code")
		return
	}

	response.SuccessResponse(c, http.StatusOK, "Authentication completed", gin.H{
		"token": exchange.Token,
		"user":  exchange.User,
	}, nil)
}
