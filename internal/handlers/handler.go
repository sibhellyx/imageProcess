package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sibhellyx/imageProccesor/internal/models"
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
	var request models.ImageRequest
	err := json.NewDecoder(c.Request.Body).Decode(&request)
	if err != nil {
		WrapError(c, err)
		return
	}

	task, err := h.service.AddImageTask(request)
	if err != nil {
		WrapError(c, err)
		return
	}

	err = h.service.Proccess(c.Request.Context(), task)

	if err != nil {
		WrapError(c, err)
		return
	}

	c.JSON(
		http.StatusOK, gin.H{
			"result": "created",
		},
	)
}

func WrapError(c *gin.Context, err error) {
	log.Println(err)
	c.JSON(http.StatusBadRequest, gin.H{
		"error": err.Error(),
	})
}
