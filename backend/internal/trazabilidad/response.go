package trazabilidad

import (
	"time"

	"sigefae/internal/db"
)

type Response struct {
	ID                  uint      `json:"id"`
	DocumentoRadicadoID uint      `json:"documento_radicado_id"`
	UsuarioID           uint      `json:"usuario_id"`
	Accion              string    `json:"accion"`
	Descripcion         string    `json:"descripcion"`
	Fecha               time.Time `json:"fecha"`
}

func toResponse(trazabilidad db.Trazabilidad) Response {

	return Response{
		ID:                  trazabilidad.ID,
		DocumentoRadicadoID: trazabilidad.DocumentoRadicadoID,
		UsuarioID:           trazabilidad.UsuarioID,
		Accion:              trazabilidad.Accion,
		Descripcion:         trazabilidad.Descripcion,
		Fecha:               trazabilidad.Fecha,
	}
}
