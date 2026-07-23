package notificacion

import "time"

type CreateDTO struct {
	UsuarioID           uint       `json:"usuario_id" binding:"required"`
	DocumentoRadicadoID *uint      `json:"documento_radicado_id"`
	Mensaje             string     `json:"mensaje" binding:"required"`
	Estado              string     `json:"estado" binding:"required"`
	Tipo                string     `json:"tipo" binding:"required"`
	FechaCreacion       time.Time  `json:"fecha_creacion" binding:"required"`
	FechaEnvio          *time.Time `json:"fecha_envio"`
	FechaLectura        *time.Time `json:"fecha_lectura"`
}

type UpdateDTO struct {
	UsuarioID           uint       `json:"usuario_id" binding:"required"`
	DocumentoRadicadoID *uint      `json:"documento_radicado_id"`
	Mensaje             string     `json:"mensaje" binding:"required"`
	Estado              string     `json:"estado" binding:"required"`
	Tipo                string     `json:"tipo" binding:"required"`
	FechaCreacion       time.Time  `json:"fecha_creacion" binding:"required"`
	FechaEnvio          *time.Time `json:"fecha_envio"`
	FechaLectura        *time.Time `json:"fecha_lectura"`
}
