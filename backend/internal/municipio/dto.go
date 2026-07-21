package municipio

type CreateDTO struct {
	Nombre         string `json:"nombre" binding:"required"`
	DepartamentoID uint   `json:"departamento_id" binding:"required"`
}

type UpdateDTO struct {
	Nombre         string `json:"nombre" binding:"required"`
	DepartamentoID uint   `json:"departamento_id" binding:"required"`
}

type UpdateStatusDTO struct {
	Activo bool `json:"activo"`
}