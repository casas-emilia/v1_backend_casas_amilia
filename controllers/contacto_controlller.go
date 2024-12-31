package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
	"v1_prefabricadas/configs"
	"v1_prefabricadas/dto"
	"v1_prefabricadas/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Función para crear datos de Contacto de un Usuario
func CrearContacto(c *gin.Context) {
	var request dto.CrearContactoRequest
	var contacto models.Contacto
	var contactoResponse dto.ContactoResponse

	idParamUsuario := c.Param("usuarioID")
	usuarioID, err := strconv.ParseUint(idParamUsuario, 10, 64)
	if err != nil {
		HandleError(c, nil, http.StatusBadRequest, "ID Usuario inválido")
		return
	}

	// Validamos de body
	if err := c.ShouldBindJSON(&request); err != nil {
		HandleError(c, nil, http.StatusBadRequest, "Error de datos")
		return
	}

	// Creamos los datos de Contacto
	contacto.EmailLaboral = request.EmailLaboral
	contacto.CelularLaboral = request.CelularLaboral
	contacto.DireccionLaboral = request.DireccionLaboral
	contacto.UsuarioID = uint(usuarioID)

	// Guardar datos de Contacto en la base de datos
	if err := configs.DB.Create(&contacto).Error; err != nil {
		HandleError(c, err, http.StatusInternalServerError, "Error, no se pudo guardar los datos de Contacto")
		return
	}

	contactoResponse = dto.ContactoResponse{
		ID:               contacto.ID,
		CreatedAt:        contacto.CreatedAt,
		UpdatedAt:        contacto.UpdatedAt,
		EmailLaboral:     contacto.EmailLaboral,
		CelularLaboral:   contacto.CelularLaboral,
		DireccionLaboral: contacto.DireccionLaboral,
		UsuarioID:        contacto.UsuarioID,
	}

	// Mostrar/enviar mensaje exitoso y datos de Contacto credo
	c.JSON(http.StatusOK, gin.H{
		"message":  "Datos de Contacto creados exitosamente",
		"contacto": contactoResponse,
	})

}

// Función para obtener todos los contactos de un Usuario
func ObtenerContactos(c *gin.Context) {
	var contactos []models.Contacto
	var contactoResponse []dto.ContactoResponse

	idParamUsuario := c.Param("usuarioID")
	usuarioID, err := strconv.ParseUint(idParamUsuario, 10, 64)
	if err != nil {
		HandleError(c, nil, http.StatusBadRequest, "ID Usuario inválido")
		return
	}

	// Buscar los Contactos del Usuario en la base de datos de acuerdo al ID del Usuario
	if err := configs.DB.Where("usuario_id = ?", usuarioID).Where("deleted_at IS NULL").Find(&contactos).Error; err != nil {
		HandleError(c, err, http.StatusInternalServerError, "Error al obtener datos de Contacto")
		return
	}

	// Devolver notFound en caso de que no existan Contactos
	if len(contactos) == 0 {
		HandleError(c, nil, http.StatusNotFound, "Datos de Contacto no encontrados")
		return
	}

	for _, contacto := range contactos {
		contactoResponse = append(contactoResponse, dto.ContactoResponse{
			ID:               contacto.ID,
			CreatedAt:        contacto.CreatedAt,
			UpdatedAt:        contacto.UpdatedAt,
			EmailLaboral:     contacto.EmailLaboral,
			CelularLaboral:   contacto.CelularLaboral,
			DireccionLaboral: contacto.DireccionLaboral,
			UsuarioID:        contacto.UsuarioID,
		})
	}

	// Mostrar/enviar Contactos de Usuario
	c.JSON(http.StatusOK, gin.H{
		"contactos": contactoResponse,
	})

}

// Función para obtener Datos de Contacto de acuerdo a su ID
func ObtenerContacto(c *gin.Context) {
	var contacto models.Contacto
	var contactoResponse dto.ContactoResponse

	idParamUsuario := c.Param("usuarioID")
	usuarioID, err := strconv.ParseUint(idParamUsuario, 10, 64)
	if err != nil {
		HandleError(c, nil, http.StatusBadRequest, "ID Usuario inválido")
		return
	}

	idParamContacto := c.Param("contactoID")
	contactoID, err := strconv.ParseUint(idParamContacto, 10, 64)
	if err != nil {
		HandleError(c, nil, http.StatusBadRequest, "ID Contacto inválido")
		return
	}

	// Buscar Datos de Contacto en la Base de datos de acuerdo a su ID
	if err := configs.DB.Where("usuario_id = ?", usuarioID).Where("deleted_at IS NULL").First(&contacto, contactoID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			HandleError(c, nil, http.StatusNotFound, "Datos de Contacto no encontrados")
			return
		}
		HandleError(c, err, http.StatusInternalServerError, "Error al obtener datos de Contacto")
		return
	}

	contactoResponse = dto.ContactoResponse{
		ID:               contacto.ID,
		CreatedAt:        contacto.CreatedAt,
		UpdatedAt:        contacto.UpdatedAt,
		EmailLaboral:     contacto.EmailLaboral,
		CelularLaboral:   contacto.CelularLaboral,
		DireccionLaboral: contacto.DireccionLaboral,
		UsuarioID:        contacto.UsuarioID,
	}

	// Mostrar/enviar datos de Contacto
	c.JSON(http.StatusOK, gin.H{
		"contacto": contactoResponse,
	})
}

