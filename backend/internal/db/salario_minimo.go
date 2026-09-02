package db

import "time"

// SalarioMinimo define el catálogo de Salario Mínimo Mensual Legal Vigente (SMMLV) por año.
type SalarioMinimo struct {
	ID        uint      `gorm:"primaryKey;column:id" json:"id"`
	Ano       int       `gorm:"column:ano;unique" json:"ano"`
	Valor     float64   `gorm:"column:valor" json:"valor"`
	Activo    bool      `gorm:"column:activo;default:true" json:"activo"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (SalarioMinimo) TableName() string { return "salario_minimo" }
