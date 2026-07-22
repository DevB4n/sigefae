package db

// ---------------------------------------------------------------------------
// Geografía
// ---------------------------------------------------------------------------

type Departamento struct {
	ID     uint   `gorm:"primaryKey;column:id" json:"id"`
	Nombre string `gorm:"column:nombre;type:varchar(255)" json:"nombre"`
	PaisID uint   `gorm:"column:pais_id;index:idx_departamento_pais_id" json:"pais_id"`
	Pais   *Pais  `gorm:"foreignKey:PaisID;references:ID" json:"pais,omitempty"`
	Activo bool   `gorm:"column:activo;default:true" json:"activo"`
}

func (Departamento) TableName() string { return "departamento" }

type Municipio struct {
	ID             uint          `gorm:"primaryKey;column:id" json:"id"`
	Nombre         string        `gorm:"column:nombre;type:varchar(255)" json:"nombre"`
	DepartamentoID uint          `gorm:"column:departamento_id;index:idx_municipio_departamento_id" json:"departamento_id"`
	Departamento   *Departamento `gorm:"foreignKey:DepartamentoID;references:ID" json:"departamento,omitempty"`
	Activo         bool          `gorm:"column:activo;default:true" json:"activo"`
}

func (Municipio) TableName() string { return "municipio" }

type Direccion struct {
	ID           uint       `gorm:"primaryKey;column:id" json:"id"`
	Nombre       string     `gorm:"column:nombre;type:varchar(255)" json:"nombre"`
	Linea1       string     `gorm:"column:linea_1;type:varchar(255)" json:"linea_1"`
	Linea2       string     `gorm:"column:linea_2;type:varchar(255)" json:"linea_2"`
	CodigoPostal string     `gorm:"column:codigo_postal;type:varchar(20)" json:"codigo_postal"`
	IDMunicipio  uint       `gorm:"column:id_municipio;index:idx_direccion_municipio_id" json:"id_municipio"`
	Municipio    *Municipio `gorm:"foreignKey:IDMunicipio;references:ID" json:"municipio,omitempty"`
	Activo       bool       `gorm:"column:activo;default:true" json:"activo"`
}

func (Direccion) TableName() string { return "direccion" }
