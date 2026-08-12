package proveedor

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

	// ── Defaults para creación rápida (si no mandan catálogo, toma el primero) ──
	if req.CategoriaID == 0 {
		var first db.CategoriaProveedor
		if err := s.db.First(&first).Error; err == nil {
			req.CategoriaID = first.ID
		}
	}
	if req.TipoPersonaID == 0 {
		var first db.TipoPersona
		if err := s.db.First(&first).Error; err == nil {
			req.TipoPersonaID = first.ID
		}
	}
	if req.ActividadEconomicaID == 0 {
		var first db.ActividadEconomica
		if err := s.db.First(&first).Error; err == nil {
			req.ActividadEconomicaID = first.ID
		}
	}
	if req.DireccionID == 0 {
		var first db.Direccion
		if err := s.db.First(&first).Error; err == nil {
			req.DireccionID = first.ID
		}
	}

	// ==========================
	// Validar Tipo Documento
	// ==========================

	var tipoDocumento db.TipoDocumento

	err := s.db.First(&tipoDocumento, req.TipoDocumentoID).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("tipo de documento no encontrado")
	}

	if err != nil {
		return nil, err
	}

	// ==========================
	// Validar Categoría
	// ==========================

	var categoria db.CategoriaProveedor

	err = s.db.First(&categoria, req.CategoriaID).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("categoría no encontrada")
	}

	if err != nil {
		return nil, err
	}

	// ==========================
	// Validar Tipo Persona
	// ==========================

	var tipoPersona db.TipoPersona

	err = s.db.First(&tipoPersona, req.TipoPersonaID).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("tipo de persona no encontrado")
	}

	if err != nil {
		return nil, err
	}

	// ==========================
	// Validar Actividad Económica
	// ==========================

	var actividad db.ActividadEconomica

	err = s.db.First(&actividad, req.ActividadEconomicaID).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("actividad económica no encontrada")
	}

	if err != nil {
		return nil, err
	}

	// ==========================
	// Validar Dirección
	// ==========================

	var direccion db.Direccion

	err = s.db.First(&direccion, req.DireccionID).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("dirección no encontrada")
	}

	if err != nil {
		return nil, err
	}

	// ==========================
	// Validar Ruta (Opcional)
	// ==========================

	if req.RutaPredeterminadaID != nil {

		var ruta db.Ruta

		err = s.db.First(&ruta, *req.RutaPredeterminadaID).Error

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("ruta predeterminada no encontrada")
		}

		if err != nil {
			return nil, err
		}
	}

	// ==========================
	// Validar Documento Repetido
	// ==========================

	var existing db.Proveedor

	err = s.db.
		Where("numero_documento = ?", req.NumeroDocumento).
		First(&existing).Error

	if err == nil {
		return nil, errors.New("el número de documento ya existe")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// ==========================
	// Crear
	// ==========================

	proveedor := db.Proveedor{
		TipoDocumentoID:      req.TipoDocumentoID,
		NumeroDocumento:      req.NumeroDocumento,
		CategoriaID:          req.CategoriaID,
		RutaPredeterminadaID: req.RutaPredeterminadaID,
		RazonSocial:          req.RazonSocial,
		NombreComercial:      req.NombreComercial,
		TipoPersonaID:        req.TipoPersonaID,
		ActividadEconomicaID: req.ActividadEconomicaID,
		DireccionID:          req.DireccionID,
		Activo:               true,
	}

	if err := s.db.Create(&proveedor).Error; err != nil {
		return nil, err
	}

	// ==========================
	// Recargar con relaciones
	// ==========================

	err = s.db.
		Preload("TipoDocumento").
		Preload("Categoria").
		Preload("RutaPredeterminada").
		Preload("TipoPersona").
		Preload("ActividadEconomica").
		Preload("Direccion").
		First(&proveedor, proveedor.ID).Error

	if err != nil {
		return nil, err
	}

	response := toResponse(proveedor)

	return &response, nil
}
func (s *Service) List() ([]Response, error) {

	var proveedores []db.Proveedor

	if err := s.db.
		Preload("TipoDocumento").
		Preload("Categoria").
		Preload("RutaPredeterminada").
		Preload("TipoPersona").
		Preload("ActividadEconomica").
		Preload("Direccion").
		Order("razon_social ASC").
		Find(&proveedores).Error; err != nil {

		return nil, err
	}

	response := make([]Response, 0, len(proveedores))

	for _, proveedor := range proveedores {
		response = append(
			response,
			toResponse(proveedor),
		)
	}

	return response, nil
}

