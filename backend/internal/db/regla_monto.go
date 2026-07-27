package db

import "time"

// ReglaMontoRuta define condiciones para agregar pasos adicionales a un flujo
// cuando una factura supere cierto monto. Se puede configurar a nivel de Área o Global.
type ReglaMontoRuta struct {
	ID                 uint      `gorm:"primaryKey;column:id" json:"id"`
	AreaID             *uint     `gorm:"column:area_id;index:idx_regla_area" json:"area_id,omitempty"` // Si es null, aplica globalmente
	Area               *Area     `gorm:"foreignKey:AreaID;references:ID" json:"area,omitempty"`
	RutaID             *uint     `gorm:"column:ruta_id;index:idx_regla_ruta" json:"ruta_id,omitempty"` // Para aplicar a una ruta específica, o global
	Ruta               *Ruta     `gorm:"foreignKey:RutaID;references:ID" json:"ruta,omitempty"`
	MontoMinimo        float64   `gorm:"column:monto_minimo" json:"monto_minimo"`
	MonedaID           uint      `gorm:"column:moneda_id" json:"moneda_id"`
	Moneda             *Moneda   `gorm:"foreignKey:MonedaID;references:ID" json:"moneda,omitempty"`
	UsuarioAprobadorID *uint     `gorm:"column:usuario_aprobador_id;index:idx_regla_aprobador" json:"usuario_aprobador_id,omitempty"`
	UsuarioAprobador   *Usuario  `gorm:"foreignKey:UsuarioAprobadorID;references:ID" json:"usuario_aprobador,omitempty"`
	RolAprobadorID     *uint     `gorm:"column:rol_aprobador_id;index:idx_regla_rol" json:"rol_aprobador_id,omitempty"`
	RolAprobador       *Rol      `gorm:"foreignKey:RolAprobadorID;references:ID" json:"rol_aprobador,omitempty"`
	Activo             bool      `gorm:"column:activo;default:true" json:"activo"`
	CreatedAt          time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt          time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (ReglaMontoRuta) TableName() string { return "regla_monto_ruta" }
