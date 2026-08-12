package proveedor

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

	var request CreateRequest

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

func (h *Handler) ListNormasReparto(c *gin.Context) {
	proveedorID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id de proveedor inválido"})
		return
	}

	rutaID, err := strconv.ParseUint(c.Query("ruta_id"), 10, 64)
	if err != nil || rutaID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ruta_id inválido"})
		return
	}

	normas, err := h.service.ListNormasReparto(uint(proveedorID), uint(rutaID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, normas)
}

func (h *Handler) Update(c *gin.Context) {

	id, err := strconv.ParseUint(
		c.Param("id"),
		10,
		64,
	)

	if err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "id inválido",
		})

		return
	}

	var request UpdateRequest

	if err := c.ShouldBindJSON(&request); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	response, err := h.service.Update(
		uint(id),
		request,
	)

	if err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handler) UpdateStatus(c *gin.Context) {

	id, err := strconv.ParseUint(
		c.Param("id"),
		10,
		64,
	)

	if err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "id inválido",
		})

		return
	}

	var request UpdateStatusRequest

	if err := c.ShouldBindJSON(&request); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	if err := h.service.UpdateStatus(
		uint(id),
		request.Activo,
	); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "estado actualizado correctamente",
	})
}
