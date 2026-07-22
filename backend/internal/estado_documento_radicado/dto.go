package estado_documento_radicado

type CreateRequest struct {
	Nombre string `json:"nombre" binding:"required"`
}

type UpdateRequest struct {
	Nombre string `json:"nombre" binding:"required"`
}

type UpdateStatusRequest struct {
	Activo bool `json:"activo"`
}
