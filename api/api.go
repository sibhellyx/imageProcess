package api

import (
	"github.com/gin-gonic/gin"
	"github.com/sibhellyx/imageProccesor/internal/handlers"
)

func CreateRoutes(handler *handlers.Handler) *gin.Engine {
	r := gin.Default()

	r.POST("/image", handler.Image)
	r.POST("/download", handler.Download)

	return r
}
