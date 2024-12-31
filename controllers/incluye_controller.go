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

// Función para Crear un Incluye
func CrearIncluye(c *gin.Context) {
	var request dto.CrearIncluyeRequest
	var incluye models.Incluye
	var incluyeResponse dto.IncluyeResponse

	idParamPrecio := c.Param("precioID")
	precioID, err := strconv.ParseUint(idParamPrecio, 10, 64)
	if err != nil {
		HandleError(c, nil, http.StatusBadRequest, "ID Precio inválido")
		return
	}

	// Validamos el body
	if err := c.ShouldBindJSON(&request); err != nil {
		HandleError(c, nil, http.StatusBadRequest, "Error de datos")
		return
	}

	// Creamos el Incluye
	incluye.NombreIncluye = request.NombreIncluye
	incluye.PrecioID = uint(precioID)

	// Guardamos en la base de datos
	if err := configs.DB.Create(&incluye).Error; err != nil {
		HandleError(c, err, http.StatusInternalServerError, "Error, no se pudo guardar Incluye")
		return
	}

	incluyeResponse = dto.IncluyeResponse{
		ID:            incluye.ID,
		CreatedAt:     incluye.CreatedAt,
		UpdatedAt:     incluye.UpdatedAt,
		NombreIncluye: incluye.NombreIncluye,
		PrecioID:      incluye.PrecioID,
	}

	// Responder/enviar mensaje exitoso e incluye
	c.JSON(http.StatusOK, gin.H{
		"mesaage": "Incluye creado exitosamente",
		"Incluye": incluyeResponse,
	})
}

// Función para obtener todos los Incluyes
func ObtenerIncluyes(c *gin.Context) {
	var incluyes []models.Incluye
	var incluyesResponse []dto.IncluyeResponse

	idParamPrecio := c.Param("precioID")
	precioID, err := strconv.ParseUint(idParamPrecio, 10, 64)
	if err != nil {
		HandleError(c, nil, http.StatusBadRequest, "ID Precio inválido")
		return
	}

	// Buscar todos los incluyes de un precio
	if err := configs.DB.Where("precio_id = ?", precioID).Where("deleted_at IS NULL").Find(&incluyes).Error; err != nil {
		HandleError(c, err, http.StatusInternalServerError, "Error al obtener Incluyes")
		return
	}

	// Verificar si no se encontraron registros
	if len(incluyes) == 0 {
		HandleError(c, nil, http.StatusNotFound, "Incluyes no encontrados")
		return
	}

	for _, incluye := range incluyes {
		incluyesResponse = append(incluyesResponse, dto.IncluyeResponse{
			ID:            incluye.ID,
			CreatedAt:     incluye.CreatedAt,
			UpdatedAt:     incluye.UpdatedAt,
			NombreIncluye: incluye.NombreIncluye,
			PrecioID:      incluye.PrecioID,
		})
	}

	// Mostrar/enviar incluyes
	c.JSON(http.StatusOK, gin.H{
		"incluyes": incluyesResponse,
	})
}

// Función para obtener un incluye de acuerdo a su ID y precio
func ObtenerIncluye(c *gin.Context) {
	var incluye models.Incluye
	var incluyeResponse dto.IncluyeResponse

	idParamPrecio := c.Param("precioID")
	precioID, err := strconv.ParseUint(idParamPrecio, 10, 64)
	if err != nil {
		HandleError(c, nil, http.StatusBadRequest, "ID Precio inválido")
		return
	}

	idParamIncluye := c.Param("incluyeID")
	incluyeID, err := strconv.ParseUint(idParamIncluye, 10, 64)
	if err != nil {
		HandleError(c, nil, http.StatusBadRequest, "ID Incluye inválido")
		return
	}

	// Buscar Incluye en la base de datos de acuerdo a su ID y al Precio que pertenece
	if err := configs.DB.Where("precio_id = ?", precioID).Where("deleted_at IS NULL").First(&incluye, incluyeID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			HandleError(c, nil, http.StatusNotFound, "Incluye no encontrado")
			return
		}
		HandleError(c, err, http.StatusInternalServerError, "Error al obtener Incluye del Precio")
		return
	}

	incluyeResponse = dto.IncluyeResponse{
		ID:            incluye.ID,
		CreatedAt:     incluye.CreatedAt,
		UpdatedAt:     incluye.UpdatedAt,
		NombreIncluye: incluye.NombreIncluye,
		PrecioID:      incluye.PrecioID,
	}

	// Mostrar/enviar Incluye
	c.JSON(http.StatusBadRequest, gin.H{
		"incluye": incluyeResponse,
	})
}

