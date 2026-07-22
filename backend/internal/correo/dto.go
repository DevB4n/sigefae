package correo

import (
	"time"
)

type CreateRequest struct {
	Asunto         string    `json:"asunto" binding:"required"`
	De             string    `json:"de" binding:"required"`
	Para           string    `json:"para" binding:"required"`
	FechaRecepcion time.Time `json:"fecha_recepcion" binding:"required"`
	IDMensaje      string    `json:"id_mensaje" binding:"required"`
	Cuerpo         string    `json:"cuerpo"`
	Cc             string    `json:"cc"`
	Bcc            string    `json:"bcc"`
	ReplyTo        string    `json:"reply_to"`
	IDEstado       uint      `json:"id_estado" binding:"required"`
}

type UpdateRequest struct {
	Asunto         string    `json:"asunto" binding:"required"`
	De             string    `json:"de" binding:"required"`
	Para           string    `json:"para" binding:"required"`
	FechaRecepcion time.Time `json:"fecha_recepcion" binding:"required"`
	IDMensaje      string    `json:"id_mensaje" binding:"required"`
	Cuerpo         string    `json:"cuerpo"`
	Cc             string    `json:"cc"`
	Bcc            string    `json:"bcc"`
	ReplyTo        string    `json:"reply_to"`
	IDEstado       uint      `json:"id_estado" binding:"required"`
}
type UpdateStatusRequest struct {
	IDEstado uint `json:"id_estado"`
}