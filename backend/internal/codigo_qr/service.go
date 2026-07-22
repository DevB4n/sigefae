package codigo_qr

import(
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
func (s *Service) Create(dto CreateDTO) (*Response, error) {

	// ==========================
	// Validar URL única
	// ==========================

	var existing db.CodigoQr

	err := s.db.
		Where("url = ?", dto.Url).
		First(&existing).Error

	if err == nil {
		return nil, errors.New("la URL ya existe")
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// ==========================
	// Crear
	// ==========================

	codigo := db.CodigoQr{
		Url:    dto.Url,
		Activo: true,
	}

	if err := s.db.Create(&codigo).Error; err != nil {
		return nil, err
	}

	if err := s.db.First(&codigo, codigo.ID).Error; err != nil {
		return nil, err
	}

	response := toResponse(codigo)

	return &response, nil
}
func (s *Service) List() ([]Response, error) {

	var codigos []db.CodigoQr

	if err := s.db.
		Order("url ASC").
		Find(&codigos).Error; err != nil {

		return nil, err
	}

	response := make([]Response, 0, len(codigos))

	for _, codigo := range codigos {
		response = append(response, toResponse(codigo))
	}

	return response, nil
}
func (s *Service) Update(id uint, dto UpdateDTO) (*Response, error) {

	// ==========================
	// Validar código QR
	// ==========================

	var codigo db.CodigoQr

	err := s.db.First(&codigo, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("el código QR no existe")
	}

	if err != nil {
		return nil, err
	}

	// ==========================
	// Validar URL única
	// ==========================

	var existing db.CodigoQr

	err = s.db.
		Where("url = ? AND id <> ?", dto.Url, id).
		First(&existing).Error

	if err == nil {
		return nil, errors.New("la URL ya existe")
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// ==========================
	// Actualizar
	// ==========================

	if err := s.db.Model(&codigo).Updates(map[string]interface{}{
		"url": dto.Url,
	}).Error; err != nil {

		return nil, err
	}

	if err := s.db.First(&codigo, id).Error; err != nil {
		return nil, err
	}

	response := toResponse(codigo)

	return &response, nil
}
func (s *Service) Delete(id uint) error {

	var codigo db.CodigoQr

	err := s.db.First(&codigo, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("el código QR no existe")
	}

	if err != nil {
		return err
	}

	codigo.Activo = !codigo.Activo

	return s.db.Save(&codigo).Error
}