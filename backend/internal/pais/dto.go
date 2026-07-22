package pais

type CreateDTO struct {
	Nombre string `json:"nombre" binding:"required"`
	Codigo string `json:"codigo" binding:"required"`
}

type UpdateDTO struct {
	Nombre string `json:"nombre" binding:"required"`
	Codigo string `json:"codigo" binding:"required"`
}

type UpdateStatusDTO struct {
	Activo bool `json:"activo"`
}
