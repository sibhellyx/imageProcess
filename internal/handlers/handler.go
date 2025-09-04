package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/sibhellyx/imageProccesor/internal/service"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Image(c *gin.Context) {

}
