package user

import "sigefae/internal/db"

func toResponse(user db.Usuario) Response {

	response := Response{
		ID:     user.ID,
		Nombre: user.Nombre,
		Email:  user.Email,
		Cargo:  user.Cargo,
		Activo: user.Activo,
	}

	if user.Rol != nil {
		response.Rol = user.Rol.Nombre
	}

	return response
}
