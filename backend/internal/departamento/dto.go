package departamento

type CreateDTO struct {
	Nombre string `json:"nombre" binding:"required"`
	PaisID uint   `json:"pais_id" binding:"required"`
}

type UpdateDTO struct {
	Nombre string `json:"nombre" binding:"required"`
	PaisID uint   `json:"pais_id" binding:"required"`
}

type UpdateStatusDTO struct {
	Activo bool `json:"activo"`
}
