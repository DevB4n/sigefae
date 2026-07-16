package db

// ConfiguracionSistema almacena parámetros globales de la aplicación.
// Se asume un registro por clave.
type ConfiguracionSistema struct {
	ID          uint   `gorm:"primaryKey;column:id" json:"id"`
	Clave       string `gorm:"column:clave;type:varchar(100);uniqueIndex:uk_config_clave" json:"clave"`
	Valor       string `gorm:"column:valor;type:varchar(255)" json:"valor"`
	Descripcion string `gorm:"column:descripcion;type:text" json:"descripcion"`
}

func (ConfiguracionSistema) TableName() string { return "configuracion_sistema" }
