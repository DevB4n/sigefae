package moneda

import (
	"errors"

	"gorm.io/gorm"

	"sigefae/internal/db"
)

type Service struct {
	db *gorm.DB
}

func New(database *gorm.DB) *Service {

	return &Service{
		db: database,
	}
}

func (s *Service) Create(req CreateRequest) (*Response, error) {

	// ==========================
	// Validar nombre repetido
	// ==========================

	var existing db.Moneda

	err := s.db.
		Where("nombre = ?", req.Nombre).
		First(&existing).Error

	if err == nil {
		return nil, errors.New("la moneda ya existe")
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// ==========================
	// Validar código repetido
	// ==========================

	err = s.db.
		Where("codigo = ?", req.Codigo).
		First(&existing).Error

	if err == nil {
		return nil, errors.New("el código ya existe")
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// ==========================
	// Crear
	// ==========================

	moneda := db.Moneda{
		Nombre: req.Nombre,
		Codigo: req.Codigo,
		Activo: true,
	}

	if err := s.db.Create(&moneda).Error; err != nil {
		return nil, err
	}

	response := toResponse(moneda)

	return &response, nil
}

func (s *Service) List() ([]Response, error) {

	var monedas []db.Moneda

	if err := s.db.
		Order("nombre ASC").
		Find(&monedas).Error; err != nil {

		return nil, err
	}

	response := make([]Response, 0, len(monedas))

	for _, moneda := range monedas {
		response = append(response, toResponse(moneda))
	}

	return response, nil
}

func (s *Service) UpdateStatus(id uint, activo bool) error {

	var moneda db.Moneda

	err := s.db.First(&moneda, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("moneda no encontrada")
	}

	if err != nil {
		return err
	}

	moneda.Activo = activo

	return s.db.Save(&moneda).Error
}

func (s *Service) Update(id uint, req UpdateRequest) (*Response, error) {

	var moneda db.Moneda

	err := s.db.First(&moneda, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("moneda no encontrada")
	}

	if err != nil {
		return nil, err
	}

	// ==========================
	// Validar nombre repetido
	// ==========================

	var existing db.Moneda

	err = s.db.
		Where("nombre = ? AND id <> ?", req.Nombre, id).
		First(&existing).Error

	if err == nil {
		return nil, errors.New("la moneda ya existe")
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// ==========================
	// Validar código repetido
	// ==========================

	err = s.db.
		Where("codigo = ? AND id <> ?", req.Codigo, id).
		First(&existing).Error

	if err == nil {
		return nil, errors.New("el código ya existe")
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// ==========================
	// Actualizar
	// ==========================

	if err := s.db.Model(&moneda).Updates(map[string]interface{}{
		"nombre": req.Nombre,
		"codigo": req.Codigo,
	}).Error; err != nil {

		return nil, err
	}

	if err := s.db.First(&moneda, moneda.ID).Error; err != nil {
		return nil, err
	}

	response := toResponse(moneda)

	return &response, nil
}