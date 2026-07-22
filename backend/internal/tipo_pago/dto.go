package tipo_pago

type CreateRequest struct {
	Nombre string `json:"nombre" binding:"required"`
}

type UpdateStatusRequest struct {
	Activo bool `json:"activo"`
}

type UpdateRequest struct {
	Nombre string `json:"nombre" binding:"required"`
}
