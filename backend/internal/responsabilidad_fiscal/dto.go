package responsabilidad_fiscal

type CreateRequest struct {
	IDProveedor uint   `json:"id_proveedor" binding:"required"`
	Codigo      string `json:"codigo" binding:"required"`
}

type UpdateRequest struct {
	IDProveedor uint   `json:"id_proveedor" binding:"required"`
	Codigo      string `json:"codigo" binding:"required"`
}

type UpdateStatusRequest struct {
	Activo bool `json:"activo"`
}
