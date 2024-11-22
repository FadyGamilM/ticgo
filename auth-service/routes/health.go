package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HealthRouter struct {
}

func (hr *HealthRouter) InstallRouteHandlers(r *gin.RouterGroup) {
	r.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"status": "pong",
		})
	})
}
