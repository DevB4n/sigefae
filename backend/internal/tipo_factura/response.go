package tipo_factura

import "sigefae/internal/db"

type AreaResponse struct {
	ID     uint   `json:"id"`
	Nombre string `json:"nombre"`
}

type Response struct {
	ID     uint          `json:"id"`
	AreaID uint          `json:"area_id"`
	Area   *AreaResponse `json:"area,omitempty"`
	Nombre string        `json:"nombre"`
	Activo bool          `json:"activo"`
}

func toResponse(tipoFactura db.TipoFactura) Response {

	response := Response{
		ID:     tipoFactura.ID,
		AreaID: tipoFactura.AreaID,
		Nombre: tipoFactura.Nombre,
		Activo: tipoFactura.Activo,
	}

	if tipoFactura.Area != nil {
		response.Area = &AreaResponse{
			ID:     tipoFactura.Area.ID,
			Nombre: tipoFactura.Area.Nombre,
		}
	}

	return response
}