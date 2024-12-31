package controllers

import (
	"net/http"
	"strconv"
	"time"
	"v1_prefabricadas/configs"
	"v1_prefabricadas/dto"
	"v1_prefabricadas/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Función auxiliar para manejar errores
func handleErrorCategoria(c *gin.Context, err error, statusCode int, message string) {
	if err != nil {
		c.JSON(statusCode, gin.H{"error": message})
	}
}

func CrearCategoria(c *gin.Context) {
	var request dto.CrearCategoriaRequest

	// Validamos el Body
	if err := c.ShouldBindJSON(&request); err != nil {
		handleErrorCategoria(c, err, http.StatusBadRequest, err.Error())
		return
	}

	// Creación de Categoria
	categoria := models.Categoria{
		NombreCategoria:      request.NombreCategoria,
		DescripcionCategoria: request.DescripcionCategoria,
	}
	if err := configs.DB.Create(&categoria).Error; err != nil {
		handleErrorCategoria(c, err, http.StatusInternalServerError, "No se pudo crear la Categoria")
		return
	}

	// Crear el Tipo de Categoria
	for _, tipo_categoriaReq := range request.Tipos {
		tipo_categoria := models.Tipo_categoria{
			CategoriaID: categoria.ID,
			TipoID:      tipo_categoriaReq.TipoID,
		}

		if err := configs.DB.Create(&tipo_categoria).Error; err != nil {
			handleErrorCategoria(c, err, http.StatusInternalServerError, "No se pudo crear Tipo_Categoria")
			return
		}
	}

	// Retornamos mensaje con éxito
	c.JSON(http.StatusOK, gin.H{"message": "Categoria creada con éxito"})

	var tiposResponse []dto.Tipo_categoriaResponse
	// Response
	response := dto.CategoriaResponse{
		ID:                   categoria.ID,
		CreatedAt:            categoria.CreatedAt,
		UpdatedAt:            categoria.UpdatedAt,
		NombreCategoria:      categoria.NombreCategoria,
		DescripcionCategoria: categoria.DescripcionCategoria,
		Tipos:                tiposResponse,
	}
	// Respondemos con Categoria creada
	c.JSON(http.StatusCreated, gin.H{"categoria": response})

}

// Función para obtener todas las Categorias
func ObtenerCategorias(c *gin.Context) {
	var categorias []models.Categoria
	var categoriasResponse []dto.CategoriaResponse

	//Obtener todas las Categorias
	if err := configs.DB.
		Where("deleted_at IS NULL").
		Preload("Tipo_categoria.Tipo").
		Find(&categorias).
		Error; err != nil {
		handleErrorCategoria(c, err, http.StatusInternalServerError, "No se pudieron obtener las Categorias")
		return
	}

	for _, categoria := range categorias {
		var tipoCategoria []dto.Tipo_categoriaResponse

		for _, nombreTipo := range categoria.Tipo_categoria {
			tipoCategoria = append(tipoCategoria, dto.Tipo_categoriaResponse{
				CategoriaID:        nombreTipo.CategoriaID,
				TipoID:             nombreTipo.Tipo.ID,
				MaterialEstructura: nombreTipo.Tipo.MaterialEstructura,
				// Otros campos si es necesario
			})
		}

		// Response
		categoriasResponse = append(categoriasResponse, dto.CategoriaResponse{
			ID:                   categoria.ID,
			CreatedAt:            categoria.CreatedAt,
			UpdatedAt:            categoria.UpdatedAt,
			NombreCategoria:      categoria.NombreCategoria,
			DescripcionCategoria: categoria.DescripcionCategoria,
			Tipos:                tipoCategoria,
		})
	}
	c.JSON(http.StatusOK, gin.H{"categorias": categoriasResponse})
}

// Función para obtener una Categoria
func ObtenerCategoria(c *gin.Context) {
	var categoria models.Categoria
	var categoriasResponse dto.CategoriaResponse
	var tipoCategoria []dto.Tipo_categoriaResponse
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64) // Convertir ID a uint
	if err != nil {
		handleErrorEstilo(c, nil, http.StatusBadRequest, "ID inválido")
		return
	}

	// Obtener la Categoria de acuerdo al ID enviado por el PATH
	if err := configs.DB.
		Preload("Tipo_categoria.Tipo").
		Where("deleted_at IS NULL").
		First(&categoria, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			HandleError(c, nil, http.StatusNotFound, "Categoria no encontrado")
			return
		}
		HandleError(c, err, http.StatusInternalServerError, "Error al obtener la Categoria")
		return
	}

	for _, nombreTipo := range categoria.Tipo_categoria {
		tipoCategoria = append(tipoCategoria, dto.Tipo_categoriaResponse{
			CategoriaID:        nombreTipo.CategoriaID,
			TipoID:             nombreTipo.Tipo.ID,
			MaterialEstructura: nombreTipo.Tipo.MaterialEstructura,
			// Otros campos si es necesario
		})
	}

	// Response
	categoriasResponse = dto.CategoriaResponse{
		ID:                   categoria.ID,
		CreatedAt:            categoria.CreatedAt,
		UpdatedAt:            categoria.UpdatedAt,
		NombreCategoria:      categoria.NombreCategoria,
		DescripcionCategoria: categoria.DescripcionCategoria,
		Tipos:                tipoCategoria,
	}

	//
	c.JSON(http.StatusOK, gin.H{"categoria": categoriasResponse})
}

