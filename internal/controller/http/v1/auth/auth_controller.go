package auth

import (
	"net/http"
	"time"

	"gin-real-time-talk/config"
	"gin-real-time-talk/internal/entity/interfaces"

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

func (ac *AuthController) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := ac.authUsecase.Register(req.Email, req.Password, req.FirstName, req.LastName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "user registered successfully, verification code sent to email",
		"user":    user,
	})
}

func (ac *AuthController) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accessToken, refreshToken, user, err := ac.authUsecase.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	setTokenCookies(c, accessToken, refreshToken)

	c.JSON(http.StatusOK, gin.H{
		"message": "login successful",
		"user":    user,
	})
}

func (ac *AuthController) VerifyCode(c *gin.Context) {
	var req VerifyCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accessToken, refreshToken, user, err := ac.authUsecase.VerifyTwoFactorCode(req.Email, req.Code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	setTokenCookies(c, accessToken, refreshToken)

	c.JSON(http.StatusOK, gin.H{
		"message": "email verified successfully",
		"user":    user,
	})
}

func (ac *AuthController) ResendCode(c *gin.Context) {
	var req ResendCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := ac.authUsecase.SendTwoFactorCode(req.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "verification code sent to email",
	})
}

func (ac *AuthController) Refresh(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "refresh token not found"})
		return
	}

	accessToken, newRefreshToken, user, err := ac.authUsecase.RefreshToken(refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	setTokenCookies(c, accessToken, newRefreshToken)

	c.JSON(http.StatusOK, gin.H{
		"message": "tokens refreshed successfully",
		"user":    user,
	})
}

func (ac *AuthController) Me(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}
