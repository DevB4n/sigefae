package categoria_proveedor

type CreateRequest struct {
	Nombre      string `json:"nombre" binding:"required"`
	Descripcion string `json:"descripcion"`
}

type UpdateRequest struct {
	Nombre      string `json:"nombre" binding:"required"`
	Descripcion string `json:"descripcion"`
}

type UpdateStatusRequest struct {
	Activo bool `json:"activo"`
}
