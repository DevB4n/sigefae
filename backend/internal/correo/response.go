package correo

import (
	"time"

	"sigefae/internal/db"
)

type EstadoCorreoResponse struct {
	ID     uint   `json:"id"`
	Nombre string `json:"nombre"`
}

type Response struct {
	ID             uint                  `json:"id"`
	Asunto         string                `json:"asunto"`
	De             string                `json:"de"`
	Para           string                `json:"para"`
	FechaRecepcion time.Time             `json:"fecha_recepcion"`
	IDMensaje      string                `json:"id_mensaje"`
	Cuerpo         string                `json:"cuerpo"`
	Cc             string                `json:"cc"`
	Bcc            string                `json:"bcc"`
	ReplyTo        string                `json:"reply_to"`
	IDEstado       uint                  `json:"id_estado"`
	EstadoCorreo   *EstadoCorreoResponse `json:"estado_correo,omitempty"`
	Archivos       []string              `json:"archivos,omitempty"`
	Activo         bool                  `json:"activo"`
}

func toResponse(correo db.Correo) Response {

	response := Response{
		ID:             correo.ID,
		Asunto:         correo.Asunto,
		De:             correo.De,
		Para:           correo.Para,
		FechaRecepcion: correo.FechaRecepcion,
		IDMensaje:      correo.IDMensaje,
		Cuerpo:         correo.Cuerpo,
		Cc:             correo.Cc,
		Bcc:            correo.Bcc,
		ReplyTo:        correo.ReplyTo,
		IDEstado:       correo.IDEstado,
		Activo:         correo.Activo,
	}

	if correo.EstadoCorreo != nil {

		response.EstadoCorreo = &EstadoCorreoResponse{
			ID:     correo.EstadoCorreo.ID,
			Nombre: correo.EstadoCorreo.Nombre,
		}
	}

	return response
}
