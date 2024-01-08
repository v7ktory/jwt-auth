package rest

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/v7ktory/fullstack/internal/model"
	"github.com/v7ktory/fullstack/pkg/logger"
	"go.uber.org/zap"
)

var log *zap.Logger

func init() {
	log = logger.NewLogger()
}
func (h *Handler) register(c *gin.Context) {

	var regReq model.RegisterRequest
	if err := c.BindJSON(&regReq); err != nil {
		log.Info("Failed to bind JSON for registration request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		return
	}

	user := model.User{
		Username: regReq.Username,
		Email:    regReq.Email,
		Password: regReq.Password,
	}

	token, err := h.services.SignUp(c, user)
	if err != nil {
		handleRegistrationError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"token": token, "message": "User registered successfully"})
}

func handleRegistrationError(c *gin.Context, err error) {

	switch {
	case errors.Is(err, model.ErrUserAlreadyExists):
		log.Info("User registration failed: user already exists", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": "User already exists"})
	default:
		log.Info("Failed to register user", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
	}
}

func (h *Handler) login(c *gin.Context) {

	var logReq model.LoginRequest
	if err := c.BindJSON(&logReq); err != nil {
		log.Info("Failed to bind JSON for login request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		return
	}

	token, err := h.services.SignIn(c, logReq.Email, logReq.Password)
	if err != nil {
		log.Info("Failed to sign in", zap.Error(err))
		handleLoginError(c, err)
		return
	}

	setTokenCookie(c, token)
	c.JSON(http.StatusOK, gin.H{"token": token, "message": "Login successful"})
}

func handleLoginError(c *gin.Context, err error) {

	switch {
	case errors.Is(err, model.ErrUserNotFound):
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
	default:
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
	}
}

func setTokenCookie(c *gin.Context, token string) {

	cookie := http.Cookie{
		Name:     "token",
		Value:    token,
		MaxAge:   int((time.Hour * 12).Seconds()),
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(c.Writer, &cookie)
}