func (s *Service) ListNormasReparto(proveedorID, rutaID uint) ([]db.ProveedorNormaReparto, error) {
	var normas []db.ProveedorNormaReparto

	if err := s.db.
		Preload("NormaReparto").
		Where("proveedor_id = ? AND ruta_id = ?", proveedorID, rutaID).
		Find(&normas).Error; err != nil {
		return nil, err
	}

	return normas, nil
}

func (s *Service) UpdateStatus(id uint, activo bool) error {

	var proveedor db.Proveedor

	err := s.db.First(&proveedor, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("proveedor no encontrado")
	}

	if err != nil {
		return err
	}

	proveedor.Activo = activo

	return s.db.Save(&proveedor).Error
}
func (s *Service) Update(id uint, req UpdateRequest) (*Response, error) {

	var proveedor db.Proveedor

	err := s.db.First(&proveedor, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("proveedor no encontrado")
	}

	if err != nil {
		return nil, err
	}

	// ==========================
	// Validar Tipo Documento
	// ==========================

	var tipoDocumento db.TipoDocumento

	err = s.db.First(&tipoDocumento, req.TipoDocumentoID).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("tipo de documento no encontrado")
	}

	if err != nil {
		return nil, err
	}

	// ==========================
	// Validar Categoría
	// ==========================

	var categoria db.CategoriaProveedor

	err = s.db.First(&categoria, req.CategoriaID).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("categoría no encontrada")
	}

	if err != nil {
		return nil, err
	}

	// ==========================
	// Validar Tipo Persona
	// ==========================

	var tipoPersona db.TipoPersona

	err = s.db.First(&tipoPersona, req.TipoPersonaID).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("tipo de persona no encontrado")
	}

	if err != nil {
		return nil, err
	}

	// ==========================
	// Validar Actividad Económica
	// ==========================

	var actividad db.ActividadEconomica

	err = s.db.First(&actividad, req.ActividadEconomicaID).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("actividad económica no encontrada")
	}

	if err != nil {
		return nil, err
	}

	// ==========================
	// Validar Dirección
	// ==========================

	var direccion db.Direccion

	err = s.db.First(&direccion, req.DireccionID).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("dirección no encontrada")
	}

	if err != nil {
		return nil, err
	}

	// ==========================
	// Validar Ruta (Opcional)
	// ==========================

	if req.RutaPredeterminadaID != nil {

		var ruta db.Ruta

		err = s.db.First(&ruta, *req.RutaPredeterminadaID).Error

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("ruta predeterminada no encontrada")
		}

		if err != nil {
			return nil, err
		}
	}

	// ==========================
	// Validar Documento Repetido
	// ==========================

	var existing db.Proveedor

	err = s.db.
		Where("numero_documento = ? AND id <> ?", req.NumeroDocumento, id).
		First(&existing).Error

	if err == nil {
		return nil, errors.New("el número de documento ya existe")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// ==========================
	// Actualizar
	// ==========================

	if err := s.db.Model(&proveedor).Updates(map[string]any{
		"tipo_documento_id":      req.TipoDocumentoID,
		"numero_documento":       req.NumeroDocumento,
		"categoria_id":           req.CategoriaID,
		"ruta_predeterminada_id": req.RutaPredeterminadaID,
		"razon_social":           req.RazonSocial,
		"nombre_comercial":       req.NombreComercial,
		"tipo_persona_id":        req.TipoPersonaID,
		"actividad_economica_id": req.ActividadEconomicaID,
		"direccion_id":           req.DireccionID,
	}).Error; err != nil {

		return nil, err
	}

	// ==========================
	// Recargar con relaciones
	// ==========================

	err = s.db.
		Preload("TipoDocumento").
		Preload("Categoria").
		Preload("RutaPredeterminada").
		Preload("TipoPersona").
		Preload("ActividadEconomica").
		Preload("Direccion").
		First(&proveedor, proveedor.ID).Error

	if err != nil {
		return nil, err
	}

	response := toResponse(proveedor)

	return &response, nil
}
