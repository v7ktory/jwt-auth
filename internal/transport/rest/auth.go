package rest

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
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
