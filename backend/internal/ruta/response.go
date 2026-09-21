package ruta

import (
	"sigefae/internal/db"
	"time"
)

type Response struct {
	ID      uint    `json:"id"`
	Nombre  string  `json:"nombre"`
	Zona    string  `json:"zona"`
	Version float64 `json:"version"`
	Activo  bool    `json:"activo"`

	AreaID uint   `json:"area_id"`
	Area   string `json:"area"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func toResponse(r db.Ruta) Response {

	area := ""

	if r.Area != nil {
		area = r.Area.Nombre
	}

	zona := r.Zona
	if zona == "" {
		zona = "BUCARAMANGA"
	}

	return Response{
		ID:        r.ID,
		Nombre:    r.Nombre,
		Zona:      zona,
		Version:   r.Version,
		Activo:    r.Activo,
		AreaID:    r.AreaID,
		Area:      area,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}
}
