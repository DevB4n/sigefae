package salario_minimo

import (
	"gorm.io/gorm"
	"sigefae/internal/db"
)

type Service struct {
	db *gorm.DB
}

func New(database *gorm.DB) *Service {
	return &Service{db: database}
}

func (s *Service) List() ([]db.SalarioMinimo, error) {
	var items []db.SalarioMinimo
	err := s.db.Order("ano desc").Find(&items).Error
	return items, err
}

func (s *Service) Create(item *db.SalarioMinimo) error {
	// Desactivar el actual si es que se quiere que solo haya uno activo
	// O simplemente lo creamos
	return s.db.Create(item).Error
}

func (s *Service) Update(id uint, item *db.SalarioMinimo) error {
	return s.db.Model(&db.SalarioMinimo{}).Where("id = ?", id).Updates(item).Error
}

func (s *Service) UpdateStatus(id uint, activo bool) error {
	return s.db.Model(&db.SalarioMinimo{}).Where("id = ?", id).Update("activo", activo).Error
}

func (s *Service) Delete(id uint) error {
	return s.db.Model(&db.SalarioMinimo{}).Where("id = ?", id).Delete(&db.SalarioMinimo{}).Error
}
