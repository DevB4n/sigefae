package zip

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

type Extractor struct{}

func NewExtractor() *Extractor {
	return &Extractor{}
}

func (e *Extractor) Extract(zipPath string, destination string) error {
	// Abrir el archivo ZIP
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer reader.Close()

	// Crear la carpeta destino
	if err := os.MkdirAll(destination, 0755); err != nil {
		return err
	}

	// Recorrer todos los archivos del ZIP
	for _, file := range reader.File {
		path := filepath.Join(destination, file.Name)

		// Si es una carpeta, crearla
		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(path, file.Mode()); err != nil {
				return err
			}
			continue
		}

		// Crear carpetas padre si existen dentro del ZIP
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return err
		}

		// Abrir archivo dentro del ZIP
		src, err := file.Open()
		if err != nil {
			return err
		}

		// Crear archivo destino
		dst, err := os.OpenFile(
			path,
			os.O_CREATE|os.O_WRONLY|os.O_TRUNC,
			file.Mode(),
		)

		if err != nil {
			src.Close()
			return err
		}

		// Copiar contenido
		if _, err := io.Copy(dst, src); err != nil {
			dst.Close()
			src.Close()
			return err
		}

		dst.Close()
		src.Close()
	}

	return nil
}