// Función para Actualizar datos de Contacto
func ActualizarContacto(c *gin.Context) {
	var request dto.ActualizarContactoRequest
	var contacto models.Contacto
	var contactoResponse dto.ContactoResponse

	idParamUsuario := c.Param("usuarioID")
	usuarioID, err := strconv.ParseUint(idParamUsuario, 10, 64)
	if err != nil {
		HandleError(c, nil, http.StatusBadRequest, "ID Usuario inválido")
		return
	}

	idParamContacto := c.Param("contactoID")
	contactoID, err := strconv.ParseUint(idParamContacto, 10, 64)
	if err != nil {
		HandleError(c, nil, http.StatusBadRequest, "ID Contacto inválido")
		return
	}

	// Bind JSON del request a la estructura del DTO
	if err := c.ShouldBindJSON(&request); err != nil {
		HandleError(c, nil, http.StatusBadRequest, "Error de Datos")
		return
	}

	// Buscar Contacto en la base de datos
	if err := configs.DB.Where("usuario_id = ?", usuarioID).Where("deleted_at IS NULL").First(&contacto, contactoID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			HandleError(c, nil, http.StatusBadRequest, "Contacto no encontrado")
			return
		}
		HandleError(c, err, http.StatusInternalServerError, "Error, no se pudo obtener datos de Contacto")
		return
	}

	// Actualizar datos de Contacto
	contacto.EmailLaboral = request.EmailLaboral
	contacto.CelularLaboral = request.CelularLaboral
	contacto.DireccionLaboral = request.DireccionLaboral

	// Guardar los datos actualizados en la Base de datos
	if err := configs.DB.Save(&contacto).Error; err != nil {
		HandleError(c, err, http.StatusInternalServerError, "Error, no se pudo actualizar los datos de Contacto")
		return
	}

	contactoResponse = dto.ContactoResponse{
		ID:               contacto.ID,
		CreatedAt:        contacto.CreatedAt,
		UpdatedAt:        contacto.UpdatedAt,
		EmailLaboral:     contacto.EmailLaboral,
		CelularLaboral:   contacto.CelularLaboral,
		DireccionLaboral: contacto.DireccionLaboral,
		UsuarioID:        contacto.UsuarioID,
	}

	// Mostrar/enviar contacto y mensaje de éxito
	c.JSON(http.StatusOK, gin.H{
		"message":  "Datos de Contacto actualizado exitosamente",
		"contacto": contactoResponse,
	})

}

// Función para eliminar lógicamente datos de Contacto
func EliminarContacto(c *gin.Context) {
	var contacto models.Contacto

	idParamUsuario := c.Param("usuarioID")
	usuarioID, err := strconv.ParseUint(idParamUsuario, 10, 64)
	if err != nil {
		HandleError(c, nil, http.StatusBadRequest, "ID Usuario inválido")
		return
	}

	idParamContacto := c.Param("contactoID")
	contactoID, err := strconv.ParseUint(idParamContacto, 10, 64)
	if err != nil {
		HandleError(c, nil, http.StatusBadRequest, "ID Contacto inválido")
		return
	}

	// Buscar datos de Contacto en la base de datos de acuerdo a su ID
	if err := configs.DB.Where("usuario_id = ? AND id = ?", usuarioID, contactoID).Where("deleted_at IS NULL").First(&contacto).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			HandleError(c, nil, http.StatusNotFound, "Datos de Contacto no encontrados")
			return
		}
		HandleError(c, err, http.StatusInternalServerError, "Error al obtener los Datos de Contacto")
		return
	}

	// Verificar si los Datos de Contacto ya estan eliminado lógicamente
	if contacto.DeletedAt != nil && !contacto.DeletedAt.IsZero() {
		HandleError(c, nil, http.StatusBadRequest, "Los Datos de Contacto ya han sido eliminados")
		return
	}

	// Poner fecha y hora de la eliminación lógica
	now := time.Now()
	contacto.DeletedAt = &now

	// Guardar eliminación lógica de la eliminación lógica
	if err := configs.DB.Save(&contacto).Error; err != nil {
		HandleError(c, err, http.StatusInternalServerError, "Error, no se pudo eliminar Datos de Contacto")
		return
	}

	// Mostrar/enviar mensaje de eliminación éxitosa
	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Contacto con ID %d eliminado exitosamente", contactoID),
	})
}
