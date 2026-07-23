package tarea

import "time"

type CreateDTO struct {
	DocumentoRadicadoID uint       `json:"documento_radicado_id" binding:"required"`
	UsuarioAsignadoID   uint       `json:"usuario_asignado_id" binding:"required"`
	EstadoID            uint       `json:"estado_id" binding:"required"`
	Descripcion         string     `json:"descripcion" binding:"required"`
	FechaLimite         *time.Time `json:"fecha_limite"`
}

type UpdateDTO struct {
	UsuarioAsignadoID uint       `json:"usuario_asignado_id" binding:"required"`
	EstadoID          uint       `json:"estado_id" binding:"required"`
	Descripcion       string     `json:"descripcion" binding:"required"`
	FechaInicio       *time.Time `json:"fecha_inicio"`
	FechaLimite       *time.Time `json:"fecha_limite"`
	FechaFinalizacion *time.Time `json:"fecha_finalizacion"`
}
