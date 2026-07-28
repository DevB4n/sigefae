package documento_radicado

type CreateDTO struct {
	DocumentoComercialID uint   `json:"documento_comercial_id" binding:"required"`
	TipoRadicacionID     uint   `json:"tipo_radicacion_id" binding:"required"`
	RutaID               uint   `json:"ruta_id" binding:"required"`
	MetodoPagoID         uint   `json:"metodo_pago_id" binding:"required"`
	NumeroRadicado       string `json:"numero_radicado"` // opcional, si vacío se autogenera
}
type UpdateDTO struct {
	TipoRadicacionID uint   `json:"tipo_radicacion_id"`
	RutaID           uint   `json:"ruta_id"`
	MetodoPagoID     uint   `json:"metodo_pago_id"`
	NumeroRadicado   string `json:"numero_radicado"`
	UsuarioActualID  uint   `json:"usuario_actual_id"`
	EstadoPosesion   string `json:"estado_posesion"`
	PasoActualID     uint   `json:"paso_actual_id"`
	EstadoID         uint   `json:"estado_id"`
}