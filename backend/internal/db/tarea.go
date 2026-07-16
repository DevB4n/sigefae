package db

import "time"

type Tarea struct {
	ID                  uint               `gorm:"primaryKey;column:id" json:"id"`
	DocumentoRadicadoID uint               `gorm:"column:documento_radicado_id;index:idx_tarea_documento" json:"documento_radicado_id"`
	DocumentoRadicado   *DocumentoRadicado `gorm:"foreignKey:DocumentoRadicadoID;references:ID" json:"documento_radicado,omitempty"`
	UsuarioAsignadoID   uint               `gorm:"column:usuario_asignado_id;index:idx_tarea_usuario" json:"usuario_asignado_id"`
	UsuarioAsignado     *Usuario           `gorm:"foreignKey:UsuarioAsignadoID;references:ID" json:"usuario_asignado,omitempty"`
	EstadoID            uint               `gorm:"column:estado_id;index:idx_tarea_estado" json:"estado_id"`
	Estado              *EstadoTarea       `gorm:"foreignKey:EstadoID;references:ID" json:"estado,omitempty"`
	Descripcion         string             `gorm:"column:descripcion;type:varchar(255)" json:"descripcion"` // Instrucción o descripción de la actividad
	FechaAsignacion     time.Time          `gorm:"column:fecha_asignacion" json:"fecha_asignacion"`
	FechaInicio         *time.Time         `gorm:"column:fecha_inicio" json:"fecha_inicio"`
	FechaLimite         *time.Time         `gorm:"column:fecha_limite;index:idx_tarea_limite" json:"fecha_limite"`
	FechaFinalizacion   *time.Time         `gorm:"column:fecha_finalizacion" json:"fecha_finalizacion"`
	CreatedAt           time.Time          `gorm:"column:created_at" json:"created_at"`
}

func (Tarea) TableName() string { return "tarea" }
