package receptor

import "sigefae/internal/db"

type Response struct {
	ID              uint               `json:"id"`
	Nombre          string             `json:"nombre"`
	TipoDocumentoID uint               `json:"tipo_documento_id"`
	TipoDocumento   *db.TipoDocumento  `json:"tipo_documento,omitempty"`
	NumeroDocumento string             `json:"numero_documento"`
	Activo          bool               `json:"activo"`
}

func toResponse(receptor db.Receptor) Response {

	return Response{
		ID:              receptor.ID,
		Nombre:          receptor.Nombre,
		TipoDocumentoID: receptor.TipoDocumentoID,
		TipoDocumento:   receptor.TipoDocumento,
		NumeroDocumento: receptor.NumeroDocumento,
		Activo:          receptor.Activo,
	}
}