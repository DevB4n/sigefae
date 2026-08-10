package archivo_origen

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

	var existing db.ArchivoOrigen

	err := s.db.
		Where("nombre = ?", req.Nombre).
		First(&existing).Error

	if err == nil {
		return nil, errors.New("el archivo de origen ya existe")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	archivo := db.ArchivoOrigen{
		Nombre: req.Nombre,
		Activo: true,
	}

	if err := s.db.Create(&archivo).Error; err != nil {
		return nil, err
	}

	response := toResponse(archivo)

	return &response, nil
}

func (s *Service) List() ([]Response, error) {

	var archivos []db.ArchivoOrigen

	if err := s.db.
		Order("nombre ASC").
		Find(&archivos).Error; err != nil {

		return nil, err
	}

	response := make([]Response, 0, len(archivos))

	for _, archivo := range archivos {
		response = append(response, toResponse(archivo))
	}

	return response, nil
}

func (s *Service) UpdateStatus(id uint, activo bool) error {

	var archivo db.ArchivoOrigen

	err := s.db.First(&archivo, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("archivo de origen no encontrado")
	}

	if err != nil {
		return err
	}

	archivo.Activo = activo

	return s.db.Save(&archivo).Error
}

func (s *Service) Update(id uint, req UpdateRequest) (*Response, error) {

	var archivo db.ArchivoOrigen

	err := s.db.First(&archivo, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("archivo de origen no encontrado")
	}

	if err != nil {
		return nil, err
	}

	var existing db.ArchivoOrigen

	err = s.db.
		Where("nombre = ? AND id <> ?", req.Nombre, id).
		First(&existing).Error

	if err == nil {
		return nil, errors.New("el archivo de origen ya existe")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if err := s.db.Model(&archivo).Updates(map[string]any{
		"nombre": req.Nombre,
	}).Error; err != nil {

		return nil, err
	}

	if err := s.db.First(&archivo, archivo.ID).Error; err != nil {
		return nil, err
	}

	response := toResponse(archivo)

	return &response, nil
}
