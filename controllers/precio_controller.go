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

// Función para Crear un nuevo Precio
func CrearPrecio(c *gin.Context) {
	var request dto.CrearPrecioRequest
	var precio models.Precio
	var precioResponse dto.PrecioResponse

	idParamPrefabricada := c.Param("prefabricadaID")
	prefabricadaID, err := strconv.ParseUint(idParamPrefabricada, 10, 64)
	if err != nil {
		HandleError(c, nil, http.StatusBadRequest, "ID Prefabricada inválido")
		return
	}

	// Validamos el body
	if err := c.ShouldBindJSON(&request); err != nil {
		HandleError(c, nil, http.StatusBadRequest, "Error de datos"+err.Error())
		return
	}

	// Creamos el precio
	precio.NombrePrecio = request.NombrePrecio
	precio.DescripcionPrecio = request.DescripcionPrecio
	precio.ValorPrefabricada = request.ValorPrefabricada
	precio.PrefabricadaID = uint(prefabricadaID)

	// Guardamos el precio en la base de datos
	if err := configs.DB.Create(&precio).Error; err != nil {
		HandleError(c, err, http.StatusInternalServerError, "Error, no se pudo guardar el Precio")
		return
	}

	precioResponse = dto.PrecioResponse{
		ID:                precio.ID,
		CreatedAt:         precio.CreatedAt,
		UpdatedAt:         precio.UpdatedAt,
		NombrePrecio:      precio.NombrePrecio,
		DescripcionPrecio: precio.DescripcionPrecio,
		ValorPrefabricada: precio.ValorPrefabricada,
		PrefabricadaID:    precio.PrefabricadaID,
	}

	// Enviar/mostrar mensaje éxitoso y el response del Precio
	c.JSON(http.StatusOK, gin.H{
		"message": "Precio guardado con éxito",
		"precio":  precioResponse,
	})
}

// Función para obtener todos los precios de una prefabricada
func ObtenerPrecios(c *gin.Context) {
	var precios []models.Precio
	var preciosResponse []dto.PrecioResponse

	idParamPrefabricada := c.Param("prefabricadaID")
	prefabricadaID, err := strconv.ParseUint(idParamPrefabricada, 10, 64)
	if err != nil {
		HandleError(c, nil, http.StatusBadRequest, "ID Prefabricada inválido")
		return
	}

	// Buscar todos los precios en la base de datos
	if err := configs.DB.
		Preload("Incluye", func(db *gorm.DB) *gorm.DB {
			return db.Where("deleted_at IS NULL") // Condición para no cargar "incluye" eliminados lógicamente
		}).
		Where("prefabricada_id = ?", prefabricadaID).
		Where("deleted_at IS NULL").
		Find(&precios).Error; err != nil {
		HandleError(c, err, http.StatusInternalServerError, "Precios no encontrados")
		return
	}

	for _, precio := range precios {

		var incluyesResponse []dto.IncluyeResponse

		for _, incluye := range precio.Incluye {
			incluyesResponse = append(incluyesResponse, dto.IncluyeResponse{
				ID:            incluye.ID,
				CreatedAt:     incluye.UpdatedAt,
				UpdatedAt:     incluye.UpdatedAt,
				NombreIncluye: incluye.NombreIncluye,
				PrecioID:      incluye.PrecioID,
			})
		}

		preciosResponse = append(preciosResponse, dto.PrecioResponse{
			ID:                precio.ID,
			CreatedAt:         precio.CreatedAt,
			UpdatedAt:         precio.UpdatedAt,
			NombrePrecio:      precio.NombrePrecio,
			DescripcionPrecio: precio.DescripcionPrecio,
			ValorPrefabricada: precio.ValorPrefabricada,
			PrefabricadaID:    precio.PrefabricadaID,
			Incluyes:          incluyesResponse,
		})
	}

	// Mostrar/enviar precios
	c.JSON(http.StatusOK, gin.H{"precios": preciosResponse})
}

// Función para obtener un solo Precio de acuerdo al ID
func ObtenerPrecio(c *gin.Context) {
	var precio models.Precio
	var precioResponse dto.PrecioResponse

	idParamPrefabricada := c.Param("prefabricadaID")
	prefabricadaID, err := strconv.ParseUint(idParamPrefabricada, 10, 64)
	if err != nil {
		HandleError(c, nil, http.StatusBadRequest, "ID Prefabricada inválido")
		return
	}

	idParamPrecio := c.Param("precioID")
	precioID, err := strconv.ParseUint(idParamPrecio, 10, 64)
	if err != nil {
		HandleError(c, nil, http.StatusBadRequest, "ID Precio inválido")
		return
	}

	// Buscar el Precio en la base de datos de acuerdo a su ID
	if err := configs.DB.Where("prefabricada_id = ?", prefabricadaID).Where("deleted_at IS NULL").First(&precio, precioID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			HandleError(c, nil, http.StatusNotFound, "Precio no encontrada")
			return
		}
		HandleError(c, err, http.StatusInternalServerError, "Error al obtener Precio")
		return
	}

	precioResponse = dto.PrecioResponse{
		ID:                precio.ID,
		CreatedAt:         precio.CreatedAt,
		UpdatedAt:         precio.UpdatedAt,
		NombrePrecio:      precio.NombrePrecio,
		DescripcionPrecio: precio.DescripcionPrecio,
		ValorPrefabricada: precio.ValorPrefabricada,
		PrefabricadaID:    precio.PrefabricadaID,
	}

	// Mostrar/enviar precio
	c.JSON(http.StatusOK, gin.H{"precio": precioResponse})
}

