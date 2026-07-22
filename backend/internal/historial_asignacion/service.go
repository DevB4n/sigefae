package historial_asignacion

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

	now := time.Now()


	// Cerrar asignación anterior si existe

	s.db.
		Model(&db.HistorialAsignacion{}).
		Where(
			"documento_radicado_id = ? AND hasta IS NULL",
			req.DocumentoRadicadoID,
		).
		Update(
			"hasta",
			now,
		)


	historial := db.HistorialAsignacion{
		DocumentoRadicadoID: req.DocumentoRadicadoID,
		UsuarioID:           req.UsuarioID,
		Desde:               now,
	}


	if err := s.db.Create(&historial).Error; err != nil {
		return nil, err
	}


	response := toResponse(historial)

	return &response, nil
}
func (s *Service) List() ([]Response, error) {

	var historial []db.HistorialAsignacion


	if err := s.db.
		Order("desde DESC").
		Find(&historial).Error; err != nil {

		return nil, err
	}


	response := make([]Response, 0, len(historial))


	for _, item := range historial {

		response = append(
			response,
			toResponse(item),
		)
	}


	return response, nil
}