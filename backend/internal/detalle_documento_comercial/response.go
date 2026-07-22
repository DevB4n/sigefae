package detalle_documento_comercial

import "sigefae/internal/db"

type DocumentoComercialResponse struct {
	ID               uint   `json:"id"`
	Tipo             string `json:"tipo"`
	NumeroDocumento  string `json:"numero_documento"`
}

type Response struct {
	ID                   uint                         `json:"id"`
	DocumentoComercialID uint                         `json:"documento_comercial_id"`
	DocumentoComercial   *DocumentoComercialResponse  `json:"documento_comercial,omitempty"`
	Descripcion          string                       `json:"descripcion"`
	ValorUnit            float64                      `json:"valor_unitario"`
	IvaUnit              float64                      `json:"iva_unitario"`
	Cantidad             float64                      `json:"cantidad"`
	Total                float64                      `json:"total"`
	Activo               bool                         `json:"activo"`
}

func toResponse(detalle db.DetalleDocumentoComercial) Response {

	response := Response{
		ID:                   detalle.ID,
		DocumentoComercialID: detalle.DocumentoComercialID,
		Descripcion:          detalle.Descripcion,
		ValorUnit:            detalle.ValorUnit,
		IvaUnit:              detalle.IvaUnit,
		Cantidad:             detalle.Cantidad,
		Total:                detalle.Total,
		Activo:               detalle.Activo,
	}

	if detalle.DocumentoComercial != nil {

		response.DocumentoComercial = &DocumentoComercialResponse{
			ID:              detalle.DocumentoComercial.ID,
			Tipo:            detalle.DocumentoComercial.Tipo,
			NumeroDocumento: detalle.DocumentoComercial.NumeroDocumento,
		}
	}

	return response
}