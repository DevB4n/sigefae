package direccion

type CreateDTO struct {
	Nombre       string `json:"nombre" binding:"required"`
	Linea1       string `json:"linea_1" binding:"required"`
	Linea2       string `json:"linea_2"`
	CodigoPostal string `json:"codigo_postal"`
	IDMunicipio  uint   `json:"municipio_id" binding:"required"`
}

type UpdateDTO struct {
	Nombre       string `json:"nombre" binding:"required"`
	Linea1       string `json:"linea_1" binding:"required"`
	Linea2       string `json:"linea_2"`
	CodigoPostal string `json:"codigo_postal"`
	IDMunicipio  uint   `json:"municipio_id" binding:"required"`
}

type UpdateStatusDTO struct {
	Activo bool `json:"activo"`
}