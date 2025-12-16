package middleware

import (
	"net/http"
	"strings"

	"gin-real-time-talk/internal/entity/interfaces"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(authUsecase interfaces.AuthUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var token string

		accessToken, err := c.Cookie("access_token")
		if err == nil && accessToken != "" {
			token = accessToken
		} else {
			authHeader := c.GetHeader("Authorization")
			if authHeader != "" {
				parts := strings.Split(authHeader, " ")
				if len(parts) == 2 && parts[0] == "Bearer" {
					token = parts[1]
				}
			}
		}

		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "access token required"})
			c.Abort()
			return
		}

		user, err := authUsecase.ValidateAccessToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "invalid or expired token"})
			c.Abort()
			return
		}

		c.Set("user", user)
		c.Set("userID", user.ID)
		c.Next()
	}
}