// Función para actualizar datos de una Categoria
func ActualizarCategoria(c *gin.Context) {
	id := c.Param("id")
	var request dto.ActualizarCategoriaRequest
	var categoria models.Categoria

	// Bind JSON del request a la estructura DTO
	if err := c.ShouldBindJSON(&request); err != nil {
		handleErrorCategoria(c, err, http.StatusBadRequest, "Datos inválidos")
		return
	}

	// Buscar la categoría por ID y cargar los tipos relacionados
	if err := configs.DB.Preload("Tipo_categoria.Tipo").First(&categoria, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Categoría no encontrada"})
		return
	}

	// Actualizar datos de la categoría
	categoria.NombreCategoria = request.NombreCategoria
	categoria.DescripcionCategoria = request.DescripcionCategoria

	// Guardar cambios de la categoría en la base de datos
	if err := configs.DB.Save(&categoria).Error; err != nil {
		handleErrorCategoria(c, err, http.StatusInternalServerError, "No se pudo actualizar la categoría")
		return
	}

	// Crear un mapa de los tipos actuales para facilitar su gestión
	existingTipos := make(map[uint]models.Tipo_categoria)
	for _, t := range categoria.Tipo_categoria {
		existingTipos[t.TipoID] = t
	}

	// Procesar los tipos enviados en la solicitud
	for _, tipoReq := range request.Tipos {
		if _, exists := existingTipos[tipoReq.TipoID]; !exists {
			// Si el tipo no existe, agregarlo
			nuevoTipo := models.Tipo_categoria{
				CategoriaID: categoria.ID,
				TipoID:      tipoReq.TipoID,
			}
			if err := configs.DB.Create(&nuevoTipo).Error; err != nil {
				handleErrorCategoria(c, err, http.StatusInternalServerError, "Error al agregar tipo")
				return
			}
		}
		// Eliminar del mapa los tipos que se mantienen
		delete(existingTipos, tipoReq.TipoID)
	}

	// Eliminar los tipos que no se incluyeron en la solicitud
	for _, tipo := range existingTipos {
		if err := configs.DB.Delete(&tipo).Error; err != nil {
			handleErrorCategoria(c, err, http.StatusInternalServerError, "Error al eliminar tipo")
			return
		}
	}

	// Recargar la categoría con los tipos actualizados
	if err := configs.DB.Preload("Tipo_categoria.Tipo").First(&categoria, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al recargar la categoría"})
		return
	}

	// Preparar la respuesta
	var tiposResponse []dto.Tipo_categoriaResponse
	for _, tipo := range categoria.Tipo_categoria {
		tiposResponse = append(tiposResponse, dto.Tipo_categoriaResponse{
			CategoriaID:        tipo.CategoriaID,
			TipoID:             tipo.TipoID,
			MaterialEstructura: tipo.Tipo.MaterialEstructura,
		})
	}

	response := dto.CategoriaResponse{
		ID:                   categoria.ID,
		CreatedAt:            categoria.CreatedAt,
		UpdatedAt:            categoria.UpdatedAt,
		NombreCategoria:      categoria.NombreCategoria,
		DescripcionCategoria: categoria.DescripcionCategoria,
		Tipos:                tiposResponse,
	}

	// Enviar la respuesta
	c.JSON(http.StatusOK, gin.H{"Datos de Categoría": response})
}

// Función para eliminar logicamente Categorias
func EliminarCategoria(c *gin.Context) {
	id := c.Param("id")
	var categoria models.Categoria

	// Buscar el usuario por ID
	if err := configs.DB.First(&categoria, id).Error; err != nil {
		handleErrorCategoria(c, err, http.StatusNotFound, "Categoría no encontrado")
		return
	}

	// Verificar si el usuario ya está eliminado
	if categoria.DeletedAt != nil {
		handleErrorCategoria(c, nil, http.StatusBadRequest, "La Categoría ya está eliminado")
		return
	}

	// Establecer la fecha de eliminación lógica
	now := time.Now()
	categoria.DeletedAt = &now

	// Actualizar el registro del usuario en la base de datos
	if err := configs.DB.Save(&categoria).Error; err != nil {
		handleErrorCategoria(c, err, http.StatusInternalServerError, "No se pudo eliminar la Categoría")
		return
	}

	// Responder con éxito
	c.JSON(http.StatusOK, gin.H{"message": "Categoría eliminada exitosamente"})
}
