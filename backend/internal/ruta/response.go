package ruta

import (
	"sigefae/internal/db"
	"time"
)

type Response struct {
	ID      uint    `json:"id"`
	Nombre  string  `json:"nombre"`
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

	return Response{
		ID:        r.ID,
		Nombre:    r.Nombre,
		Version:   r.Version,
		Activo:    r.Activo,
		AreaID:    r.AreaID,
		Area:      area,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}
}
