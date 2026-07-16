package db

import "time"

type HistorialAsignacion struct {
	ID                  uint               `gorm:"primaryKey;column:id" json:"id"`
	DocumentoRadicadoID uint               `gorm:"column:documento_radicado_id;index:idx_historial_doc" json:"documento_radicado_id"`
	DocumentoRadicado   *DocumentoRadicado `gorm:"foreignKey:DocumentoRadicadoID;references:ID" json:"documento_radicado,omitempty"`
	UsuarioID           uint               `gorm:"column:usuario_id;index:idx_historial_usuario" json:"usuario_id"`
	Usuario             *Usuario           `gorm:"foreignKey:UsuarioID;references:ID" json:"usuario,omitempty"`
	Desde               time.Time          `gorm:"column:desde;index:idx_historial_desde" json:"desde"`
	Hasta               *time.Time         `gorm:"column:hasta;index:idx_historial_hasta" json:"hasta,omitempty"`
}

func (HistorialAsignacion) TableName() string { return "historial_asignacion" }
