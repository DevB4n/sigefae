package db

import "time"

type SolicitudRechazo struct {
	ID                  uint               `gorm:"primaryKey;column:id" json:"id"`
	DocumentoRadicadoID uint               `gorm:"column:documento_radicado_id;index:idx_solicitud_doc" json:"documento_radicado_id"`
	DocumentoRadicado   *DocumentoRadicado `gorm:"foreignKey:DocumentoRadicadoID;references:ID" json:"documento_radicado,omitempty"`
	UsuarioID           uint               `gorm:"column:usuario_id;index:idx_solicitud_usuario" json:"usuario_id"`
	Usuario             *Usuario           `gorm:"foreignKey:UsuarioID;references:ID" json:"usuario,omitempty"`
	Mensaje             string             `gorm:"column:mensaje;type:text" json:"mensaje"`
	Estado              string             `gorm:"column:estado;type:varchar(50)" json:"estado"` // Pendiente, Aceptada, Rechazada
	ResueltoPorID       *uint              `gorm:"column:resuelto_por_id" json:"resuelto_por_id,omitempty"`
	ResueltoPor         *Usuario           `gorm:"foreignKey:ResueltoPorID;references:ID" json:"resuelto_por,omitempty"`
	Respuesta           string             `gorm:"column:respuesta;type:text" json:"respuesta,omitempty"`
	FechaCreacion       time.Time          `gorm:"column:fecha_creacion" json:"fecha_creacion"`
	FechaResolucion     *time.Time         `gorm:"column:fecha_resolucion" json:"fecha_resolucion,omitempty"`
}

func (SolicitudRechazo) TableName() string { return "solicitud_rechazo" }
