package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthRouter struct{}

func (ar *AuthRouter) InstallRouteHandlers(r *gin.RouterGroup) {
	r.GET("/current", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"user": "you"})
	})
}
