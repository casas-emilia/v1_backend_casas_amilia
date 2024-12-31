package controllers

import (
	"net/http"
	"time"
	"v1_prefabricadas/configs"
	"v1_prefabricadas/dto"
	"v1_prefabricadas/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Función auxiliar para manejar errores
func handleErrorTipo(c *gin.Context, err error, statusCode int, message string) {
	if err != nil {
		c.JSON(statusCode, gin.H{"error": message})
	}
}

// Función para crear un Tipo de Estructura
func CrearTipo(c *gin.Context) {
	var request dto.CrearTipoRequest

	// Validamos el body
	if err := c.ShouldBindJSON(&request); err != nil {
		handleErrorTipo(c, err, http.StatusBadRequest, err.Error())
		return
	}

	// Creación de Tipo
	tipo := models.Tipo{
		MaterialEstructura:  request.MaterialEstructura,
		DescripcionMaterial: request.DescripcionMaterial,
	}

	// Guardamos Tipo en la Base de Datos
	if err := configs.DB.Create(&tipo).Error; err != nil {
		handleErrorTipo(c, err, http.StatusInternalServerError, "No se pudo crear Empresa")
		return
	}

	// Retornamos mensaje con éxito
	c.JSON(http.StatusOK, gin.H{"message": "Tipo creada con éxito"})

	// Response
	response := dto.TipoResponse{
		ID:                  tipo.ID,
		CreatedAt:           tipo.CreatedAt,
		UpdatedAt:           tipo.UpdatedAt,
		MaterialEstructura:  tipo.MaterialEstructura,
		DescripcionMaterial: tipo.DescripcionMaterial,
	}

	// Respondemos con tipo creado
	c.JSON(http.StatusCreated, gin.H{"tipo": response})
}

// Función para obtener todos los Tipos de Estructuras exceptuando los eliminados logicamente
func ObtenerTipos(c *gin.Context) {
	var tipos []models.Tipo

	if err := configs.DB.Where("deleted_at IS NULL").Find(&tipos).Error; err != nil {
		handleErrorTipo(c, err, http.StatusInternalServerError, "No se pudieron obtener los tipos de estruturas")
		return
	}

	var tiposResponse []dto.TipoResponse

	for _, tipo := range tipos {
		tiposResponse = append(tiposResponse, dto.TipoResponse{
			ID:                  tipo.ID,
			CreatedAt:           tipo.CreatedAt,
			UpdatedAt:           tipo.UpdatedAt,
			MaterialEstructura:  tipo.MaterialEstructura,
			DescripcionMaterial: tipo.DescripcionMaterial,
		})

	}
	c.JSON(http.StatusOK, gin.H{"tipos": tiposResponse})
}

// Función para Obter el Tipo de Estructura de acuerso a su ID, excepto los eliminados logicamente
func ObtenerTipo(c *gin.Context) {
	var tipo models.Tipo
	id := c.Param("id")

	if err := configs.DB.Where("deleted_at IS NULL").First(&tipo, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			handleErrorTipo(c, nil, http.StatusNotFound, "Tipo estructura no encontrado")
			return
		}
		handleErrorTipo(c, err, http.StatusInternalServerError, "Error al obtener el tipo estructura")
		return
	}

	responseTipo := dto.TipoResponse{
		ID:                  tipo.ID,
		CreatedAt:           tipo.CreatedAt,
		UpdatedAt:           tipo.UpdatedAt,
		MaterialEstructura:  tipo.MaterialEstructura,
		DescripcionMaterial: tipo.DescripcionMaterial,
	}

	// Responder con éxito
	c.JSON(http.StatusOK, gin.H{"message": "Estructuras", "tipo": responseTipo})
}

// Función para actualizar datos de tipos de estructuras
func ActualizarTipo(c *gin.Context) {
	var request dto.ActualizarTipoRequest
	var tipo models.Tipo
	id := c.Param("id")

	// Bind JSON de la solicitud a la estructura DTO
	if err := c.ShouldBindJSON(&request); err != nil {
		handleErrorTipo(c, err, http.StatusBadRequest, err.Error())
		return
	}

	// Buscar Tipo estructura por su ID
	if err := configs.DB.Find(&tipo, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tipo estructura no encontrado"})
		return
	}

	// Actualizar los datos de Tipo de estructura
	tipo.MaterialEstructura = request.MaterialEstructura
	tipo.DescripcionMaterial = request.DescripcionMaterial

	// Guardar los datos en la base de datos
	if err := configs.DB.Save(tipo).Error; err != nil {
		handleErrorTipo(c, err, http.StatusInternalServerError, "No se pudieron actualizar los datos del Tipo de Estructura")
		return
	}

	// response
	response := dto.TipoResponse{
		MaterialEstructura:  tipo.MaterialEstructura,
		DescripcionMaterial: tipo.DescripcionMaterial,
	}

	// Responder con éxito
	c.JSON(http.StatusOK, gin.H{"message": "Tipo de Estructura actualizado exitosamente", "tipo": response})
}

// Función para eliminar(logicamente) un Tipo de estructura
func EliminarTipo(c *gin.Context) {
	var tipo models.Tipo
	id := c.Param("id")

	// Buscar el Tipo de Estructura por su ID
	if err := configs.DB.First(&tipo, id).Error; err != nil {
		handleErrorTipo(c, err, http.StatusNotFound, "Tipo de estructura no encontrado")
		return
	}

	// Verificar si el Tipo de estructura ya esta eliminado
	if tipo.DeletedAt != nil {
		handleErrorTipo(c, nil, http.StatusBadRequest, "El Tipo de estructura ya está eliminado")
		return
	}

	// Establecer Fecha y hora de la eliminación logica
	now := time.Now()
	tipo.DeletedAt = &now

	// Actualizar el registro de tipo en la base de datos
	if err := configs.DB.Save(&tipo).Error; err != nil {
		handleErrorTipo(c, err, http.StatusInternalServerError, "No se pudo eliminar el Tipo de estructura")
		return
	}

	// Responder con éxito
	c.JSON(http.StatusOK, gin.H{"message": "Tipo de estructuta eliminado exitosamente"})
}
