package comentario

import (
	"time"

	"sigefae/internal/db"
)

type Response struct {
	ID                  uint      `json:"id"`
	DocumentoRadicadoID uint      `json:"documento_radicado_id"`
	UsuarioID           uint      `json:"usuario_id"`
	UsuarioNombre       string    `json:"usuario_nombre"`  // <-- nuevo
	Descripcion         string    `json:"descripcion"`
	Fecha               time.Time `json:"fecha"`
}

func toResponse(comentario db.Comentario) Response {

	nombreUsuario := ""

	if comentario.Usuario.ID != 0 {
		nombreUsuario = comentario.Usuario.Nombre
	}

	return Response{
		ID:                  comentario.ID,
		DocumentoRadicadoID: comentario.DocumentoRadicadoID,
		Descripcion:         comentario.Descripcion,
		UsuarioID:           comentario.UsuarioID,
		UsuarioNombre:       nombreUsuario,
		Fecha:               comentario.Fecha,
	}
}