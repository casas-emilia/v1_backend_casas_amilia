package controllers

import (
	"errors"
	"net/http"
	"strconv"
	"time"
	"v1_prefabricadas/configs"
	"v1_prefabricadas/dto"
	"v1_prefabricadas/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Función para Crear una características con la siguiente estructura(clave : valor)
func CrearCaracteristica(c *gin.Context) {
	var request dto.CrearCaracteristicaRequest
	var caracteristica models.Caracteristica
	var caracteristicaResponse dto.CaracteristicaResponse

	idParamPrefabricada := c.Param("prefabricadaID")
	prefabricadaID, err := strconv.ParseUint(idParamPrefabricada, 10, 64)
	if err != nil {
		HandleError(c, nil, http.StatusBadRequest, "ID Prefabricada inválido")
		return
	}

	// Validamos el body
	if err := c.ShouldBindJSON(&request); err != nil {
		HandleError(c, err, http.StatusBadRequest, err.Error())
		return
	}

	// Creamos las caracteristicas
	caracteristica.Clave = request.Clave
	caracteristica.Valor = request.Valor
	caracteristica.PrefabricadaID = uint(prefabricadaID)

	// Guardamos en la base de datos
	if err := configs.DB.Create(&caracteristica).Error; err != nil {
		HandleError(c, err, http.StatusInternalServerError, "Error al guardar las características")
		return
	}

	caracteristicaResponse = dto.CaracteristicaResponse{
		ID:             caracteristica.ID,
		CreatedAt:      caracteristica.CreatedAt,
		UpdatedAt:      caracteristica.UpdatedAt,
		Clave:          caracteristica.Clave,
		Valor:          caracteristica.Valor,
		PrefabricadaID: caracteristica.PrefabricadaID,
	}

	// Mostrar/enviar un mensaje de éxito y el response
	c.JSON(http.StatusOK, gin.H{"message": "Característica guardada con éxito"})
	c.JSON(http.StatusOK, gin.H{"caracteristica": caracteristicaResponse})
}

// Función para obtener todas las Caracteristicas
func ObtenerCaracteristicas(c *gin.Context) {
	var caracteristicas []models.Caracteristica
	var caracteristicasResponse []dto.CaracteristicaResponse

	idParamPrefabricada := c.Param("prefabricadaID")
	prefabricadaID, err := strconv.ParseUint(idParamPrefabricada, 10, 64)
	if err != nil {
		HandleError(c, nil, http.StatusBadRequest, "ID Prefabricada inválido")
		return
	}

	// Buscar las características en la base de datos
	if err := configs.DB.Where("prefabricada_id = ?", prefabricadaID).Where("deleted_at IS NULL").Find(&caracteristicas).Error; err != nil {
		HandleError(c, err, http.StatusInternalServerError, "Características no encontradas")
		return
	}

	for _, caracteristica := range caracteristicas {
		caracteristicasResponse = append(caracteristicasResponse, dto.CaracteristicaResponse{
			ID:             caracteristica.ID,
			CreatedAt:      caracteristica.CreatedAt,
			UpdatedAt:      caracteristica.UpdatedAt,
			Clave:          caracteristica.Clave,
			Valor:          caracteristica.Valor,
			PrefabricadaID: caracteristica.PrefabricadaID,
		})
	}

	// Enviar response
	c.JSON(http.StatusOK, gin.H{"caracteristicas": caracteristicasResponse})
}

// Función para obtener característica de acuerdo al ID
func ObtenerCaracteristica(c *gin.Context) {
	var caracteristica models.Caracteristica
	var caracteristicaResponse dto.CaracteristicaResponse

	idParamPrefabricada := c.Param("prefabricadaID")
	prefabricadaID, err := strconv.ParseUint(idParamPrefabricada, 10, 64)
	if err != nil {
		HandleError(c, nil, http.StatusBadRequest, "ID Prefabricada inválido")
		return
	}

	idParamCaracteristica := c.Param("caracteristicaID")
	caracteristicaID, err := strconv.ParseUint(idParamCaracteristica, 10, 64)
	if err != nil {
		HandleError(c, nil, http.StatusBadRequest, "ID Característica inválido")
		return
	}

	// Buscar la caracteristica en la base de datos
	if err := configs.DB.Where("prefabricada_id = ?", prefabricadaID).Where("deleted_at IS NULL").First(&caracteristica, caracteristicaID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			HandleError(c, nil, http.StatusNotFound, "Característica no encontrada")
			return
		}
		HandleError(c, err, http.StatusInternalServerError, "Error al obtener Característica")
		return
	}

	caracteristicaResponse = dto.CaracteristicaResponse{
		ID:             caracteristica.ID,
		CreatedAt:      caracteristica.CreatedAt,
		UpdatedAt:      caracteristica.UpdatedAt,
		Clave:          caracteristica.Clave,
		Valor:          caracteristica.Valor,
		PrefabricadaID: caracteristica.PrefabricadaID,
	}

	// Enviar Caracteristica
	c.JSON(http.StatusOK, gin.H{"caracteristica": caracteristicaResponse})
}

