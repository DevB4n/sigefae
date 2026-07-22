package proveedor

type CreateRequest struct {
	TipoDocumentoID      uint   `json:"tipo_documento_id" binding:"required"`
	NumeroDocumento      string `json:"numero_documento" binding:"required"`
	CategoriaID          uint   `json:"categoria_id" binding:"required"`
	RutaPredeterminadaID *uint  `json:"ruta_predeterminada_id"`
	RazonSocial          string `json:"razon_social" binding:"required"`
	NombreComercial      string `json:"nombre_comercial"`
	TipoPersonaID        uint   `json:"tipo_persona_id" binding:"required"`
	ActividadEconomicaID uint   `json:"actividad_economica_id" binding:"required"`
	DireccionID          uint   `json:"direccion_id" binding:"required"`
}

type UpdateRequest struct {
	TipoDocumentoID      uint   `json:"tipo_documento_id" binding:"required"`
	NumeroDocumento      string `json:"numero_documento" binding:"required"`
	CategoriaID          uint   `json:"categoria_id" binding:"required"`
	RutaPredeterminadaID *uint  `json:"ruta_predeterminada_id"`
	RazonSocial          string `json:"razon_social" binding:"required"`
	NombreComercial      string `json:"nombre_comercial"`
	TipoPersonaID        uint   `json:"tipo_persona_id" binding:"required"`
	ActividadEconomicaID uint   `json:"actividad_economica_id" binding:"required"`
	DireccionID          uint   `json:"direccion_id" binding:"required"`
}

type UpdateStatusRequest struct {
	Activo bool `json:"activo"`
}
