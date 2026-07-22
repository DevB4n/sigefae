package archivo_origen

import "sigefae/internal/db"

type Response struct {
	ID      uint   `json:"id"`
	Nombre  string `json:"nombre"`
	Activo  bool   `json:"activo"`
}

func toResponse(archivo db.ArchivoOrigen) Response {

	return Response{
		ID:      archivo.ID,
		Nombre:  archivo.Nombre,
		Activo:  archivo.Activo,
	}
}