package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) profile(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{"message": "Profile"})
}
