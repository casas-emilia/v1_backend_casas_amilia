package dto

import "time"

// Estructura para la solicitud de creacion de Rol
type CrearRolRequest struct {
	NombreRol      string `json:"nombre_rol" binding:"required"` // nombre_rol obligatorio
	DescripcionRol string `json:"descripcion_rol"`               // descripcion_rol opcional
}

// Estructura para la solicitud de actualización de Rol
type ActualizarRolRequest struct {
	NombreRol      string `json:"nombre_rol" binding:"required"`      // nombre_rol obligatorio
	DescripcionRol string `json:"descripcion_rol" binding:"required"` // descripcion_rol opcional
}

// Estructura para mostrar la información de Rol
type RolResponse struct {
	ID             uint      `json:"id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	NombreRol      string    `json:"nombre_rol" binding:"required"` // nombre_rol obligatorio
	DescripcionRol string    `json:"descripcion_rol"`               // descripcion_rol opcional
}