// Función para actualizar datos de Característica
func ActualizarCaracteristica(c *gin.Context) {
	var request dto.ActualizarCaracteristicaRequest
	var caracteristica models.Caracteristica
	var caracteristicaResponse dto.CaracteristicaResponse

	idParamPrefabricada := c.Param("prefabricadaID")
	prefabricadaID, err := strconv.ParseUint(idParamPrefabricada, 10, 64)
	if err != nil {
		HandleError(c, nil, http.StatusBadRequest, "ID Prefabricada inválido")
		return
	}

	idParamCaracteristica := c.Param("caracteristicaID")
	caracteristicaID, err := strconv.ParseUint(idParamCaracteristica, 10, 64)
	if err != nil {
		HandleError(c, nil, http.StatusBadRequest, "ID Característica inválido")
		return
	}

	// Bind json del request a la estructura del dto
	if err := c.ShouldBindJSON(&request); err != nil {
		HandleError(c, nil, http.StatusBadRequest, "Error de datos")
		return
	}

	// Buscamos la característica
	if err := configs.DB.Where("prefabricada_id = ?", prefabricadaID).Where("deleted_at IS NULL").First(&caracteristica, caracteristicaID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			HandleError(c, nil, http.StatusNotFound, "Característica no encontrada")
			return
		}
		HandleError(c, err, http.StatusInternalServerError, "Error al obtener Característica")
		return
	}

	// Actualizamos los datos de la Característica
	caracteristica.Clave = request.Clave
	caracteristica.Valor = request.Valor

	// Guardamos los cambios en la Base de datos
	if err := configs.DB.Save(&caracteristica).Error; err != nil {
		HandleError(c, err, http.StatusInternalServerError, "No se pudo actualizar la Característica")
		return
	}

	caracteristicaResponse = dto.CaracteristicaResponse{
		ID:             caracteristica.ID,
		CreatedAt:      caracteristica.CreatedAt,
		UpdatedAt:      caracteristica.UpdatedAt,
		Clave:          caracteristica.Clave,
		Valor:          caracteristica.Valor,
		PrefabricadaID: caracteristica.PrefabricadaID,
	}

	// Enviar/mostrar mensaje éxito y la caraceristica actualizada
	c.JSON(http.StatusOK, gin.H{"message": "Característica actualizada exitosamente"})
	c.JSON(http.StatusOK, gin.H{"caracteristica": caracteristicaResponse})
}

// Función para eliminar lógicamente Características
func EliminarCaracteristica(c *gin.Context) {
	var caracteristica models.Caracteristica

	idParamPrefabricada := c.Param("prefabricadaID")
	prefabricadaID, err := strconv.ParseUint(idParamPrefabricada, 10, 64)
	if err != nil {
		HandleError(c, nil, http.StatusBadRequest, "ID Prefavricada inválido")
		return
	}

	idParamCaracteristica := c.Param("caracteristicaID")
	caracteristicaID, err := strconv.ParseUint(idParamCaracteristica, 10, 64)
	if err != nil {
		HandleError(c, nil, http.StatusBadRequest, "ID Característica inválido")
		return
	}

	// Buscar caracteíristica en la base de datos
	if err := configs.DB.Where("prefabricada_id = ?", prefabricadaID).Where("deleted_at IS NULL").Unscoped().First(&caracteristica, caracteristicaID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			HandleError(c, nil, http.StatusNotFound, "Característica no encontrada")
			return
		}
		HandleError(c, err, http.StatusInternalServerError, "Error al obtener Característica")
		return
	}

	// Verificamos si la característica ya se encuentra eliminada
	if caracteristica.DeletedAt != nil && !caracteristica.DeletedAt.IsZero() {
		HandleError(c, nil, http.StatusBadRequest, "La Característica ya se encuentra eliminada")
		return
	}

	// Guardar fecha y hora de la eliminación lógica
	now := time.Now()
	caracteristica.DeletedAt = &now

	// Guardamos la fecha y hora de la eliminación lógica en la base de datos
	if err := configs.DB.Save(&caracteristica).Error; err != nil {
		HandleError(c, err, http.StatusInternalServerError, "Error, no se pudo eliminar Característica")
		return
	}

	// Mostrar/enviar mensaje de eliminación exitosa
	c.JSON(http.StatusOK, gin.H{"message": "Característica eliminada exitosamente"})
}
