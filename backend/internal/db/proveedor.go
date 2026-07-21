package db

import "time"

// ---------------------------------------------------------------------------
// Proveedores y contactos
// ---------------------------------------------------------------------------

type Proveedor struct {
	ID                   uint                     `gorm:"primaryKey;column:id" json:"id"`
	TipoDocumentoID      uint                     `gorm:"column:tipo_documento_id;index:idx_proveedor_tipo_documento_id" json:"tipo_documento_id"`
	TipoDocumento        *TipoDocumento           `gorm:"foreignKey:TipoDocumentoID;references:ID" json:"tipo_documento,omitempty"`
	NumeroDocumento      string                   `gorm:"column:numero_documento;type:varchar(50);uniqueIndex:uk_proveedor_documento" json:"numero_documento"`
	CategoriaID          uint                     `gorm:"column:categoria_id;index:idx_proveedor_categoria_id" json:"categoria_id"`
	Categoria            *CategoriaProveedor      `gorm:"foreignKey:CategoriaID;references:ID" json:"categoria,omitempty"`
	RutaPredeterminadaID *uint                    `gorm:"column:ruta_predeterminada_id;index:idx_proveedor_ruta" json:"ruta_predeterminada_id,omitempty"`
	RutaPredeterminada   *Ruta                    `gorm:"foreignKey:RutaPredeterminadaID;references:ID" json:"ruta_predeterminada,omitempty"`
	RazonSocial          string                   `gorm:"column:razon_social;type:varchar(255)" json:"razon_social"`
	NombreComercial      string                   `gorm:"column:nombre_comercial;type:varchar(255)" json:"nombre_comercial"`
	TipoPersonaID        uint                     `gorm:"column:tipo_persona_id;index:idx_proveedor_tipo_persona_id" json:"tipo_persona_id"`
	TipoPersona          *TipoPersona             `gorm:"foreignKey:TipoPersonaID;references:ID" json:"tipo_persona,omitempty"`
	ActividadEconomicaID uint                     `gorm:"column:actividad_economica_id;index:idx_proveedor_actividad_id" json:"actividad_economica_id"`
	ActividadEconomica   *ActividadEconomica      `gorm:"foreignKey:ActividadEconomicaID;references:ID" json:"actividad_economica,omitempty"`
	DireccionID          uint                     `gorm:"column:direccion_id;index:idx_proveedor_direccion_id" json:"direccion_id"`
	Direccion            *Direccion               `gorm:"foreignKey:DireccionID;references:ID" json:"direccion,omitempty"`
	Contactos            []Contacto               `gorm:"foreignKey:ProveedorID" json:"contactos,omitempty"`
	Responsabilidades    []ResponsabilidadFiscal  `gorm:"foreignKey:IDProveedor" json:"responsabilidades,omitempty"`
	DocumentosComerciales []DocumentoComercial    `gorm:"foreignKey:IDProveedor" json:"documentos_comerciales,omitempty"`
	Activo 				 bool 					  `gorm:"column:activo;default:true" json:"activo"`
	CreatedAt            time.Time                `gorm:"column:created_at" json:"created_at"`
	UpdatedAt            time.Time                `gorm:"column:updated_at" json:"updated_at"`
}

func (Proveedor) TableName() string { return "proveedor" }

type ResponsabilidadFiscal struct {
	ID          uint       `gorm:"primaryKey;column:id" json:"id"`
	IDProveedor uint       `gorm:"column:id_proveedor;index:idx_responsabilidad_proveedor_id" json:"id_proveedor"`
	Proveedor   *Proveedor `gorm:"foreignKey:IDProveedor;references:ID" json:"proveedor,omitempty"`
	Codigo      string     `gorm:"column:codigo;type:varchar(50);index:idx_responsabilidad_codigo" json:"codigo"`
}

func (ResponsabilidadFiscal) TableName() string { return "responsabilidad_fiscal" }

type Contacto struct {
	ID          uint       `gorm:"primaryKey;column:id" json:"id"`
	ProveedorID uint       `gorm:"column:proveedor_id;index:idx_contacto_proveedor_id" json:"proveedor_id"`
	Proveedor   *Proveedor `gorm:"foreignKey:ProveedorID;references:ID" json:"proveedor,omitempty"`
	Nombre      string     `gorm:"column:nombre;type:varchar(255)" json:"nombre"`
	Cargo       string     `gorm:"column:cargo;type:varchar(255)" json:"cargo"`
	Telefono    string     `gorm:"column:telefono;type:varchar(50)" json:"telefono"`
	Correo      string     `gorm:"column:correo;type:varchar(255)" json:"correo"`
	Activo      bool       `gorm:"column:activo;default:true" json:"activo"`
}

func (Contacto) TableName() string { return "contacto" }
