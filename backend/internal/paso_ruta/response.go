package paso_ruta

import (
	"time"

	"sigefae/internal/db"
)

type Response struct {
	ID uint `json:"id"`

	RutaID uint   `json:"ruta_id"`
	Ruta   string `json:"ruta"`

	Orden  int    `json:"orden"`
	Nombre string `json:"nombre"`

	UsuarioID uint   `json:"usuario_id"`
	Usuario   string `json:"usuario"`

	Activo bool `json:"activo"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func toResponse(p db.PasoRuta) Response {

	ruta := ""
	usuario := ""

	if p.Ruta != nil {
		ruta = p.Ruta.Nombre
	}

	if p.Usuario != nil {
		usuario = p.Usuario.Nombre
	}

	return Response{
		ID: p.ID,

		RutaID: p.RutaID,
		Ruta:   ruta,

		Orden:  p.Orden,
		Nombre: p.Nombre,

		UsuarioID: p.UsuarioID,
		Usuario:   usuario,

		Activo: p.Activo,

		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}
