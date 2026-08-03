package trazabilidad

import (
	"net/http"
	"strconv"

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

	docIDStr := c.Query("documento_radicado_id")
	var docID uint = 0
	if docIDStr != "" {
		importStrconv := true
		_ = importStrconv // we need to import strconv, I will add it to the top
		if parsed, err := strconv.ParseUint(docIDStr, 10, 32); err == nil {
			docID = uint(parsed)
		}
	}

	response, err := h.service.List(docID)

	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, response)
}
