package detalle_documento_comercial

type CreateDTO struct {
	DocumentoComercialID uint    `json:"documento_comercial_id" binding:"required"`
	Descripcion          string  `json:"descripcion" binding:"required"`
	ValorUnit            float64 `json:"valor_unitario"`
	IvaUnit              float64 `json:"iva_unitario"`
	Cantidad             float64 `json:"cantidad"`
	Total                float64 `json:"total"`
	Activo               bool    `json:"activo"`
}

type UpdateDTO struct {
	DocumentoComercialID uint    `json:"documento_comercial_id" binding:"required"`
	Descripcion          string  `json:"descripcion" binding:"required"`
	ValorUnit            float64 `json:"valor_unitario"`
	IvaUnit              float64 `json:"iva_unitario"`
	Cantidad             float64 `json:"cantidad"`
	Total                float64 `json:"total"`
}

type UpdateStatusDTO struct {
	Activo bool `json:"activo"`
}
