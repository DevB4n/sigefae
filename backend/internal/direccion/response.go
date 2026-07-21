package direccion

import "sigefae/internal/db"

type Response struct {
	ID             uint   `json:"id"`

	Nombre         string `json:"nombre"`

	Linea1         string `json:"linea_1"`
	Linea2         string `json:"linea_2"`

	CodigoPostal   string `json:"codigo_postal"`

	IDMunicipio    uint   `json:"id_municipio"`
	Municipio      string `json:"municipio"`

	Activo         bool   `json:"activo"`
}

func toResponse(d db.Direccion) Response {

	municipio := ""

	if d.Municipio != nil {
		municipio = d.Municipio.Nombre
	}

	return Response{
		ID:           d.ID,
		Nombre:       d.Nombre,
		Linea1:       d.Linea1,
		Linea2:       d.Linea2,
		CodigoPostal: d.CodigoPostal,
		IDMunicipio:  d.IDMunicipio,
		Municipio:    municipio,
		Activo:       d.Activo,
	}
}