package tipo_factura

type CreateRequest struct {
	AreaID uint   `json:"area_id" binding:"required"`
	Nombre string `json:"nombre" binding:"required"`
}

type UpdateStatusRequest struct {
	Activo bool `json:"activo"`
}