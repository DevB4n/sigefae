package db

import "time"

// ---------------------------------------------------------------------------
// Archivo
// ---------------------------------------------------------------------------

// Archivo representa un archivo adjunto a una factura.
// NOTA: en el diagrama "origen_id" aparece con tipo VARCHAR, pero por la
// relación con archivo_origen (PK id INTEGER) se asume que en realidad es un
// INTEGER FK. Ver aclaración al final de la respuesta.
type Archivo struct {
	ID                  uint               `gorm:"primaryKey;column:id" json:"id"`
	DocumentoRadicadoID uint               `gorm:"column:documento_radicado_id;index:idx_archivo_radicado" json:"documento_radicado_id"`
	DocumentoRadicado   *DocumentoRadicado `gorm:"foreignKey:DocumentoRadicadoID;references:ID" json:"documento_radicado,omitempty"`
	Nombre              string             `gorm:"column:nombre;type:varchar(255)" json:"nombre"`
	Extension           string             `gorm:"column:extension;type:varchar(20)" json:"extension"`
	Peso                int64              `gorm:"column:peso" json:"peso"`
	Ruta                string             `gorm:"column:ruta;type:varchar(500)" json:"ruta"`
	OrigenID            uint               `gorm:"column:origen_id;index:idx_archivo_origen" json:"origen_id"`
	Origen              *ArchivoOrigen     `gorm:"foreignKey:OrigenID;references:ID" json:"origen,omitempty"`
	CreadoPorID         uint               `gorm:"column:creado_por_id;index:idx_archivo_creado_por" json:"creado_por_id"`
	CreadoPor           *Usuario           `gorm:"foreignKey:CreadoPorID;references:ID" json:"creado_por,omitempty"`
	CreatedAt           time.Time          `gorm:"column:created_at" json:"created_at"`
	UpdatedAt           time.Time          `gorm:"column:updated_at" json:"updated_at"`
}

func (Archivo) TableName() string { return "archivo" }
