package auth

import (
	"net/http"
	"time"

	"gin-real-time-talk/config"
	"gin-real-time-talk/internal/entity/interfaces"
	"gin-real-time-talk/pkg/validator"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authUsecase interfaces.AuthUsecase
}

func NewAuthController(authUsecase interfaces.AuthUsecase) *AuthController {
	return &AuthController{
		authUsecase: authUsecase,
	}
}

type RegisterRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type VerifyCodeRequest struct {
	Email string `json:"email" binding:"required,email"`
	Code  string `json:"code" binding:"required"`
}

type ResendCodeRequest struct {
	Email string `json:"email" binding:"required,email"`
}

func setTokenCookies(c *gin.Context, accessToken, refreshToken string) {
	accessExpiry, _ := time.ParseDuration(config.Env.JWT.AccessExpiry)
	if accessExpiry == 0 {
		accessExpiry = 15 * time.Minute
	}

	refreshExpiry, _ := time.ParseDuration(config.Env.JWT.RefreshExpiry)
	if refreshExpiry == 0 {
		refreshExpiry = 7 * 24 * time.Hour
	}

	isSecure := config.Env.App.Environment == "production"

	c.SetCookie("access_token", accessToken, int(accessExpiry.Seconds()), "/", "", isSecure, true)
	c.SetCookie("refresh_token", refreshToken, int(refreshExpiry.Seconds()), "/", "", isSecure, true)
}

// Register godoc
// @Summary Register new user
// @Description Registers a new user and sends verification code to email
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Registration data"
// @Success 201 {object} map[string]interface{} "User successfully registered"
// @Failure 400 {object} map[string]string "Validation or registration error"
// @Router /auth/register [post]
func (ac *AuthController) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": validator.FormatErrors(err)})
		return
	}

	user, err := ac.authUsecase.Register(req.Email, req.Password, req.FirstName, req.LastName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"user":    user,
	})
}

// Login godoc
// @Summary User login
// @Description Authenticates user and returns access tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login credentials"
// @Success 200 {object} map[string]interface{} "Successful login"
// @Failure 400 {object} map[string]string "Validation error"
// @Failure 401 {object} map[string]string "Invalid credentials"
// @Router /auth/login [post]
func (ac *AuthController) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": validator.FormatErrors(err)})
		return
	}

	accessToken, refreshToken, user, err := ac.authUsecase.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": err.Error()})
		return
	}

	setTokenCookies(c, accessToken, refreshToken)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"user":    user,
	})
}

// VerifyCode godoc
// @Summary Verify code
// @Description Verifies email verification code and issues access tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param request body VerifyCodeRequest true "Email and verification code"
// @Success 200 {object} map[string]interface{} "Email successfully verified"
// @Failure 400 {object} map[string]string "Validation error or invalid code"
// @Router /auth/verify [post]
func (ac *AuthController) VerifyCode(c *gin.Context) {
	var req VerifyCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": validator.FormatErrors(err)})
		return
	}

	accessToken, refreshToken, user, err := ac.authUsecase.VerifyTwoFactorCode(req.Email, req.Code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	setTokenCookies(c, accessToken, refreshToken)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"user":    user,
	})
}

// ResendCode godoc
// @Summary Resend verification code
// @Description Sends a new verification code to the specified email
// @Tags auth
// @Accept json
// @Produce json
// @Param request body ResendCodeRequest true "Email to send code to"
// @Success 200 {object} map[string]string "Code sent to email"
// @Failure 400 {object} map[string]string "Validation or sending error"
// @Router /auth/resend-code [post]
func (ac *AuthController) ResendCode(c *gin.Context) {
	var req ResendCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": validator.FormatErrors(err)})
		return
	}

	err := ac.authUsecase.SendTwoFactorCode(req.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}

// Refresh godoc
// @Summary Refresh access tokens
// @Description Refreshes access and refresh tokens using current refresh token from cookies
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Tokens successfully refreshed"
// @Failure 401 {object} map[string]string "Refresh token not found or invalid"
// @Router /auth/refresh [post]
func (ac *AuthController) Refresh(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "refresh token not found"})
		return
	}

	accessToken, newRefreshToken, user, err := ac.authUsecase.RefreshToken(refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": err.Error()})
		return
	}

	setTokenCookies(c, accessToken, newRefreshToken)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"user":    user,
	})
}

// Me godoc
// @Summary Get current user
// @Description Returns information about the current authenticated user
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "User information"
// @Failure 401 {object} map[string]string "User not authenticated"
// @Router /auth/me [get]
func (ac *AuthController) Me(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"user":    user,
	})
}
