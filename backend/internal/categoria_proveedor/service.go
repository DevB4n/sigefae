package categoria_proveedor

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

	var existing db.CategoriaProveedor

	err := s.db.
		Where("nombre = ?", req.Nombre).
		First(&existing).Error

	if err == nil {
		return nil, errors.New("la categoría de proveedor ya existe")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	categoria := db.CategoriaProveedor{
		Nombre:      req.Nombre,
		Descripcion: req.Descripcion,
		Activo:      true,
	}

	if err := s.db.Create(&categoria).Error; err != nil {
		return nil, err
	}

	response := toResponse(categoria)

	return &response, nil
}

func (s *Service) List() ([]Response, error) {

	var categorias []db.CategoriaProveedor

	if err := s.db.
		Order("nombre ASC").
		Find(&categorias).Error; err != nil {

		return nil, err
	}

	response := make([]Response, 0, len(categorias))

	for _, categoria := range categorias {

		response = append(
			response,
			toResponse(categoria),
		)
	}

	return response, nil
}

func (s *Service) UpdateStatus(id uint, activo bool) error {

	var categoria db.CategoriaProveedor

	err := s.db.First(&categoria, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("categoría de proveedor no encontrada")
	}

	if err != nil {
		return err
	}

	categoria.Activo = activo

	return s.db.Save(&categoria).Error
}

func (s *Service) Update(id uint, req UpdateRequest) (*Response, error) {

	var categoria db.CategoriaProveedor

	err := s.db.First(&categoria, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("categoría de proveedor no encontrada")
	}

	if err != nil {
		return nil, err
	}

	var existing db.CategoriaProveedor

	err = s.db.
		Where("nombre = ? AND id <> ?", req.Nombre, id).
		First(&existing).Error

	if err == nil {
		return nil, errors.New("la categoría de proveedor ya existe")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if err := s.db.Model(&categoria).Updates(map[string]any{
		"nombre":      req.Nombre,
		"descripcion": req.Descripcion,
	}).Error; err != nil {

		return nil, err
	}

	if err := s.db.First(&categoria, categoria.ID).Error; err != nil {
		return nil, err
	}

	response := toResponse(categoria)

	return &response, nil
}
