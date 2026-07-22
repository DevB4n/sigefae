package db

import "time"

// ---------------------------------------------------------------------------
// Correo
// ---------------------------------------------------------------------------

type Correo struct {
	ID                    uint                 `gorm:"primaryKey;column:id" json:"id"`
	Asunto                string               `gorm:"column:asunto;type:varchar(255)" json:"asunto"`
	De                    string               `gorm:"column:de;type:varchar(255);index:idx_correo_de" json:"de"`
	Para                  string               `gorm:"column:para;type:varchar(255);index:idx_correo_para" json:"para"`
	FechaRecepcion        time.Time            `gorm:"column:fecha_recepcion;index:idx_correo_fecha" json:"fecha_recepcion"`
	IDMensaje             string               `gorm:"column:id_mensaje;type:varchar(255);uniqueIndex:uk_correo_id_mensaje" json:"id_mensaje"`
	Cuerpo                string               `gorm:"column:cuerpo;type:text" json:"cuerpo"`
	Cc                    string               `gorm:"column:cc;type:varchar(255)" json:"cc"`
	Bcc                   string               `gorm:"column:bcc;type:varchar(255)" json:"bcc"`
	ReplyTo               string               `gorm:"column:reply_to;type:varchar(255)" json:"reply_to"`
	IDEstado              uint                 `gorm:"column:id_estado;index:idx_correo_estado" json:"id_estado"`
	EstadoCorreo          *EstadoCorreo        `gorm:"foreignKey:IDEstado;references:ID" json:"estado_correo,omitempty"`
	DocumentosComerciales []DocumentoComercial `gorm:"foreignKey:CorreoID" json:"documentos_comerciales,omitempty"`
	Activo                bool                 `gorm:"column:activo;default:true" json:"activo"`
}

func (Correo) TableName() string { return "correo" }
