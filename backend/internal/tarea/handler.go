package tarea

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

func (h *Handler) ListByRadicado(c *gin.Context) {
	id := c.Param("id")
	radicadoID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return
	}

	var tareas []db.Tarea
	if err := h.db.Where("documento_radicado_id = ?", radicadoID).
		Preload("Estado").
		Preload("UsuarioAsignado").
		Order("id asc").
		Find(&tareas).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tareas)
}

func (h *Handler) Completar(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return
	}

	user := c.MustGet("user").(db.Usuario)

	var tarea db.Tarea
	if err := h.db.Preload("DocumentoRadicado").First(&tarea, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tarea no encontrada"})
		return
	}

	if tarea.UsuarioAsignadoID != user.ID && user.Rol.Nombre != "Superadministrador" {
		c.JSON(http.StatusForbidden, gin.H{"error": "No estás asignado a esta tarea"})
		return
	}

	now := time.Now()

	var estadoCompletado db.EstadoTarea
	h.db.Where("nombre = ?", "Completada").First(&estadoCompletado)
	if estadoCompletado.ID == 0 {
		estadoCompletado.ID = 3
	}

	if tarea.EstadoID == estadoCompletado.ID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "La tarea ya está completada"})
		return
	}

	if err := h.db.Model(&tarea).Updates(map[string]interface{}{
		"estado_id":          estadoCompletado.ID,
		"fecha_finalizacion": now,
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error actualizando tarea"})
		return
	}

	var siguiente db.Tarea
	var notificarA uint = 0
	var numRadicado string = ""

	if tarea.DocumentoRadicado != nil {
		numRadicado = tarea.DocumentoRadicado.NumeroRadicado
	}

	if err := h.db.Where("documento_radicado_id = ? AND id > ? AND estado_id != ?",
		tarea.DocumentoRadicadoID, tarea.ID, estadoCompletado.ID).
		Order("id asc").First(&siguiente).Error; err == nil {

		var estadoEnProceso db.EstadoTarea
		h.db.Where("nombre = ?", "En Proceso").First(&estadoEnProceso)
		if estadoEnProceso.ID == 0 {
			estadoEnProceso.ID = 2
		}

		h.db.Model(&siguiente).Updates(map[string]interface{}{
			"estado_id":    estadoEnProceso.ID,
			"fecha_inicio": now,
		})

		h.db.Model(&db.DocumentoRadicado{}).Where("id = ?", tarea.DocumentoRadicadoID).Updates(map[string]interface{}{
			"usuario_actual_id": siguiente.UsuarioAsignadoID,
			"estado_posesion":   "EnProceso",
		})

		notificarA = siguiente.UsuarioAsignadoID
	} else {
		h.db.Model(&db.DocumentoRadicado{}).Where("id = ?", tarea.DocumentoRadicadoID).Updates(map[string]interface{}{
			"estado_posesion": "Completado",
		})
	}

	// ── NOTIFICAR al siguiente usuario ──
	if h.notifSvc != nil && notificarA != 0 && notificarA != user.ID {
		docID := tarea.DocumentoRadicadoID
		h.notifSvc.CreateFromEvent(notificacion.CreateDTO{
			UsuarioID:           notificarA,
			DocumentoRadicadoID: &docID,
			Mensaje:             "Te asignaron el radicado " + numRadicado,
			Estado:              "Pendiente",
			Tipo:                "Asignacion",
			FechaCreacion:       time.Now(),
		})
	}

	// ── REGISTRAR Trazabilidad ──
	var accionTrazabilidad string
	var descTrazabilidad string
	if notificarA != 0 {
		var usuarioSiguiente db.Usuario
		h.db.First(&usuarioSiguiente, notificarA)
		nombreSiguiente := usuarioSiguiente.Nombre
		if nombreSiguiente == "" {
			nombreSiguiente = "Usuario Desconocido"
		}
		
		accionTrazabilidad = "Paso Completado y Asignado"
		descTrazabilidad = fmt.Sprintf("El usuario completó su tarea. El documento fue asignado a: %s", nombreSiguiente)
	} else {
		accionTrazabilidad = "Proceso Finalizado"
		descTrazabilidad = "El usuario completó la última tarea. El proceso ha finalizado."
	}
	
	h.db.Create(&db.Trazabilidad{
		DocumentoRadicadoID: tarea.DocumentoRadicadoID,
		UsuarioID:           user.ID,
		Accion:              accionTrazabilidad,
		Descripcion:         descTrazabilidad,
		Fecha:               time.Now(),
	})

	c.JSON(http.StatusOK, gin.H{"message": "Tarea completada"})
}

func (h *Handler) Create(c *gin.Context) {
	var dto CreateDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	response, err := h.service.Create(dto)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
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