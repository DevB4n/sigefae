package receptor

type CreateRequest struct {
	Nombre          string `json:"nombre" binding:"required"`
	TipoDocumentoID uint   `json:"tipo_documento_id" binding:"required"`
	NumeroDocumento string `json:"numero_documento" binding:"required"`
}

type UpdateRequest struct {
	Nombre          string `json:"nombre" binding:"required"`
	TipoDocumentoID uint   `json:"tipo_documento_id" binding:"required"`
	NumeroDocumento string `json:"numero_documento" binding:"required"`
}

type UpdateStatusRequest struct {
	Activo bool `json:"activo"`
}
