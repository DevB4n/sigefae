package tarea

import (
	"time"

	"sigefae/internal/db"
)

type Response struct {
	ID                  uint       `json:"id"`
	DocumentoRadicadoID uint       `json:"documento_radicado_id"`
	UsuarioAsignadoID   uint       `json:"usuario_asignado_id"`
	EstadoID            uint       `json:"estado_id"`
	Descripcion         string     `json:"descripcion"`
	FechaAsignacion     time.Time  `json:"fecha_asignacion"`
	FechaInicio         *time.Time `json:"fecha_inicio"`
	FechaLimite         *time.Time `json:"fecha_limite"`
	FechaFinalizacion   *time.Time `json:"fecha_finalizacion"`
	CreatedAt           time.Time  `json:"created_at"`
}

func toResponse(tarea db.Tarea) Response {

	return Response{
		ID:                  tarea.ID,
		DocumentoRadicadoID: tarea.DocumentoRadicadoID,
		UsuarioAsignadoID:   tarea.UsuarioAsignadoID,
		EstadoID:            tarea.EstadoID,
		Descripcion:         tarea.Descripcion,
		FechaAsignacion:     tarea.FechaAsignacion,
		FechaInicio:         tarea.FechaInicio,
		FechaLimite:         tarea.FechaLimite,
		FechaFinalizacion:   tarea.FechaFinalizacion,
		CreatedAt:           tarea.CreatedAt,
	}
}