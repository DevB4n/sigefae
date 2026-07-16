package db

import "time"

type RegistroAprobacion struct {
	ID                  uint               `gorm:"primaryKey;column:id" json:"id"`
	DocumentoRadicadoID uint               `gorm:"column:documento_radicado_id;index:idx_aprobacion_documento" json:"documento_radicado_id"`
	DocumentoRadicado   *DocumentoRadicado `gorm:"foreignKey:DocumentoRadicadoID;references:ID" json:"documento_radicado,omitempty"`
	ResponsableID       uint               `gorm:"column:responsable_id;index:idx_aprobacion_responsable" json:"responsable_id"`
	Responsable         *Usuario           `gorm:"foreignKey:ResponsableID;references:ID" json:"responsable,omitempty"`
	RolID               uint               `gorm:"column:rol_id;index:idx_aprobacion_rol" json:"rol_id"`
	Rol                 *Rol               `gorm:"foreignKey:RolID;references:ID" json:"rol,omitempty"`
	Estado              string             `gorm:"column:estado;type:varchar(50)" json:"estado"` // "APROBADO", "RECHAZADO", "DEVUELTO"
	Observacion         string             `gorm:"column:observacion;type:text" json:"observacion"`
	Fecha               time.Time          `gorm:"column:fecha" json:"fecha"`
}

func (RegistroAprobacion) TableName() string { return "registro_aprobacion" }
