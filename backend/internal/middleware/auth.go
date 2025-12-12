package middleware

import (
	"jp-ru-dict/backend/internal/service"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware создает middleware для проверки JWT аутентификации
func AuthMiddleware(authService service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "требуется заголовок Authorization"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "неверный формат заголовка Authorization"})
			c.Abort()
			return
		}

		tokenString := parts[1]
		token, err := authService.ValidateToken(tokenString)
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "неверный или просроченный токен"})
			c.Abort()
			return
		}

		user, err := authService.GetUserFromToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "не удалось получить данные пользователя"})
			c.Abort()
			return
		}

		// Сохраняем пользователя в контексте
		c.Set("user", user)
		c.Set("user_id", user.ID)
		c.Next()
	}
}
