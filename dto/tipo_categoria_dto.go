package dto

type CrearTipo_categoriaRequest struct {
	CategoriaID uint `json:"categoria_id"`
	TipoID      uint `json:"tipo_id"`
}

type ActualizarTipo_categoriaRequest struct {
	CategoriaID uint `json:"categoria_id"`
	TipoID      uint `json:"tipo_id"`
}

type Tipo_categoriaResponse struct {
	//CreatedAt   time.Time `json:"created_at"`
	//UpdatedAt   time.Time `json:"updated_at"`
	CategoriaID        uint   `json:"categoria_id"`
	TipoID             uint   `json:"tipo_id"`
	MaterialEstructura string `json:"material_estructura"`
}
