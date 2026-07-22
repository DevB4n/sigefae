package contacto

import "sigefae/internal/db"

type Response struct {
	ID uint `json:"id"`

	ProveedorID uint   `json:"proveedor_id"`
	Proveedor   string `json:"proveedor"`

	Nombre   string `json:"nombre"`
	Cargo    string `json:"cargo"`
	Telefono string `json:"telefono"`
	Correo   string `json:"correo"`

	Activo bool `json:"activo"`
}

func toResponse(c db.Contacto) Response {

	proveedor := ""

	if c.Proveedor != nil {
		proveedor = c.Proveedor.RazonSocial
	}

	return Response{
		ID: c.ID,

		ProveedorID: c.ProveedorID,
		Proveedor:   proveedor,

		Nombre:   c.Nombre,
		Cargo:    c.Cargo,
		Telefono: c.Telefono,
		Correo:   c.Correo,

		Activo: c.Activo,
	}
}
