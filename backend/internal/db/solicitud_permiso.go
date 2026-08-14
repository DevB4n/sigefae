package db

import "time"

type SolicitudPermiso struct {
	ID              uint       `gorm:"primaryKey;column:id" json:"id"`
	Tipo            string     `gorm:"column:tipo;type:varchar(50);not null" json:"tipo"` // norma_reparto, archivo
	ObjetoID        uint       `gorm:"column:objeto_id;not null" json:"objeto_id"`
	Accion          string     `gorm:"column:accion;type:varchar(30);not null" json:"accion"` // editar, eliminar
	SolicitanteID   uint       `gorm:"column:solicitante_id;not null;index:idx_solicitud_permiso_solicitante" json:"solicitante_id"`
	Solicitante     *Usuario   `gorm:"foreignKey:SolicitanteID;references:ID" json:"solicitante,omitempty"`
	PropietarioID   uint       `gorm:"column:propietario_id;not null;index:idx_solicitud_permiso_propietario" json:"propietario_id"`
	Propietario     *Usuario   `gorm:"foreignKey:PropietarioID;references:ID" json:"propietario,omitempty"`
	Mensaje         string     `gorm:"column:mensaje;type:text" json:"mensaje"`
	Estado          string     `gorm:"column:estado;type:varchar(50);default:'Pendiente'" json:"estado"` // Pendiente, Aprobada, Rechazada
	ResueltoPorID   *uint      `gorm:"column:resuelto_por_id" json:"resuelto_por_id,omitempty"`
	ResueltoPor     *Usuario   `gorm:"foreignKey:ResueltoPorID;references:ID" json:"resuelto_por,omitempty"`
	Respuesta       string     `gorm:"column:respuesta;type:text" json:"respuesta,omitempty"`
	FechaCreacion   time.Time  `gorm:"column:fecha_creacion" json:"fecha_creacion"`
	FechaResolucion *time.Time `gorm:"column:fecha_resolucion" json:"fecha_resolucion,omitempty"`
}

func (SolicitudPermiso) TableName() string { return "solicitud_permiso" }
