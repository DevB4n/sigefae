package trazabilidad

import (
	"time"

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
func (s *Service) Create(req CreateDTO) (*Response, error) {

	trazabilidad := db.Trazabilidad{
		DocumentoRadicadoID: req.DocumentoRadicadoID,
		UsuarioID:           req.UsuarioID,
		Accion:              req.Accion,
		Descripcion:         req.Descripcion,
		Fecha:               time.Now(),
	}

	if err := s.db.Create(&trazabilidad).Error; err != nil {
		return nil, err
	}

	response := toResponse(trazabilidad)

	return &response, nil
}
func (s *Service) List() ([]Response, error) {

	var trazabilidades []db.Trazabilidad

	if err := s.db.
		Order("fecha DESC").
		Find(&trazabilidades).Error; err != nil {

		return nil, err
	}


	response := make([]Response, 0, len(trazabilidades))

	for _, trazabilidad := range trazabilidades {

		response = append(
			response,
			toResponse(trazabilidad),
		)
	}


	return response, nil
}