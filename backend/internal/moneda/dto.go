package moneda

type CreateRequest struct {
	Nombre string `json:"nombre" binding:"required"`
	Codigo string `json:"codigo" binding:"required"`
}

type UpdateRequest struct {
	Nombre string `json:"nombre" binding:"required"`
	Codigo string `json:"codigo" binding:"required"`
}

type UpdateStatusRequest struct {
	Activo bool `json:"activo"`
}