package notificacion

import (
	"time"

	"sigefae/internal/db"
)

type UsuarioResponse struct {
	ID     uint   `json:"id"`
	Nombre string `json:"nombre"`
	Email  string `json:"email"`
}

type DocumentoRadicadoResponse struct {
	ID             uint   `json:"id"`
	NumeroRadicado string `json:"numero_radicado"`
	EstadoPosesion string `json:"estado_posesion"`
}

type Response struct {
	ID                  uint                       `json:"id"`
	UsuarioID           uint                       `json:"usuario_id"`
	Usuario             *UsuarioResponse           `json:"usuario,omitempty"`
	DocumentoRadicadoID *uint                      `json:"documento_radicado_id,omitempty"`
	DocumentoRadicado   *DocumentoRadicadoResponse `json:"documento_radicado,omitempty"`
	Mensaje             string                     `json:"mensaje"`
	Estado              string                     `json:"estado"`
	Tipo                string                     `json:"tipo"`
	FechaCreacion       time.Time                  `json:"fecha_creacion"`
	FechaEnvio          *time.Time                 `json:"fecha_envio,omitempty"`
	FechaLectura        *time.Time                 `json:"fecha_lectura,omitempty"`
}

func toResponse(n db.Notificacion) Response {

	response := Response{
		ID:                  n.ID,
		UsuarioID:           n.UsuarioID,
		DocumentoRadicadoID: n.DocumentoRadicadoID,
		Mensaje:             n.Mensaje,
		Estado:              n.Estado,
		Tipo:                n.Tipo,
		FechaCreacion:       n.FechaCreacion,
		FechaEnvio:          n.FechaEnvio,
		FechaLectura:        n.FechaLectura,
	}

	if n.Usuario != nil {
		response.Usuario = &UsuarioResponse{
			ID:     n.Usuario.ID,
			Nombre: n.Usuario.Nombre,
			Email:  n.Usuario.Email,
		}
	}

	if n.DocumentoRadicado != nil {
		response.DocumentoRadicado = &DocumentoRadicadoResponse{
			ID:             n.DocumentoRadicado.ID,
			NumeroRadicado: n.DocumentoRadicado.NumeroRadicado,
			EstadoPosesion: n.DocumentoRadicado.EstadoPosesion,
		}
	}

	return response
}
