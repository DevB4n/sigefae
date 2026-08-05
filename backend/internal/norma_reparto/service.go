package norma_reparto

import (
	"errors"
	"sigefae/internal/db"

	"gorm.io/gorm"
)

type Service struct {
	db *gorm.DB
}

func New(database *gorm.DB) *Service {
	return &Service{db: database}
}

func (s *Service) Create(dto CreateDTO) (*db.NormaReparto, error) {
	norma := db.NormaReparto{
		Codigo:       dto.Codigo,
		Nombre:       dto.Nombre,
		Sucursal:     dto.Sucursal,
		Departamento: dto.Departamento,
		Tipo:         dto.Tipo,
		TarifaIva:    dto.TarifaIva,
		Activo:       true,
	}
	if err := s.db.Create(&norma).Error; err != nil {
		return nil, err
	}
	return &norma, nil
}

func (s *Service) List(activo *bool) ([]db.NormaReparto, error) {
	var normas []db.NormaReparto
	query := s.db.Order("codigo asc")
	if activo != nil {
		query = query.Where("activo = ?", *activo)
	}
	if err := query.Find(&normas).Error; err != nil {
		return nil, err
	}
	return normas, nil
}

func (s *Service) GetByID(id uint) (*db.NormaReparto, error) {
	var norma db.NormaReparto
	if err := s.db.First(&norma, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("norma de reparto no encontrada")
		}
		return nil, err
	}
	return &norma, nil
}

func (s *Service) Update(id uint, dto UpdateDTO) (*db.NormaReparto, error) {
	var norma db.NormaReparto
	if err := s.db.First(&norma, id).Error; err != nil {
		return nil, err
	}
	updates := map[string]interface{}{}
	if dto.Codigo != "" {
		updates["codigo"] = dto.Codigo
	}
	if dto.Nombre != "" {
		updates["nombre"] = dto.Nombre
	}
	if dto.Sucursal != "" {
		updates["sucursal"] = dto.Sucursal
	}
	if dto.Departamento != "" {
		updates["departamento"] = dto.Departamento
	}
	if dto.Tipo != nil {
		updates["tipo"] = *dto.Tipo
	} else {
		updates["tipo"] = nil
	}
	if dto.TarifaIva != nil {
		updates["tarifa_iva"] = *dto.TarifaIva
	} else {
		updates["tarifa_iva"] = nil
	}

	if len(updates) > 0 {
		if err := s.db.Model(&norma).Updates(updates).Error; err != nil {
			return nil, err
		}
	}
	return &norma, nil
}

func (s *Service) UpdateStatus(id uint, dto UpdateStatusDTO) (*db.NormaReparto, error) {
	var norma db.NormaReparto
	if err := s.db.First(&norma, id).Error; err != nil {
		return nil, err
	}
	if err := s.db.Model(&norma).Update("activo", dto.Activo).Error; err != nil {
		return nil, err
	}
	return &norma, nil
}