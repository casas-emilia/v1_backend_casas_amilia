package dto

import (
	"mime/multipart"
	"time"
)

type CrearUsuarioRequest struct {
	PrimerNombre    string                `form:"primer_nombre" binding:"required"`
	SegundoNombre   string                `form:"segundo_nombre"`
	PrimerApellido  string                `form:"primer_apellido" binding:"required"`
	SegundoApellido string                `form:"segundo_nombre"`
	Image           *multipart.FileHeader `form:"image"`
}

type ActualizarUsuarioRequest struct {
	PrimerNombre    string                `form:"primer_nombre" binding:"required"`
	SegundoNombre   string                `form:"segundo_nombre"`
	PrimerApellido  string                `form:"primer_apellido" binding:"required"`
	SegundoApellido string                `form:"segundo_apellido"`
	Image           *multipart.FileHeader `form:"image"`
}

type UsuarioResponse struct {
	ID              uint      `json:"id"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	PrimerNombre    string    `json:"primer_nombre" binding:"required"`   // Campo obligatorio
	SegundoNombre   string    `json:"segundo_nombre"`                     // Opcional
	PrimerApellido  string    `json:"primer_apellido" binding:"required"` // Campo obligatorio
	SegundoApellido string    `json:"segundo_apellido"`                   // Opcional
	Image           string    `json:"image"`
	EmpresaID       uint      `json:"empresa_id"`
}
