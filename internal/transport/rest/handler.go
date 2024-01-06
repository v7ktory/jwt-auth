package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/v7ktory/fullstack/internal/service"
	"github.com/v7ktory/fullstack/pkg/token"
)

type Handler struct {
	services *service.Service
	jwt      *token.JWTService
}

func NewHandler(services *service.Service, jwt *token.JWTService) *Handler {
	return &Handler{
		services: services,
		jwt:      jwt,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	r := gin.New()

	auth := r.Group("/auth")
	{
		auth.POST("/register", h.register)
		auth.POST("/login", h.login)
	}

	secured := auth.Group("/secured").Use(h.Auth())
	{
		secured.GET("/profile/:id", h.profile)
	}

	return r
}
