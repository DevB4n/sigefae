package historial_asignacion

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {

	return &Handler{
		service: service,
	}
}
func (h *Handler) Create(c *gin.Context) {

	var request CreateDTO


	if err := c.ShouldBindJSON(&request); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}


	response, err := h.service.Create(request)


	if err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}


	c.JSON(http.StatusCreated, response)
}
func (h *Handler) List(c *gin.Context) {

	response, err := h.service.List()


	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}


	c.JSON(http.StatusOK, response)
}