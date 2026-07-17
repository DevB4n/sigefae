package ruta

type CreateRequest struct {
	Nombre string `json:"nombre" binding:"required"`
	AreaID uint   `json:"area_id" binding:"required"`
}

type UpdateStatusRequest struct {
	Activo bool `json:"activo"`
}