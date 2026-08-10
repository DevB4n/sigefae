package norma_reparto

type CreateDTO struct {
	Codigo       string  `json:"codigo" binding:"required"`
	Nombre       string  `json:"nombre" binding:"required"`
	Sucursal     string  `json:"sucursal" binding:"required"`
	Departamento string  `json:"departamento" binding:"required"`
	Tipo         *string `json:"tipo"`
	TarifaIva    *string `json:"tarifa_iva"`
}

type UpdateDTO struct {
	Codigo       string  `json:"codigo"`
	Nombre       string  `json:"nombre"`
	Sucursal     string  `json:"sucursal"`
	Departamento string  `json:"departamento"`
	Tipo         *string `json:"tipo"`
	TarifaIva    *string `json:"tarifa_iva"`
}

type UpdateStatusDTO struct {
	Activo bool `json:"activo"`
}
