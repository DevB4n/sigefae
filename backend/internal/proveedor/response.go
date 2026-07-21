package proveedor

import "sigefae/internal/db"

type Response struct {
	ID uint `json:"id"`

	TipoDocumentoID uint `json:"tipo_documento_id"`
	TipoDocumento   *db.TipoDocumento `json:"tipo_documento,omitempty"`

	NumeroDocumento string `json:"numero_documento"`

	CategoriaID uint `json:"categoria_id"`
	Categoria   *db.CategoriaProveedor `json:"categoria,omitempty"`

	RutaPredeterminadaID *uint    `json:"ruta_predeterminada_id,omitempty"`
	RutaPredeterminada   *db.Ruta `json:"ruta_predeterminada,omitempty"`

	RazonSocial     string `json:"razon_social"`
	NombreComercial string `json:"nombre_comercial"`

	TipoPersonaID uint `json:"tipo_persona_id"`
	TipoPersona   *db.TipoPersona `json:"tipo_persona,omitempty"`

	ActividadEconomicaID uint `json:"actividad_economica_id"`
	ActividadEconomica   *db.ActividadEconomica `json:"actividad_economica,omitempty"`

	DireccionID uint `json:"direccion_id"`
	Direccion   *db.Direccion `json:"direccion,omitempty"`

	Activo bool `json:"activo"`
}

func toResponse(proveedor db.Proveedor) Response {

	return Response{
		ID: proveedor.ID,

		TipoDocumentoID: proveedor.TipoDocumentoID,
		TipoDocumento:   proveedor.TipoDocumento,

		NumeroDocumento: proveedor.NumeroDocumento,

		CategoriaID: proveedor.CategoriaID,
		Categoria:   proveedor.Categoria,

		RutaPredeterminadaID: proveedor.RutaPredeterminadaID,
		RutaPredeterminada:   proveedor.RutaPredeterminada,

		RazonSocial:     proveedor.RazonSocial,
		NombreComercial: proveedor.NombreComercial,

		TipoPersonaID: proveedor.TipoPersonaID,
		TipoPersona:   proveedor.TipoPersona,

		ActividadEconomicaID: proveedor.ActividadEconomicaID,
		ActividadEconomica:   proveedor.ActividadEconomica,

		DireccionID: proveedor.DireccionID,
		Direccion:   proveedor.Direccion,

		Activo: proveedor.Activo,
	}
}