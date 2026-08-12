package db

import "time"

type SolicitudCambioNormaReparto struct {
	ID                     uint                  `gorm:"primaryKey;column:id" json:"id"`
	DocumentoRadicadoID    uint                  `gorm:"column:documento_radicado_id;not null;index:idx_sol_cambio_doc" json:"documento_radicado_id"`
	DocumentoRadicado      *DocumentoRadicado    `gorm:"foreignKey:DocumentoRadicadoID;references:ID" json:"documento_radicado,omitempty"`
	RadicadoNormaRepartoID uint                  `gorm:"column:radicado_norma_reparto_id;not null;index:idx_sol_cambio_rnr" json:"radicado_norma_reparto_id"`
	RadicadoNormaReparto   *RadicadoNormaReparto `gorm:"foreignKey:RadicadoNormaRepartoID;references:ID" json:"radicado_norma_reparto,omitempty"`
	NormaRepartoID         uint                  `gorm:"column:norma_reparto_id;not null" json:"norma_reparto_id"`
	NormaReparto           *NormaReparto         `gorm:"foreignKey:NormaRepartoID;references:ID" json:"norma_reparto,omitempty"`
	NuevoPorcentaje        float64               `gorm:"column:nuevo_porcentaje;type:decimal(5,2);not null" json:"nuevo_porcentaje"`
	PorcentajeAnterior     float64               `gorm:"column:porcentaje_anterior;type:decimal(5,2);not null" json:"porcentaje_anterior"`
	UsuarioID              uint                  `gorm:"column:usuario_id;not null;index:idx_sol_cambio_usr" json:"usuario_id"`
	Usuario                *Usuario              `gorm:"foreignKey:UsuarioID;references:ID" json:"usuario,omitempty"`
	Justificacion          string                `gorm:"column:justificacion;type:text" json:"justificacion"`
	Estado                 string                `gorm:"column:estado;type:varchar(50);default:'Pendiente'" json:"estado"` // Pendiente, Aprobada, Rechazada
	ResueltoPorID          *uint                 `gorm:"column:resuelto_por_id" json:"resuelto_por_id,omitempty"`
	ResueltoPor            *Usuario              `gorm:"foreignKey:ResueltoPorID;references:ID" json:"resuelto_por,omitempty"`
	Respuesta              string                `gorm:"column:respuesta;type:text" json:"respuesta,omitempty"`
	FechaCreacion          time.Time             `gorm:"column:fecha_creacion" json:"fecha_creacion"`
	FechaResolucion        *time.Time            `gorm:"column:fecha_resolucion" json:"fecha_resolucion,omitempty"`
}

func (SolicitudCambioNormaReparto) TableName() string { return "solicitud_cambio_norma_reparto" }
