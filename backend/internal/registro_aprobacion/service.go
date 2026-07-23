package registro_aprobacion

import (
	"errors"
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

	// ==========================
	// Validar Documento Radicado
	// ==========================

	var documento db.DocumentoRadicado

	err := s.db.First(
		&documento,
		req.DocumentoRadicadoID,
	).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("documento radicado no encontrado")
	}

	if err != nil {
		return nil, err
	}

	// ==========================
	// Validar Responsable
	// ==========================

	var usuario db.Usuario

	err = s.db.First(
		&usuario,
		req.ResponsableID,
	).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("responsable no encontrado")
	}

	if err != nil {
		return nil, err
	}

	// ==========================
	// Validar Rol
	// ==========================

	var rol db.Rol

	err = s.db.First(
		&rol,
		req.RolID,
	).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("rol no encontrado")
	}

	if err != nil {
		return nil, err
	}

	// ==========================
	// Crear
	// ==========================

	registro := db.RegistroAprobacion{
		DocumentoRadicadoID: req.DocumentoRadicadoID,
		ResponsableID:       req.ResponsableID,
		RolID:               req.RolID,
		Estado:              req.Estado,
		Observacion:         req.Observacion,
		Fecha:               time.Now(),
	}

	if err := s.db.Create(&registro).Error; err != nil {
		return nil, err
	}

	response := toResponse(registro)

	return &response, nil
}

func (s *Service) List() ([]Response, error) {

	var registros []db.RegistroAprobacion

	if err := s.db.
		Order("fecha DESC").
		Find(&registros).Error; err != nil {

		return nil, err
	}

	response := make([]Response, 0, len(registros))

	for _, registro := range registros {
		response = append(
			response,
			toResponse(registro),
		)
	}

	return response, nil
}