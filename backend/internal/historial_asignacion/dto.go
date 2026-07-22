package historial_asignacion

import "time"

type CreateDTO struct {
	DocumentoRadicadoID uint `json:"documento_radicado_id" binding:"required"`
	UsuarioID           uint `json:"usuario_id" binding:"required"`
}


type CerrarAsignacionDTO struct {
	Hasta time.Time `json:"hasta"`//se cierra automaticamente cuando se hace una nueva asignacion a ese documento radicado, no hay otro metodo
}