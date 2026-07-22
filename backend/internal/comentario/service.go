package comentario

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

	comentario := db.Comentario{
		DocumentoRadicadoID: req.DocumentoRadicadoID,
		UsuarioID:           req.UsuarioID,
		Descripcion:         req.Descripcion,
		Fecha:               time.Now(),
	}


	if err := s.db.Create(&comentario).Error; err != nil {
		return nil, err
	}


	response := toResponse(comentario)

	return &response, nil
}
func (s *Service) List() ([]Response, error) {

	var comentarios []db.Comentario


	if err := s.db.
		Order("fecha DESC").
		Find(&comentarios).Error; err != nil {

		return nil, err
	}


	response := make([]Response, 0, len(comentarios))


	for _, comentario := range comentarios {

		response = append(
			response,
			toResponse(comentario),
		)
	}


	return response, nil
}