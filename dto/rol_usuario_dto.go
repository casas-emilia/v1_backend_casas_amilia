package dto

import "time"

type CrearRol_usuarioRequest struct {
	UsuarioID uint `json:"usuario_id" binding:"required"`
	//RolID     uint `json:"rol_id" binding:"required"`
}

type ActualizarRol_usuarioRequest struct {
	UsuarioID uint `json:"usuario_id" binding:"required"`
	//RolID     uint `json:"rol_id" binding:"required"`
}

type Rol_usuarioResponse struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UsuarioID uint      `json:"usuario_id" binding:"required"`
	RolID     uint      `json:"rol_id" binding:"required"`
	//NombreRol string `json:"nombre_rol"`
}
