package user

type CreateRequest struct {
	Nombre      string `json:"nombre" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	Contrasena  string `json:"contrasena" binding:"required,min=6"`
	Cargo       string `json:"cargo" binding:"required"`
	RolID       uint   `json:"rol_id" binding:"required"`
}

type Response struct {
	ID     uint   `json:"id"`
	Nombre string `json:"nombre"`
	Email  string `json:"email"`
	Cargo  string `json:"cargo"`
	Rol    string `json:"rol"`
}