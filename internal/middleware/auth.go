package middleware

import (
	"movie-ticket-booking/internal/service"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func JWTAuthMiddleware(authService *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			tokenString, _ = c.Cookie("accessToken")
		}
		if tokenString == "" {
			tokenString = c.Query("accessToken")
		}
		// Kiểm tra xem token có bắt đầu bằng "Bearer "
		if len(tokenString) > 7 && strings.ToUpper(tokenString[0:7]) == "BEARER " {
			tokenString = tokenString[7:]
		}
		token, err := authService.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		claims, _ := token.Claims.(jwt.MapClaims)
		c.Set("email", claims["email"])
		c.Set("role", claims["role"])

		c.Next()
	}
}
