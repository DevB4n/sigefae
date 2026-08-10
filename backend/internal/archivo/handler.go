package archivo

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"sigefae/internal/db"
)

type Handler struct {
	service *Service
	db      *gorm.DB
}

func NewHandler(service *Service, database *gorm.DB) *Handler {

	return &Handler{
		service: service,
		db:      database,
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

	var request UpdateDTO

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

func (h *Handler) Delete(c *gin.Context) {

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

	user, hasUser := c.Get("user")
	var userID uint
	if hasUser {
		userID = user.(db.Usuario).ID
	}

	// Obtener info del archivo antes de borrarlo para trazabilidad
	var archivo db.Archivo
	if h.db != nil {
		h.db.First(&archivo, uint(id))
	}

	if err := h.service.Delete(uint(id)); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "archivo eliminado correctamente",
	})

	if h.db != nil && archivo.DocumentoRadicadoID != 0 && userID != 0 {
		h.db.Create(&db.Trazabilidad{
			DocumentoRadicadoID: archivo.DocumentoRadicadoID,
			UsuarioID:           userID,
			Accion:              "Archivo Eliminado",
			Descripcion:         "Se eliminó el archivo: " + archivo.Nombre,
			Fecha:               time.Now(),
		})
	}
}
func (h *Handler) UploadAnexo(c *gin.Context) {
	radicadoID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id de documento radicado inválido"})
		return
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "archivo requerido"})
		return
	}

	rutaBase := "storage"

	response, err := h.service.UploadAnexo(uint(radicadoID), fileHeader, rutaBase)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user, hasUser := c.Get("user")
	if hasUser {
		h.db.Create(&db.Trazabilidad{
			DocumentoRadicadoID: uint(radicadoID),
			UsuarioID:           user.(db.Usuario).ID,
			Accion:              "Archivo Subido",
			Descripcion:         "Se adjuntó el archivo: " + fileHeader.Filename,
			Fecha:               time.Now(),
		})
	}

	c.JSON(http.StatusCreated, response)
}
func (h *Handler) Download(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return
	}

	archivo, err := h.service.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// inline = preview en navegador, attachment = descarga forzada
	disposition := "inline"
	if c.Query("download") == "1" {
		disposition = "attachment"
	}

	c.Header("Content-Disposition", fmt.Sprintf("%s; filename=\"%s\"", disposition, archivo.Nombre))
	c.File(archivo.Ruta)
}

func (h *Handler) Reemplazar(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "archivo requerido"})
		return
	}

	response, err := h.service.Reemplazar(uint(id), fileHeader)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user, hasUser := c.Get("user")
	if hasUser && response != nil && response.DocumentoRadicadoID != 0 {
		h.db.Create(&db.Trazabilidad{
			DocumentoRadicadoID: response.DocumentoRadicadoID,
			UsuarioID:           user.(db.Usuario).ID,
			Accion:              "Archivo Editado",
			Descripcion:         "Se editó/reemplazó el archivo: " + response.Nombre,
			Fecha:               time.Now(),
		})
	}

	c.JSON(http.StatusOK, response)
}
