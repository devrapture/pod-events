package handler

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/devrapture/pod-events/internal/config"
	"github.com/devrapture/pod-events/internal/dto"
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

type authExchangeRequest dto.AuthExchangeRequest

const (
	googleOAuthStateCookieName = "pod_events_google_oauth_state"
	oauthStateCookieMaxAge     = 5 * 60
)

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
//	@Summary     Get current user
//	@Description Get the authenticated user's profile
//	@Tags        Auth
//	@Security    BearerAuth
//	@Success     200 {object} response.APIResponse "User fetched successfully"
//	@Failure     401 {object} response.APIResponse "unauthorized"
//	@Router      /auth/me [get]
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

// GoogleLogin starts the public sign-in flow.
//
//	@Summary     Initiate Google OAuth login
//	@Description Generates an OAuth state token and redirects to Google's authorization page
//	@Tags        Auth
//	@Success     307
//	@Router      /auth/google/login [get]
func (h *AuthHandler) GoogleLogin(c *gin.Context) {
	state, err := h.authService.GenerateState()
	if err != nil {
		h.logger.Warn("failed to generate google oauth state", zap.Error(err))
		response.ErrorResponse(c, http.StatusInternalServerError, "failed to start google sign-in")
		return
	}
	h.authService.RememberGoogleOAuthState(state)
	h.setOAuthStateCookie(c, googleOAuthStateCookieName, state)
	c.Redirect(http.StatusTemporaryRedirect, h.authService.GetGoogleAuthorizationURL(state))
}

// GoogleCallback completes public sign-in and redirects the browser to the SPA
// with a short-lived, one-time exchange code.
//
//	@Summary     Google OAuth callback
//	@Description Exchanges the Google code, links or creates a user, and redirects to the frontend
//	@Tags        Auth
//	@Param       code  query string true "Authorization code from Google"
//	@Param       state query string true "OAuth state token for CSRF protection"
//	@Success     307
//	@Router      /auth/google/callback [get]
func (h *AuthHandler) GoogleCallback(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")
	providerError := c.Query("error")
	browserState, _ := c.Cookie(googleOAuthStateCookieName)
	h.clearOAuthStateCookie(c, googleOAuthStateCookieName)

	if providerError != "" {
		h.redirectAuthError(c, "Google sign-in rejected")
		return
	}
	if code == "" || state == "" {
		h.redirectAuthError(c, "Missing code or state parameter")
		return
	}

	user, token, err := h.authService.HandleGoogleCallback(c.Request.Context(), code, state, browserState)
	if err != nil {
		h.logger.Warn("failed to authenticate with google", zap.Error(err))
		h.redirectAuthError(c, "Authentication failed")
		return
	}
	exchangeCode, err := h.authService.CreateAuthExchangeCode(token, user)
	if err != nil {
		h.logger.Error("failed to create auth exchange code", zap.Error(err))
		h.redirectAuthError(c, "Authentication failed")
		return
	}

	redirectURL := fmt.Sprintf("%s/auth/callback?code=%s", h.cfg.FrontendURL, url.QueryEscape(exchangeCode))
	c.Redirect(http.StatusTemporaryRedirect, redirectURL)
}

