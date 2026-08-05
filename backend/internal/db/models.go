package db
// internal/db/migrate.go

// Fase 1: Tablas base (sin dependencias circulares)
func BaseModels() []interface{} {
	return []interface{}{
		&Pais{}, &Departamento{}, &Municipio{}, &Direccion{},
		&TipoDocumento{}, &TipoPersona{}, &CategoriaProveedor{},
		&ActividadEconomica{}, &Proveedor{}, &ResponsabilidadFiscal{},
		&Contacto{}, &Rol{}, &Usuario{}, &Ruta{}, &PasoRuta{},
		&Moneda{}, &Area{}, &Receptor{}, &EstadoCorreo{},
		&Correo{}, &ArchivoOrigen{}, &Archivo{}, &DocumentoComercial{},
		&DetalleDocumentoComercial{}, &TipoFactura{}, &TipoPago{},
		&MetodoPago{}, &TipoRadicacion{}, &CodigoQr{},
		&EstadoDocumentoRadicado{}, &EstadoTarea{},
		&NormaReparto{}, 
	}
}

// Fase 2: Tablas con dependencias circulares
func CircularModels() []interface{} {
	return []interface{}{
		&DocumentoRadicado{},
		&Tarea{},
		&Comentario{}, &Trazabilidad{}, &RegistroAprobacion{},
		&Notificacion{}, &ReglaMontoRuta{},
		&RadicadoNormaReparto{},
	}
}