package role

type CreateRequest struct {
	Nombre string `json:"nombre" binding:"required"`
}

type UpdateRequest struct {
	Nombre string `json:"nombre" binding:"required"`
}

type Response struct {
	ID     uint   `json:"id"`
	Nombre string `json:"nombre"`
	Activo bool   `json:"activo"`
}

type UpdateStatusRequest struct {
	Activo bool `json:"activo"`
}