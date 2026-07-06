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
	oauthStateCookieName   = "pod_events_oauth_state"
	oauthStateCookieMaxAge = 5 * 60
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

// SpotifyLogin generates a state token, stores it in a cookie, and
// redirects the user to Spotify's authorization page.
//
//	@Summary     Initiate Spotify OAuth login
//	@Description Generates an OAuth state token and redirects to Spotify's authorization page
//	@Tags        Auth
//	@Success     307
//	@Router      /auth/spotify/login [get]
func (h *AuthHandler) SpotifyLogin(c *gin.Context) {
	state, err := h.authService.GenerateState()
	if err != nil {
		h.logger.Warn("failed to generate spotify oauth state", zap.Error(err))
		response.ErrorResponse(c, http.StatusInternalServerError, "failed to generate spotify oauth state")
		return
	}
	h.authService.RememberOAuthState(state)
	h.setOAuthStateCookie(c, state)

	authURL := h.authService.GetAuthorizationURL(state)
	c.Redirect(http.StatusTemporaryRedirect, authURL)
}

// SpotifyCallback handles the redirect back from Spotify after the user logs in.
// Spotify sends ?code=XXX&state=YYY as query parameters.
//
//	@Summary     Spotify OAuth callback
//	@Description Handles the redirect from Spotify after user authorization, exchanges code for tokens, and redirects to frontend with an exchange code
//	@Tags        Auth
//	@Param       code  query string true "Authorization code from Spotify"
//	@Param       state query string true "OAuth state token for CSRF protection"
//	@Success     307
//	@Router      /auth/spotify/callback [get]
func (h *AuthHandler) SpotifyCallback(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")
	spotifyError := c.Query("error")
	browserState, _ := c.Cookie(oauthStateCookieName)
	h.clearOAuthStateCookie(c)

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
	user, token, err := h.authService.HandleCallback(c.Request.Context(), code, state, browserState)
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

// ExchangeAuthCode exchanges a temporary auth code for a JWT token.
//
//	@Summary     Exchange auth code for JWT
//	@Description Exchanges the temporary auth code (obtained from the frontend callback redirect) for a JWT token
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

func (h *AuthHandler) setOAuthStateCookie(c *gin.Context, state string) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(oauthStateCookieName, state, oauthStateCookieMaxAge, "/", "", h.secureCookie(c), true)
}

func (h *AuthHandler) clearOAuthStateCookie(c *gin.Context) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(oauthStateCookieName, "", -1, "/", "", h.secureCookie(c), true)
}

func (h *AuthHandler) secureCookie(c *gin.Context) bool {
	return h.cfg.IsProduction() || c.Request.TLS != nil || c.GetHeader("X-Forwarded-Proto") == "https"
}
