package metodo_pago

type CreateRequest struct {
	TipoPagoID uint   `json:"tipo_pago_id" binding:"required"`
	Nombre     string `json:"nombre" binding:"required"`
}

type UpdateStatusRequest struct {
	Activo bool `json:"activo"`
}