// Función para actualizar Precio
func ActualizarPrecio(c *gin.Context) {
	var request dto.ActualizarPrecioRequest
	var precio models.Precio
	var precioResponse dto.PrecioResponse

	idParamPrefabricada := c.Param("prefabricadaID")
	prefabricadaID, err := strconv.ParseUint(idParamPrefabricada, 10, 64)
	if err != nil {
		HandleError(c, nil, http.StatusBadRequest, "ID Prefabricada inválido")
		return
	}

	idParamPrecio := c.Param("precioID")
	precioID, err := strconv.ParseUint(idParamPrecio, 10, 64)
	if err != nil {
		HandleError(c, nil, http.StatusBadRequest, "ID Precio inválido")
		return
	}

	// Bind json del request a la estructura del dto
	if err := c.ShouldBindJSON(&request); err != nil {
		HandleError(c, nil, http.StatusBadRequest, "Error de Datos"+err.Error())
		return
	}

	// Buscar Precio en la Base de Datos de acuerdo a su ID
	if err := configs.DB.Where("prefabricada_id = ?", prefabricadaID).Where("deleted_at IS NULL").First(&precio, precioID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			HandleError(c, nil, http.StatusNotFound, "Precio no encontrado")
			return
		}
		HandleError(c, err, http.StatusInternalServerError, "Error al obtener Precio")
		return
	}

	// Actualizar los datos del Precio
	precio.NombrePrecio = request.NombrePrecio
	precio.DescripcionPrecio = request.DescripcionPrecio
	precio.ValorPrefabricada = request.ValorPrefabricada

	// Guardar datos actualizados en la base de datos
	if err := configs.DB.Save(&precio).Error; err != nil {
		HandleError(c, err, http.StatusBadRequest, "Error, no se pudo actualizar la información"+err.Error())
	}

	precioResponse = dto.PrecioResponse{
		ID:                precio.ID,
		CreatedAt:         precio.CreatedAt,
		UpdatedAt:         precio.UpdatedAt,
		NombrePrecio:      precio.NombrePrecio,
		DescripcionPrecio: precio.DescripcionPrecio,
		ValorPrefabricada: precio.ValorPrefabricada,
		PrefabricadaID:    precio.PrefabricadaID,
	}

	// Mostrar/enviar mensaje exitoso y el Precio
	c.JSON(http.StatusOK, gin.H{
		"message": "Precio actualizado exitosamente",
		"precio":  precioResponse,
	})
}

// Función para eliminar un Precio
func EliminarPrecio(c *gin.Context) {
	var precio models.Precio

	idParamPrefabricada := c.Param("prefabricadaID")
	prefabricadaID, err := strconv.ParseUint(idParamPrefabricada, 10, 64)
	if err != nil {
		HandleError(c, nil, http.StatusBadRequest, "ID Prefabricada inválido")
		return
	}

	idParamPrecio := c.Param("precioID")
	precioID, err := strconv.ParseUint(idParamPrecio, 10, 64)
	if err != nil {
		HandleError(c, nil, http.StatusBadRequest, "ID Precio inválido")
		return
	}

	// Buscar el Precio en la base de datos
	if err := configs.DB.Where("prefabricada_id = ?", prefabricadaID).Where("deleted_at IS NULL").Unscoped().First(&precio, precioID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			HandleError(c, nil, http.StatusNotFound, "Precio no encontrado")
			return
		}
		HandleError(c, err, http.StatusInternalServerError, "Error al obtener Precio")
		return
	}

	// Verificar si el Precio ya se encuentra eliminado
	if precio.DeletedAt != nil && !precio.DeletedAt.IsZero() {
		HandleError(c, nil, http.StatusBadRequest, "El Precio ya se encuentra eliminado")
		return
	}

	// Poner fecha y hora de la eliminación lógica
	now := time.Now()
	precio.DeletedAt = &now

	// Guardar eliminación lógica en la base de datos
	if err := configs.DB.Save(&precio).Error; err != nil {
		HandleError(c, err, http.StatusInternalServerError, "Error, no se pudo eliminar el Precio")
		return
	}

	// Mostrar/enviar mensaje de eliminación exitosa
	c.JSON(http.StatusOK, gin.H{
		"message": "Precio eliminado exitosamente",
	})
}
