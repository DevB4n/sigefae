package contacto

type CreateDTO struct {
	ProveedorID uint   `json:"proveedor_id" binding:"required"`
	Nombre      string `json:"nombre" binding:"required"`
	Cargo       string `json:"cargo" binding:"required"`
	Telefono    string `json:"telefono" binding:"required"`
	Correo      string `json:"correo" binding:"required,email"`
}

type UpdateDTO struct {
	ProveedorID uint   `json:"proveedor_id" binding:"required"`
	Nombre      string `json:"nombre" binding:"required"`
	Cargo       string `json:"cargo" binding:"required"`
	Telefono    string `json:"telefono" binding:"required"`
	Correo      string `json:"correo" binding:"required,email"`
}

type UpdateStatusDTO struct {
	Activo bool `json:"activo"`
}