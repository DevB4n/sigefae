package documento_radicado

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"sigefae/internal/db"
	"sigefae/internal/notificacion"
)

type Handler struct {
	service  *Service
	notifSvc *notificacion.Service
	db       *gorm.DB
}

func NewHandler(service *Service, notifSvc *notificacion.Service, database *gorm.DB) *Handler {
	return &Handler{
		service:  service,
		notifSvc: notifSvc,
		db:       database,
	}
}

func (h *Handler) Create(c *gin.Context) {
	user := c.MustGet("user").(db.Usuario)

	var dto CreateDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.service.Create(dto, user.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ── NOTIFICAR al primer responsable ──
	if h.notifSvc != nil && response.UsuarioActualID != 0 && response.UsuarioActualID != user.ID {
		docID := response.ID
		h.notifSvc.CreateFromEvent(notificacion.CreateDTO{
			UsuarioID:           response.UsuarioActualID,
			DocumentoRadicadoID: &docID,
			Mensaje:             "Nuevo radicado " + response.NumeroRadicado + " requiere tu revisión",
			Estado:              "Pendiente",
			Tipo:                "Asignacion",
			FechaCreacion:       time.Now(),
		})
	}

	// ── REGISTRAR Trazabilidad ──
	descTrazabilidad := "Documento radicado e iniciado en el flujo."
	if response.UsuarioActualID != 0 {
		var usuarioSiguiente db.Usuario
		h.db.First(&usuarioSiguiente, response.UsuarioActualID)
		nombreSiguiente := usuarioSiguiente.Nombre
		if nombreSiguiente == "" {
			nombreSiguiente = "Usuario Desconocido"
		}
		descTrazabilidad = fmt.Sprintf("Documento radicado. Asignado inicialmente a: %s", nombreSiguiente)
	}

	h.db.Create(&db.Trazabilidad{
		DocumentoRadicadoID: response.ID,
		UsuarioID:           user.ID,
		Accion:              "Radicación",
		Descripcion:         descTrazabilidad,
		Fecha:               time.Now(),
	})

	c.JSON(http.StatusCreated, response)
}

func (h *Handler) List(c *gin.Context) {
	response, err := h.service.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response)
}

func (h *Handler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return
	}
	response, err := h.service.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response)
}

func (h *Handler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return
	}
	var dto UpdateDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	response, err := h.service.Update(uint(id), dto)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response)
}

func (h *Handler) VerificarPublico(c *gin.Context) {
	numero := c.Param("numero_radicado")

	radicado, err := h.service.GetByNumeroRadicado(numero)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"valido": false,
			"error":  "El documento no existe en nuestros registros",
		})
		return
	}

	// Preparar respuesta pública (no exponemos todo el modelo)
	response := gin.H{
		"valido":           true,
		"numero_radicado":  radicado.NumeroRadicado,
		"fecha_radicacion": radicado.FechaRadicacion,
		"tipo_radicacion":  radicado.TipoRadicacion.Nombre,
		"estado_posesion":  radicado.EstadoPosesion,
		"documento": gin.H{
			"tipo":          radicado.DocumentoComercial.Tipo,
			"numero":        radicado.DocumentoComercial.NumeroDocumento,
			"fecha_emision": radicado.DocumentoComercial.FechaDocumento,
			"total":         radicado.DocumentoComercial.Total,
			"moneda":        radicado.DocumentoComercial.Moneda.Nombre,
			"cufe":          radicado.DocumentoComercial.Cufe, // ← verifica que este campo exista en tu modelo
		},
		"proveedor": gin.H{
			"razon_social": radicado.DocumentoComercial.Proveedor.RazonSocial,
			"nit":          radicado.DocumentoComercial.Proveedor.NumeroDocumento,
		},
		"receptor": gin.H{
			"nombre": radicado.DocumentoComercial.Receptor.Nombre,
			"nit":    radicado.DocumentoComercial.Receptor.NumeroDocumento,
		},
		"qr_url": radicado.Qr.Url,
	}

	c.JSON(http.StatusOK, response)
}
func (h *Handler) GetNormasReparto(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return
	}
	normas, err := h.service.GetNormasReparto(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, normas)
}

func (h *Handler) AsignarNormasReparto(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return
	}
	var dto struct {
		Normas []NormaRepartoInputDTO `json:"normas" binding:"required,dive"`
	}
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.AsignarNormasReparto(uint(id), dto.Normas); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "normas asignadas correctamente"})
}
