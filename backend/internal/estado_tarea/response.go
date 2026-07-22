package estado_tarea

import "sigefae/internal/db"

type Response struct {
	ID      uint   `json:"id"`
	Nombre  string `json:"nombre"`
	Activo  bool   `json:"activo"`
}

func toResponse(estado db.EstadoTarea) Response {

	return Response{
		ID:      estado.ID,
		Nombre:  estado.Nombre,
		Activo:  estado.Activo,
	}
}