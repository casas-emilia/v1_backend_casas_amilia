package dto

import "time"

type CrearCategoriaRequest struct {
	NombreCategoria      string                       `json:"nombre_categoria" binding:"required"`
	DescripcionCategoria string                       `json:"descripcion_categoria"`
	Tipos                []CrearTipo_categoriaRequest `json:"tipos"`
}

type ActualizarCategoriaRequest struct {
	NombreCategoria      string                            `json:"nombre_categoria" binding:"required"`
	DescripcionCategoria string                            `json:"descripcion_categoria"`
	Tipos                []ActualizarTipo_categoriaRequest `json:"tipos"`
}

type CategoriaResponse struct {
	ID                   uint                     `json:"id"`
	CreatedAt            time.Time                `json:"created_At"`
	UpdatedAt            time.Time                `json:"updated_at"`
	NombreCategoria      string                   `json:"nombre_categoria"`
	DescripcionCategoria string                   `json:"descripcion_categoria"`
	Tipos                []Tipo_categoriaResponse `json:"tipos"`
}
