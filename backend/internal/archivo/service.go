package archivo

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
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
	// Validar Origen
	// ==========================

	var origen db.ArchivoOrigen

	err = s.db.First(
		&origen,
		req.OrigenID,
	).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("origen no encontrado")
	}

	if err != nil {
		return nil, err
	}

	// ==========================
	// Validar ruta repetida
	// ==========================

	var existing db.Archivo

	err = s.db.
		Where("ruta = ?", req.Ruta).
		First(&existing).Error

	if err == nil {
		return nil, errors.New("ya existe un archivo con esa ruta")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	// ==========================
	// Crear
	// ==========================

	archivo := db.Archivo{
		DocumentoRadicadoID: req.DocumentoRadicadoID,
		Nombre:              req.Nombre,
		Extension:           req.Extension,
		Peso:                req.Peso,
		Ruta:                req.Ruta,
		OrigenID:            req.OrigenID,
	}

	if err := s.db.Create(&archivo).Error; err != nil {
		return nil, err
	}

	response := toResponse(archivo)

	return &response, nil
}

func (s *Service) List() ([]Response, error) {

	var archivos []db.Archivo

	if err := s.db.
		Order("created_at DESC").
		Find(&archivos).Error; err != nil {

		return nil, err
	}

	response := make([]Response, 0, len(archivos))

	for _, archivo := range archivos {
		response = append(response, toResponse(archivo))
	}

	return response, nil
}

func (s *Service) Update(id uint, req UpdateDTO) (*Response, error) {

	var archivo db.Archivo

	err := s.db.First(&archivo, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("archivo no encontrado")
	}

	if err != nil {
		return nil, err
	}

	// ==========================
	// Validar Origen
	// ==========================

	var origen db.ArchivoOrigen

	err = s.db.First(
		&origen,
		req.OrigenID,
	).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("origen no encontrado")
	}

	if err != nil {
		return nil, err
	}

	// ==========================
	// Validar ruta repetida
	// ==========================

	var existing db.Archivo

	err = s.db.
		Where("ruta = ? AND id <> ?", req.Ruta, id).
		First(&existing).Error

	if err == nil {
		return nil, errors.New("ya existe un archivo con esa ruta")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	// ==========================
	// Actualizar
	// ==========================

	if err := s.db.Model(&archivo).Updates(map[string]any{
		"nombre":    req.Nombre,
		"extension": req.Extension,
		"peso":      req.Peso,
		"ruta":      req.Ruta,
		"origen_id": req.OrigenID,
	}).Error; err != nil {

		return nil, err
	}

	if err := s.db.First(&archivo, archivo.ID).Error; err != nil {
		return nil, err
	}

	response := toResponse(archivo)

	return &response, nil
}

func (s *Service) Delete(id uint) error {

	var archivo db.Archivo

	err := s.db.First(&archivo, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("archivo no encontrado")
	}

	if err != nil {
		return err
	}

	// ==========================
	// Eliminar archivo físico
	// ==========================

	if archivo.Ruta != "" {

		if err := os.Remove(archivo.Ruta); err != nil && !os.IsNotExist(err) {
			return err
		}
	}

	// ==========================
	// Eliminar registro
	// ==========================

	return s.db.Delete(&archivo).Error
}
func (s *Service) UploadAnexo(documentoRadicadoID uint, fileHeader *multipart.FileHeader, rutaBase string) (*Response, error) {
	// Crear carpeta si no existe
	dir := filepath.Join(rutaBase, "anexos", fmt.Sprintf("%d", documentoRadicadoID))
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return nil, err
	}

	// Nombre único para evitar colisiones
	ext := filepath.Ext(fileHeader.Filename)
	nombreBase := strings.TrimSuffix(fileHeader.Filename, ext)
	nombreLimpio := strings.ReplaceAll(nombreBase, " ", "_")
	nombreFinal := fmt.Sprintf("%s_%d%s", nombreLimpio, time.Now().Unix(), ext)
	rutaFinal := filepath.Join(dir, nombreFinal)

	// Guardar archivo físico
	src, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	dst, err := os.Create(rutaFinal)
	if err != nil {
		return nil, err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return nil, err
	}

	// Obtener tamaño
	info, err := dst.Stat()
	if err != nil {
		return nil, err
	}

	// Buscar o crear origen "Anexo"
	var origen db.ArchivoOrigen
	if err := s.db.Where("nombre = ?", "Anexo").First(&origen).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			origen = db.ArchivoOrigen{Nombre: "Anexo", Activo: true}
			if err := s.db.Create(&origen).Error; err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	// Crear registro en BD
	archivo := db.Archivo{
		DocumentoRadicadoID: documentoRadicadoID,
		Nombre:              fileHeader.Filename,
		Extension:           strings.TrimPrefix(ext, "."),
		Peso:                info.Size(),
		Ruta:                rutaFinal,
		OrigenID:            origen.ID,
	}

	if err := s.db.Create(&archivo).Error; err != nil {
		return nil, err
	}

	resp := toResponse(archivo)
	return &resp, nil
}
func (s *Service) GetByID(id uint) (*db.Archivo, error) {
	var archivo db.Archivo
	if err := s.db.First(&archivo, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("archivo no encontrado")
		}
		return nil, err
	}
	return &archivo, nil
}

func (s *Service) Reemplazar(id uint, fileHeader *multipart.FileHeader) (*Response, error) {
	var archivo db.Archivo
	if err := s.db.First(&archivo, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("archivo no encontrado")
		}
		return nil, err
	}

	// Validar que sea PDF
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if ext != ".pdf" {
		return nil, errors.New("solo se permiten archivos PDF")
	}

	// Borrar archivo físico viejo
	if archivo.Ruta != "" {
		if err := os.Remove(archivo.Ruta); err != nil && !os.IsNotExist(err) {
			return nil, err
		}
	}

	// Guardar archivo nuevo en la misma carpeta
	dir := filepath.Dir(archivo.Ruta)
	if dir == "" || dir == "." {
		dir = filepath.Join("storage", "anexos", fmt.Sprintf("%d", archivo.DocumentoRadicadoID))
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return nil, err
		}
	}

	nombreBase := strings.TrimSuffix(fileHeader.Filename, ext)
	nombreLimpio := strings.ReplaceAll(nombreBase, " ", "_")
	nombreFinal := fmt.Sprintf("%s_%d%s", nombreLimpio, time.Now().Unix(), ext)
	rutaFinal := filepath.Join(dir, nombreFinal)

	src, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	dst, err := os.Create(rutaFinal)
	if err != nil {
		return nil, err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return nil, err
	}

	info, err := dst.Stat()
	if err != nil {
		return nil, err
	}

	// Actualizar registro (mismo ID, mismo radicado, mismo origen)
	archivo.Nombre = fileHeader.Filename
	archivo.Extension = strings.TrimPrefix(ext, ".")
	archivo.Peso = info.Size()
	archivo.Ruta = rutaFinal

	if err := s.db.Save(&archivo).Error; err != nil {
		return nil, err
	}

	resp := toResponse(archivo)
	return &resp, nil
}
