package dto

import "time"

// Estructura para la solicitud de actualización de Credenciales(email y password) de usuario
type CrearCredencial struct {
	Email    string `json:"email" binding:"required,email"`    // Validación de email
	Password string `json:"password" binding:"required,min=6"` // Validar que la contraseña tenga al menos 6 caracteres
}

type ActualizarCredencialRequest struct {
	Email    string `json:"email" binding:"required,email"` // Validación de email
	Password string `json:"password,omitempty"`             // Permitir que el password sea opcional
}

// DTO para mostrar la información de un usuario
type CredencialResponse struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	UsuarioID uint      `json:"usuario_id" binding:"required"` // FK de Usuario requerido
}
