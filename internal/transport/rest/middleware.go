package rest

import "github.com/gin-gonic/gin"

func (h *Handler) Auth() gin.HandlerFunc {

	return func(context *gin.Context) {
		tokenString := context.GetHeader("Authorization")
		if tokenString == "" {
			context.AbortWithStatusJSON(401, gin.H{"error": "request does not contain an access token"})
			return
		}
		if err := h.jwt.ValidateToken(tokenString); err != nil {
			context.AbortWithStatusJSON(401, gin.H{"error": err.Error()})
			return
		}
		context.Next()
	}
}
