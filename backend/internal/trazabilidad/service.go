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
func (s *Service) List(documentoRadicadoID uint) ([]Response, error) {

	var trazabilidades []db.Trazabilidad

	query := s.db.Preload("Usuario").Order("fecha DESC")
	if documentoRadicadoID > 0 {
		query = query.Where("documento_radicado_id = ?", documentoRadicadoID)
	}

	if err := query.Find(&trazabilidades).Error; err != nil {

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
func (s *Service) ListPorArea(areaID uint, fechaDesde, fechaHasta string) ([]RadicadoConTrazabilidadResponse, error) {
	// 1. Rutas que pertenecen al área seleccionada
	var rutaIDs []uint
	if err := s.db.Model(&db.Ruta{}).Where("area_id = ?", areaID).Pluck("id", &rutaIDs).Error; err != nil {
		return nil, err
	}
	if len(rutaIDs) == 0 {
		return []RadicadoConTrazabilidadResponse{}, nil
	}

	// 2. Radicados que usen esas rutas, con trazabilidad
	var radicados []db.DocumentoRadicado
	if err := s.db.
		Preload("DocumentoComercial").
		Preload("DocumentoComercial.Proveedor").
		Preload("DocumentoComercial.Area").
		Preload("Ruta").
		Preload("Ruta.Area").
		Preload("Trazabilidades", func(db *gorm.DB) *gorm.DB {
			return db.Order("fecha ASC")
		}).
		Preload("Trazabilidades.Usuario").
		Where("ruta_id IN ? AND DATE(fecha_radicacion) BETWEEN ? AND ?", rutaIDs, fechaDesde, fechaHasta).
		Order("fecha_radicacion DESC").
		Find(&radicados).Error; err != nil {
		return nil, err
	}

	// 3. Mapear respuesta (área siempre desde la ruta, que es lo que el usuario eligió)
	response := make([]RadicadoConTrazabilidadResponse, 0, len(radicados))
	for _, r := range radicados {
		item := RadicadoConTrazabilidadResponse{
			ID:              r.ID,
			NumeroRadicado:  r.NumeroRadicado,
			FechaRadicacion: r.FechaRadicacion,
			EstadoPosesion:  r.EstadoPosesion,
		}

		if r.DocumentoComercial != nil && r.DocumentoComercial.Proveedor != nil {
			item.Proveedor = ProveedorMiniResponse{
				ID:              r.DocumentoComercial.Proveedor.ID,
				RazonSocial:     r.DocumentoComercial.Proveedor.RazonSocial,
				NumeroDocumento: r.DocumentoComercial.Proveedor.NumeroDocumento,
			}
		}

		// El área "oficial" del radicado es la de la ruta que se eligió al radicar
		if r.Ruta != nil && r.Ruta.Area != nil {
			item.Area = AreaMiniResponse{
				ID:     r.Ruta.Area.ID,
				Nombre: r.Ruta.Area.Nombre,
			}
		} else if r.DocumentoComercial != nil && r.DocumentoComercial.Area != nil {
			// Fallback solo si la ruta no tiene área
			item.Area = AreaMiniResponse{
				ID:     r.DocumentoComercial.Area.ID,
				Nombre: r.DocumentoComercial.Area.Nombre,
			}
		}

		for _, t := range r.Trazabilidades {
			item.Trazabilidad = append(item.Trazabilidad, toResponse(t))
		}
		response = append(response, item)
	}

	return response, nil
}