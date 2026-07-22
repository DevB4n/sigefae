package categoria_proveedor

import "sigefae/internal/db"

type Response struct {
	ID          uint   `json:"id"`
	Nombre      string `json:"nombre"`
	Descripcion string `json:"descripcion"`
	Activo      bool   `json:"activo"`
}

func toResponse(categoria db.CategoriaProveedor) Response {

	return Response{
		ID:          categoria.ID,
		Nombre:      categoria.Nombre,
		Descripcion: categoria.Descripcion,
		Activo:      categoria.Activo,
	}
}
