package responsabilidad_fiscal

import "sigefae/internal/db"

type Response struct {
	ID          uint          `json:"id"`
	IDProveedor uint          `json:"id_proveedor"`
	Proveedor   *db.Proveedor `json:"proveedor,omitempty"`
	Codigo      string        `json:"codigo"`
	Activo      bool          `json:"activo"`
}

func toResponse(responsabilidad db.ResponsabilidadFiscal) Response {

	return Response{
		ID:          responsabilidad.ID,
		IDProveedor: responsabilidad.IDProveedor,
		Proveedor:   responsabilidad.Proveedor,
		Codigo:      responsabilidad.Codigo,
		Activo:      responsabilidad.Activo,
	}
}