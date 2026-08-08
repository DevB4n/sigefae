package tarea

import (
	"fmt"
	"math"          // ← AGREGAR
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

	var tareaRetornoID *uint
	if tarea.DocumentoRadicado != nil {
		tareaRetornoID = tarea.DocumentoRadicado.TareaPendienteRetornoID
	}
	var saltoDirecto = false

	// VERIFICAR RETORNO DIRECTO
	if tareaRetornoID != nil && *tareaRetornoID != 0 {
		if err := h.db.First(&siguiente, *tareaRetornoID).Error; err == nil {
			saltoDirecto = true
		}
	}

	if saltoDirecto {
		var estadoEnProceso db.EstadoTarea
		h.db.Where("nombre = ?", "En Proceso").First(&estadoEnProceso)
		if estadoEnProceso.ID == 0 {
			estadoEnProceso.ID = 2
		}

		h.db.Model(&siguiente).Updates(map[string]interface{}{
			"estado_id":          estadoEnProceso.ID,
			"fecha_inicio":       now,
			"fecha_finalizacion": gorm.Expr("NULL"),
		})

		// Eliminar la marca de retorno (para que si se devuelve de nuevo, se pueda)
		h.db.Model(&db.DocumentoRadicado{}).Where("id = ?", tarea.DocumentoRadicadoID).Updates(map[string]interface{}{
			"usuario_actual_id":           siguiente.UsuarioAsignadoID,
			"estado_posesion":             "EnProceso",
			"tarea_pendiente_retorno_id": gorm.Expr("NULL"),
		})
		notificarA = siguiente.UsuarioAsignadoID
	} else {
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
		// ── VALIDACIÓN: es el último paso, las normas deben sumar 100 % ──
		var totalPorcentaje float64
		h.db.Raw("SELECT COALESCE(SUM(porcentaje),0) FROM documento_radicado_norma_repartos WHERE documento_radicado_id = ?",
			tarea.DocumentoRadicadoID).Scan(&totalPorcentaje)

		if math.Abs(totalPorcentaje-100.0) > 0.01 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("Las normas de reparto suman %.2f %%. Deben sumar exactamente 100 %% antes de finalizar el proceso.", totalPorcentaje),
			})
			return
		}

		h.db.Model(&db.DocumentoRadicado{}).Where("id = ?", tarea.DocumentoRadicadoID).Updates(map[string]interface{}{
			"estado_posesion": "Completado",
		})
		}
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

