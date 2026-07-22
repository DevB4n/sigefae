package tipo_radicacion

type CreateRequest struct {
	Nombre string `json:"nombre" binding:"required"`
}

type UpdateRequest struct {
	Nombre string `json:"nombre" binding:"required"`
}

type UpdateStatusRequest struct {
	Activo bool `json:"activo"`
}
