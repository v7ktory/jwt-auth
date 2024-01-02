package rest

import (
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/v7ktory/fullstack/internal/model"
)

func (h *Handler) register(c *gin.Context) {
	var regReq model.RegisterRequest

	if err := c.BindJSON(&regReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		return
	}

	user := model.User{
		Username: regReq.Username,
		Email:    regReq.Email,
		Password: regReq.Password,
	}

	if err := h.services.RegisterUser(c, user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

func (h *Handler) login(c *gin.Context) {
	var logReq model.LoginRequest

	if err := c.BindJSON(&logReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		return
	}

	token, err := h.services.AuthenticateUser(c, logReq.Email, logReq.Password)
	if err != nil {
		if errors.Is(err, model.ErrUserNotFound) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
		return
	}

	c.SetCookie("token", token, int((time.Hour * 24).Seconds()), "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}

// ProfileHandler - хендлер для получения профиля пользователя
func (h *Handler) profileHandler(c *gin.Context) {
	// Извлечение информации о пользователе из контекста Gin
	jwtPayload, ok := jwtPayloadFromRequest(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Ваша логика получения информации о пользователе
	// Замените этот блок кода своим способом получения информации о пользователе

	userInfo := gin.H{
		"name":  jwtPayload["sub"],
		"email": jwtPayload["email"],
	}

	c.JSON(http.StatusOK, userInfo)
}

// AuthRequired - middleware для проверки авторизации с использованием JWT
func authRequired() gin.HandlerFunc {

	return func(c *gin.Context) {
		token, err := c.Request.Cookie("token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		jwtToken, err := jwt.Parse(token.Value, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET_KEY")), nil
		})
		if err != nil || !jwtToken.Valid {
			c.JSON(http.StatusUnauthorized, model.ErrUnauthorized)
			c.Abort()
			return
		}

		c.Set("user", jwtToken)

		c.Next()
	}
}

func jwtPayloadFromRequest(c *gin.Context) (jwt.MapClaims, bool) {
	jwtToken, exists := c.Get("user")
	if !exists {
		return nil, false
	}

	payload, ok := jwtToken.(*jwt.Token).Claims.(jwt.MapClaims)
	if !ok {
		return nil, false
	}

	return payload, true
}