// SpotifyOwnerLogin creates an owner-bound Spotify authorization URL. This
// endpoint is intentionally absent from the public UI and restricted by email.
//
//	@Summary     Initiate owner Spotify authorization
//	@Description Returns a Spotify authorization URL for the configured application owner
//	@Tags        Auth
//	@Security    BearerAuth
//	@Success     200 {object} response.APIResponse
//	@Failure     403 {object} response.APIResponse
//	@Router      /auth/spotify/owner/login [get]
func (h *AuthHandler) SpotifyOwnerLogin(c *gin.Context) {
	userID, _ := c.Get("userID")
	isOwner, err := h.authService.IsOwner(c.Request.Context(), userID.(uuid.UUID))
	if err != nil {
		h.logger.Error("failed to verify application owner", zap.Error(err))
		response.ErrorResponse(c, http.StatusInternalServerError, "failed to verify application owner")
		return
	}
	if !isOwner {
		response.ErrorResponse(c, http.StatusForbidden, "only the application owner can connect Spotify")
		return
	}

	state, err := h.authService.GenerateState()
	if err != nil {
		h.logger.Warn("failed to generate spotify owner oauth state", zap.Error(err))
		response.ErrorResponse(c, http.StatusInternalServerError, "failed to start spotify authorization")
		return
	}
	h.authService.RememberSpotifyOwnerOAuthState(state, userID.(uuid.UUID))
	response.SuccessResponse(c, http.StatusOK, "Spotify authorization ready", gin.H{
		"authorization_url": h.authService.GetSpotifyOwnerAuthorizationURL(state),
	}, nil)
}

// SpotifyOwnerCallback stores the owner's encrypted, refreshable Spotify token.
//
//	@Summary     Owner Spotify OAuth callback
//	@Description Completes the owner-only Spotify authorization flow
//	@Tags        Auth
//	@Param       code  query string true "Authorization code from Spotify"
//	@Param       state query string true "OAuth state token for CSRF protection"
//	@Success     307
//	@Router      /auth/spotify/callback [get]
func (h *AuthHandler) SpotifyOwnerCallback(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")
	providerError := c.Query("error")

	if providerError != "" {
		h.redirectSpotifyOwner(c, "error", "Spotify authorization rejected")
		return
	}
	if code == "" || state == "" {
		h.redirectSpotifyOwner(c, "error", "Missing code or state parameter")
		return
	}
	if err := h.authService.HandleSpotifyOwnerCallback(c.Request.Context(), code, state); err != nil {
		h.logger.Warn("failed to authorize spotify owner", zap.Error(err))
		h.redirectSpotifyOwner(c, "error", "Spotify authorization failed")
		return
	}
	h.redirectSpotifyOwner(c, "connected", "1")
}

// ExchangeAuthCode exchanges a temporary auth code for a JWT token.
//
//	@Summary     Exchange auth code for JWT
//	@Description Exchanges the temporary code obtained from the frontend callback redirect
//	@Tags        Auth
//	@Accept      json
//	@Produce     json
//	@Param       request body dto.AuthExchangeRequest true "Auth exchange request"
//	@Success     200 {object} response.APIResponse "Authentication completed"
//	@Failure     400 {object} response.APIResponse "missing auth exchange code"
//	@Failure     401 {object} response.APIResponse "invalid or expired auth exchange code"
//	@Router      /auth/exchange [post]
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

func (h *AuthHandler) redirectAuthError(c *gin.Context, message string) {
	redirectURL := fmt.Sprintf("%s/auth/callback?error=%s", h.cfg.FrontendURL, url.QueryEscape(message))
	c.Redirect(http.StatusTemporaryRedirect, redirectURL)
}

func (h *AuthHandler) redirectSpotifyOwner(c *gin.Context, key, value string) {
	redirectURL := fmt.Sprintf("%s/owner/spotify?%s=%s", h.cfg.FrontendURL, url.QueryEscape(key), url.QueryEscape(value))
	c.Redirect(http.StatusTemporaryRedirect, redirectURL)
}

func (h *AuthHandler) setOAuthStateCookie(c *gin.Context, name, state string) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(name, state, oauthStateCookieMaxAge, "/", "", h.secureCookie(c), true)
}

func (h *AuthHandler) clearOAuthStateCookie(c *gin.Context, name string) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(name, "", -1, "/", "", h.secureCookie(c), true)
}

func (h *AuthHandler) secureCookie(c *gin.Context) bool {
	return h.cfg.IsProduction() || c.Request.TLS != nil || c.GetHeader("X-Forwarded-Proto") == "https"
}
