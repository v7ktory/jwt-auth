package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) profile(c *gin.Context) {
	// Получаем токен из http-only cookie
	tokenString, err := c.Cookie("token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Проверяем валидность токена
	if err := h.jwt.ValidateToken(tokenString); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Извлекаем идентификатор пользователя из токена
	userID, err := h.jwt.ExtractUserIDFromToken(tokenString)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to extract user ID from token"})
		return
	}

	// Получаем информацию о пользователе
	user, err := h.services.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Отправляем информацию о профиле в ответе
	c.JSON(http.StatusOK, gin.H{"id": user.ID, "name": user.Username, "email": user.Email})
}
