package documento_radicado

import "time"

type CreateDTO struct {
	DocumentoComercialID   uint       `json:"documento_comercial_id" binding:"required"`
	TipoRadicacionID       uint       `json:"tipo_radicacion_id" binding:"required"`
	RutaID                 uint       `json:"ruta_id" binding:"required"`
	NumeroRadicado         string     `json:"numero_radicado" binding:"required"`
	FechaRadicacion        time.Time  `json:"fecha_radicacion" binding:"required"`
	UsuarioActualID        uint       `json:"usuario_actual_id" binding:"required"`
	EstadoPosesion         string     `json:"estado_posesion" binding:"required"`
	PasoActualID           *uint      `json:"paso_actual_id"`
	PasoPendienteRetornoID *uint      `json:"paso_pendiente_retorno_id"`
	EstadoID               uint       `json:"estado_id" binding:"required"`
	UltimaActividad        *time.Time `json:"ultima_actividad"`
	QrID                   uint       `json:"qr_id" binding:"required"`
	MetodoPagoID           uint       `json:"metodo_pago_id" binding:"required"`
}

type UpdateDTO struct {
	DocumentoComercialID   uint       `json:"documento_comercial_id" binding:"required"`
	TipoRadicacionID       uint       `json:"tipo_radicacion_id" binding:"required"`
	RutaID                 uint       `json:"ruta_id" binding:"required"`
	NumeroRadicado         string     `json:"numero_radicado" binding:"required"`
	FechaRadicacion        time.Time  `json:"fecha_radicacion" binding:"required"`
	UsuarioActualID        uint       `json:"usuario_actual_id" binding:"required"`
	EstadoPosesion         string     `json:"estado_posesion" binding:"required"`
	PasoActualID           *uint      `json:"paso_actual_id"`
	PasoPendienteRetornoID *uint      `json:"paso_pendiente_retorno_id"`
	EstadoID               uint       `json:"estado_id" binding:"required"`
	UltimaActividad        *time.Time `json:"ultima_actividad"`
	QrID                   uint       `json:"qr_id" binding:"required"`
	MetodoPagoID           uint       `json:"metodo_pago_id" binding:"required"`
}
