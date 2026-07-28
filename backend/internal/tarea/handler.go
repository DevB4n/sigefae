package tarea

import (
	"net/http"
	"strconv"
	"time"

	"sigefae/internal/db"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	response, err := h.service.Create(request)
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
	var request UpdateDTO
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	response, err := h.service.Update(uint(id), request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response)
}

// ── NUEVO: Listar tareas por radicado ──
func (h *Handler) ListByRadicado(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return
	}

	var tareas []db.Tarea
	if err := h.db.Preload("UsuarioAsignado").Preload("Estado").
		Where("documento_radicado_id = ?", id).
		Order("id asc").
		Find(&tareas).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tareas)
}

// ── NUEVO: Completar tarea y avanzar flujo ──
func (h *Handler) Completar(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	user := c.MustGet("user").(db.Usuario)

	var tarea db.Tarea
	if err := h.db.Preload("DocumentoRadicado").Preload("Estado").First(&tarea, id).Error; err != nil {
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

	tarea.EstadoID = estadoCompletado.ID
	tarea.FechaFinalizacion = &now
	h.db.Save(&tarea)

	var siguiente db.Tarea
	if err := h.db.Where("documento_radicado_id = ? AND id > ? AND estado_id != ?",
		tarea.DocumentoRadicadoID, tarea.ID, estadoCompletado.ID).
		Order("id asc").First(&siguiente).Error; err == nil {

		var estadoEnProceso db.EstadoTarea
		h.db.Where("nombre = ?", "En Proceso").First(&estadoEnProceso)
		if estadoEnProceso.ID == 0 {
			estadoEnProceso.ID = 2
		}

		siguiente.EstadoID = estadoEnProceso.ID
		siguiente.FechaInicio = &now
		h.db.Save(&siguiente)

		h.db.Model(&db.DocumentoRadicado{}).Where("id = ?", tarea.DocumentoRadicadoID).Updates(map[string]interface{}{
			"usuario_actual_id": siguiente.UsuarioAsignadoID,
			"estado_posesion":   "EnProceso",
		})
	} else {
		h.db.Model(&db.DocumentoRadicado{}).Where("id = ?", tarea.DocumentoRadicadoID).Updates(map[string]interface{}{
			"estado_posesion": "Completado",
		})
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tarea completada", "tarea": tarea})
}