func (h *Handler) Devolver(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return
	}

	var dto DevolverDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := c.MustGet("user").(db.Usuario)

	var tarea db.Tarea
	if err := h.db.Preload("DocumentoRadicado").First(&tarea, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tarea actual no encontrada"})
		return
	}

	if tarea.UsuarioAsignadoID != user.ID && user.Rol.Nombre != "Superadministrador" {
		c.JSON(http.StatusForbidden, gin.H{"error": "No estás asignado a esta tarea"})
		return
	}

	var destino db.Tarea
	if err := h.db.First(&destino, dto.TareaDestinoID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tarea destino no encontrada"})
		return
	}

	// ===== VALIDACIONES QUE FALTAN =====
	// 1. Mismo documento
	if destino.DocumentoRadicadoID != tarea.DocumentoRadicadoID {
	    c.JSON(http.StatusBadRequest, gin.H{"error": "La tarea destino no pertenece a este documento"})
	    return
	}
	
	// 2. No devolver a sí misma
	if destino.ID == tarea.ID {
	    c.JSON(http.StatusBadRequest, gin.H{"error": "No puede devolver a la misma tarea"})
	    return
	}
	
	// 3. Solo devolver a tareas COMPLETADAS (no a devueltas ni pendientes)
	var estadoDestino db.EstadoTarea
	if err := h.db.First(&estadoDestino, destino.EstadoID).Error; err != nil {
	    c.JSON(http.StatusInternalServerError, gin.H{"error": "Error verificando estado de tarea destino"})
	    return
	}
	if estadoDestino.Nombre != "Completada" {
	    c.JSON(http.StatusBadRequest, gin.H{"error": "Solo puede devolver a una tarea que haya sido completada"})
	    return
	}
	// ===================================

	// Estados
	var estadoDevuelta db.EstadoTarea
	h.db.Where("nombre = ?", "Devuelta").First(&estadoDevuelta)
	if estadoDevuelta.ID == 0 {
		estadoDevuelta.ID = 4
	}

	var estadoEnProceso db.EstadoTarea
	h.db.Where("nombre = ?", "En Proceso").First(&estadoEnProceso)
	if estadoEnProceso.ID == 0 {
		estadoEnProceso.ID = 2
	}

	now := time.Now()

	// Actualizar tarea actual a Devuelta
	if err := h.db.Model(&tarea).Updates(map[string]interface{}{
		"estado_id":          estadoDevuelta.ID,
		"fecha_finalizacion": now,
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error actualizando tarea actual: " + err.Error()})
		return
	}

	// Reactivar tarea destino (limpiar fecha_finalizacion para que se muestre como activa)
	if err := h.db.Model(&destino).Updates(map[string]interface{}{
		"estado_id":          estadoEnProceso.ID,
		"fecha_inicio":       now,
		"fecha_finalizacion": gorm.Expr("NULL"),
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reactivando tarea destino: " + err.Error()})
		return
	}

	// Actualizar Radicado
	radicadoUpdates := map[string]interface{}{
		"usuario_actual_id": destino.UsuarioAsignadoID,
		"estado_posesion":   "Devuelto",
	}
	if dto.RetornoDirecto {
		radicadoUpdates["tarea_pendiente_retorno_id"] = tarea.ID
	} else {
		radicadoUpdates["tarea_pendiente_retorno_id"] = gorm.Expr("NULL")
	}

	if err := h.db.Model(&db.DocumentoRadicado{}).Where("id = ?", tarea.DocumentoRadicadoID).Updates(radicadoUpdates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error actualizando radicado: " + err.Error()})
		return
	}

	// Registrar Aprobacion, Comentario y Trazabilidad
	h.db.Create(&db.RegistroAprobacion{
		DocumentoRadicadoID: tarea.DocumentoRadicadoID,
		ResponsableID:       user.ID,
		RolID:               user.RolID,
		Estado:              "DEVUELTO",
		Observacion:         dto.Observacion,
		Fecha:               now,
	})

	h.db.Create(&db.Comentario{
		DocumentoRadicadoID: tarea.DocumentoRadicadoID,
		UsuarioID:           user.ID,
		Descripcion:         "DEVOLUCIÓN: " + dto.Observacion,
		Fecha:               now,
	})

	var usuarioDestino db.Usuario
	h.db.First(&usuarioDestino, destino.UsuarioAsignadoID)
	nombreDestino := usuarioDestino.Nombre
	if nombreDestino == "" {
		nombreDestino = "Usuario Desconocido"
	}

	h.db.Create(&db.Trazabilidad{
		DocumentoRadicadoID: tarea.DocumentoRadicadoID,
		UsuarioID:           user.ID,
		Accion:              "Tarea Devuelta",
		Descripcion:         fmt.Sprintf("El usuario devolvió el documento a %s. Motivo: %s", nombreDestino, dto.Observacion),
		Fecha:               now,
	})

	// Notificar
	if h.notifSvc != nil {
		docID := tarea.DocumentoRadicadoID
		numRad := ""
		if tarea.DocumentoRadicado != nil {
			numRad = tarea.DocumentoRadicado.NumeroRadicado
		}
		h.notifSvc.CreateFromEvent(notificacion.CreateDTO{
			UsuarioID:           destino.UsuarioAsignadoID,
			DocumentoRadicadoID: &docID,
			Mensaje:             fmt.Sprintf("Te devolvieron el radicado %s para revisión. Motivo: %s", numRad, dto.Observacion),
			Estado:              "Pendiente",
			Tipo:                "Devolucion",
			FechaCreacion:       now,
		})
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tarea devuelta correctamente"})
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
