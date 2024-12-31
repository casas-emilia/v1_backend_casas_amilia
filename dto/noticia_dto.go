package dto

import "time"

type CrearNoticiaRequest struct {
	TituloNoticia     string `json:"titulo_noticia" binding:"required"`
	DesarrolloNoticia string `json:"desarrollo_noticia" binding:"required"`
	//EmpresaID         uint   `json:"empresa_id" binding:"required"`
	//UsuarioID         uint                         `json:"usuario_id" binding:"required"`
	//Imagenes []CrearImagen_noticiaRequest `json:"imagenes"`
}

type ActualizarNoticiaRequest struct {
	TituloNoticia     string `json:"titulo_noticia" binding:"required"`
	DesarrolloNoticia string `json:"desarrollo_noticia" binding:"required"`
	//Imagenes          []CrearImagen_noticiaRequest `json:"imagenes"`
}

type NoticiaResponse struct {
	ID                uint      `json:"id"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	TituloNoticia     string    `json:"titulo_noticia" binding:"required"`
	DesarrolloNoticia string    `json:"desarrollo_noticia" binding:"required"`
	// UsuarioID         uint      `json:"usuario_id" binding:"required"`
	EmpresaID uint `json:"empresa_id"`
	//Imagenes          []CrearImagen_noticiaRequest `json:"imagenes"`
}
