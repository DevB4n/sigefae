package tipo_documento

type CreateDTO struct {
	Nombre string `json:"nombre" binding:"required"`
}

type UpdateDTO struct {
	Nombre string `json:"nombre" binding:"required"`
}

type UpdateStatusDTO struct {
	Activo bool `json:"activo"`
}