// Función para actualizar datos de un Incluye
func ActualizarIncluye(c *gin.Context) {
	var request dto.ActualizarIncluyeRequest
	var incluye models.Incluye
	var incluyeResponse dto.IncluyeResponse

	idParamPrecio := c.Param("precioID")
	precioID, err := strconv.ParseUint(idParamPrecio, 10, 64)
	if err != nil {
		HandleError(c, nil, http.StatusBadRequest, "ID Precio inválido")
		return
	}

	idParamIncluye := c.Param("incluyeID")
	incluyeID, err := strconv.ParseUint(idParamIncluye, 10, 64)
	if err != nil {
		HandleError(c, nil, http.StatusBadRequest, "ID Incluye inválido")
		return
	}

	// Bind JSON del request a la estructura del dto
	if err := c.ShouldBindJSON(&request); err != nil {
		HandleError(c, nil, http.StatusBadRequest, "Error de datos"+err.Error())
		return
	}

	// Buscamos el Incluye en la base da datos
	if err := configs.DB.Where("precio_id", precioID).Where("deleted_at IS NULL").First(&incluye, incluyeID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			HandleError(c, nil, http.StatusNotFound, "Incluye no encontrado")
			return
		}
		HandleError(c, err, http.StatusInternalServerError, "Error al obtener Incluye")
		return
	}

	// Actualizar datos de Incluye
	incluye.NombreIncluye = request.NombreIncluye

	// Guardar los cambios en la base de datos
	if err := configs.DB.Save(&incluye).Error; err != nil {
		HandleError(c, err, http.StatusBadRequest, "Error, no se pudo actualizar los datos de Incluye")
		return
	}

	incluyeResponse = dto.IncluyeResponse{
		ID:            incluye.ID,
		CreatedAt:     incluye.CreatedAt,
		UpdatedAt:     incluye.UpdatedAt,
		NombreIncluye: incluye.NombreIncluye,
		PrecioID:      incluye.PrecioID,
	}

	// Mostrar/enviar mensaje de éxito e Incluye
	c.JSON(http.StatusOK, gin.H{
		"message": "Datos actualizados con éxito",
		"Incluye": incluyeResponse,
	})
}

// Función para eliminar lógicamente un Incluye
func EliminarIncluye(c *gin.Context) {
	var incluye models.Incluye

	idParamPrecio := c.Param("precioID")
	precioID, err := strconv.ParseUint(idParamPrecio, 10, 64)
	if err != nil {
		HandleError(c, nil, http.StatusBadRequest, "ID Precio inválido")
		return
	}

	idParamIncluye := c.Param("incluyeID")
	incluyeID, err := strconv.ParseUint(idParamIncluye, 10, 64)
	if err != nil {
		HandleError(c, nil, http.StatusBadRequest, "ID Incluye inválido")
		return
	}

	// Buscar Incluye en la Base de datos
	if err := configs.DB.Where("precio_id = ?", precioID).Where("deleted_at IS NULL").First(&incluye, incluyeID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			HandleError(c, nil, http.StatusNotFound, "Incluye no encontrado")
			return
		}
		HandleError(c, err, http.StatusInternalServerError, "Error al obtener Incluye")
		return
	}

	// Verificamos si el Incluye ya se encuentra eliminado
	if incluye.DeletedAt != nil && !incluye.DeletedAt.IsZero() {
		HandleError(c, nil, http.StatusBadRequest, "El Incluye ya se encuentra eliminado")
		return
	}

	// Poner fecha y hora de la eliminación lógica
	now := time.Now()
	incluye.DeletedAt = &now

	// Guardar eliminación lógica en la base de datos
	if err := configs.DB.Save(&incluye).Error; err != nil {
		HandleError(c, err, http.StatusInternalServerError, "Error, no se pudo eliminar el Incluye")
		return
	}

	// Mostrar/enviar mensaje de eliminación lógica exitosa
	c.JSON(http.StatusOK, gin.H{
		"message": "Incluye eliminado exitosamente",
	})
}
