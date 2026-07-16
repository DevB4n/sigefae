package db

import "time"

type Notificacion struct {
	ID                  uint               `gorm:"primaryKey;column:id" json:"id"`
	UsuarioID           uint               `gorm:"column:usuario_id;index:idx_notificacion_usuario" json:"usuario_id"`
	Usuario             *Usuario           `gorm:"foreignKey:UsuarioID;references:ID" json:"usuario,omitempty"`
	DocumentoRadicadoID *uint              `gorm:"column:documento_radicado_id;index:idx_notificacion_doc" json:"documento_radicado_id,omitempty"`
	DocumentoRadicado   *DocumentoRadicado `gorm:"foreignKey:DocumentoRadicadoID;references:ID" json:"documento_radicado,omitempty"`
	Mensaje             string             `gorm:"column:mensaje;type:text" json:"mensaje"`
	Estado              string             `gorm:"column:estado;type:varchar(50)" json:"estado"` // "Pendiente", "Enviada", "Leida", "Expirada"
	Tipo                string             `gorm:"column:tipo;type:varchar(50)" json:"tipo"` // "Recordatorio", "Asignacion", "Sistema"
	FechaCreacion       time.Time          `gorm:"column:fecha_creacion" json:"fecha_creacion"`
	FechaEnvio          *time.Time         `gorm:"column:fecha_envio" json:"fecha_envio,omitempty"`
	FechaLectura        *time.Time         `gorm:"column:fecha_lectura" json:"fecha_lectura,omitempty"`
}

func (Notificacion) TableName() string { return "notificacion" }
