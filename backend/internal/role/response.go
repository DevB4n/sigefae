package role

import "sigefae/internal/db"

func toResponse(role db.Rol) Response {

	return Response{
		ID:     role.ID,
		Nombre: role.Nombre,
		Activo: role.Activo,
	}
}
