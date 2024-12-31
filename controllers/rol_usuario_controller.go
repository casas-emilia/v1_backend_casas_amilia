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

// Función para crear Roles de Usuario
func CrearRol_usuario(c *gin.Context) {
	var request dto.CrearRol_usuarioRequest

	idParamRol := c.Param("rolID")
	rolID, err := strconv.ParseUint(idParamRol, 10, 64)
	if err != nil {
		HandleError(c, nil, http.StatusBadRequest, "ID Rol inválido")
		return
	}

	// Validamos el JSON entrante a la estructura DTO
	if err := c.ShouldBindJSON(&request); err != nil {
		HandleError(c, err, http.StatusBadRequest, err.Error())
		return
	}

	// Creamos Rol_usuario
	rol_usuario := models.Rol_usuario{
		UsuarioID: request.UsuarioID,
		RolID:     uint(rolID),
	}

	// Guardamos en la base de datos en nuevo Rol de Usuario
	if err := configs.DB.Create(&rol_usuario).Error; err != nil {
		HandleError(c, err, http.StatusInternalServerError, "No se pudo guardar el Rol de Usuario")
		return
	}

	// Mostramos un mensaje de éxito
	c.JSON(http.StatusOK, gin.H{"message": "Rol_usuario creado con éxito"})

	// Response
	response := dto.Rol_usuarioResponse{
		ID:        rol_usuario.ID,
		CreatedAt: rol_usuario.CreatedAt,
		UpdatedAt: rol_usuario.UpdatedAt,
		UsuarioID: rol_usuario.UsuarioID,
		RolID:     rol_usuario.RolID,
	}

	// Mostramos el Response
	c.JSON(http.StatusCreated, gin.H{"rol_usuario": response})
}

// Función para obtener todos los usuarios de un rol
func ObtenerRoles_usuarios(c *gin.Context) {
	var roles_usuarios []models.Rol_usuario
	var roles_usuariosResponse []dto.Rol_usuarioResponse

	idParamRol := c.Param("rolID")
	rolID, err := strconv.ParseUint(idParamRol, 10, 64)
	if err != nil {
		HandleError(c, err, http.StatusBadRequest, "ID Rol inválido")
		return
	}

	// Buscar todos los usuarios de todos los roles
	if err := configs.DB.Where("deleted_at IS NULL").Where("rol_id = ?", rolID).Find(&roles_usuarios).Error; err != nil {
		HandleError(c, err, http.StatusInternalServerError, "Error, no se pudo obtener los Datos solicitados(Roles de usuarios)")
		return
	}

	if len(roles_usuarios) == 0 {
		HandleError(c, nil, http.StatusNotFound, "Datos no encontrados")
		return
	}

	for _, rol_usuario := range roles_usuarios {
		roles_usuariosResponse = append(roles_usuariosResponse, dto.Rol_usuarioResponse{
			ID:        rol_usuario.ID,
			CreatedAt: rol_usuario.CreatedAt,
			UpdatedAt: rol_usuario.UpdatedAt,
			UsuarioID: rol_usuario.UsuarioID,
			RolID:     rol_usuario.RolID,
		})
	}

	// Mostrar/enviar roles_usuarios
	c.JSON(http.StatusOK, gin.H{
		"roles_usuarios": roles_usuariosResponse,
	})

}

// Función para Obtener un Usuario de un Rol
func ObtenerRol_usuario(c *gin.Context) {
	var rol_usuario models.Rol_usuario
	var rol_usuarioResponse dto.Rol_usuarioResponse

	idParamRol := c.Param("rolID")
	rolID, err := strconv.ParseUint(idParamRol, 10, 64)
	if err != nil {
		HandleError(c, err, http.StatusBadRequest, "ID Rol inválido")
		return
	}

	idParamUsuario := c.Param("rol_usuarioID")
	rol_usuarioID, err := strconv.ParseUint(idParamUsuario, 10, 64)
	if err != nil {
		HandleError(c, nil, http.StatusBadRequest, "ID rol_usuario inválido")
		return
	}

	// Buscar el Usuario y el Rol en la base de Datos
	if err := configs.DB.Where("deleted_at IS NULL").Where("rol_id = ? AND id = ?", rolID, rol_usuarioID).First(&rol_usuario).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			HandleError(c, nil, http.StatusNotFound, "Rol_usuario no encontrado")
			return
		}
		HandleError(c, err, http.StatusInternalServerError, "Error al tratar de obtener datos, intente nuevamente más tarde")
		return
	}

	rol_usuarioResponse = dto.Rol_usuarioResponse{
		ID:        rol_usuario.ID,
		CreatedAt: rol_usuario.CreatedAt,
		UpdatedAt: rol_usuario.UpdatedAt,
		UsuarioID: rol_usuario.UsuarioID,
		RolID:     rol_usuario.RolID,
	}

	// Mostrar/enviar rol_usuarioResponse
	c.JSON(http.StatusOK, gin.H{
		"roles_usuarios": rol_usuarioResponse,
	})

}

// Función para actualizar datos de roles_usuarios
func ActualizarRol_usuario(c *gin.Context) {
	var request dto.ActualizarRol_usuarioRequest
	var rol_usuario models.Rol_usuario
	var rol_usuarioResponse dto.Rol_usuarioResponse

	idParamRol := c.Param("rolID")
	rolID, err := strconv.ParseUint(idParamRol, 10, 64)
	if err != nil {
		HandleError(c, err, http.StatusBadRequest, "ID Rol inválido")
		return
	}

	idParamRol_usuarioID := c.Param("rol_usuarioID")
	rol_usuarioID, err := strconv.ParseUint(idParamRol_usuarioID, 10, 64)
	if err != nil {
		HandleError(c, err, http.StatusBadRequest, "ID Rol_usuario inválido")
		return
	}

	// Bind JSON del request a la estructura del DTO
	if err := c.ShouldBindJSON(&request); err != nil {
		HandleError(c, err, http.StatusBadRequest, "Error de Datos "+err.Error())
		return
	}

	// Buscar en la base de datos
	if err := configs.DB.Where("deleted_at IS NULL").Where("rol_id = ? AND id = ?", rolID, rol_usuarioID).First(&rol_usuario).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			HandleError(c, nil, http.StatusNotFound, "Rol_usuario no encontrado")
			return
		}
		HandleError(c, err, http.StatusInternalServerError, "Error al tratar de obtener los datos")
		return
	}

	// actualizar los datos
	rol_usuario.UsuarioID = request.UsuarioID
	rol_usuario.RolID = uint(rolID)

	// Guardar los cambios en la base de datos
	if err := configs.DB.Save(&rol_usuario).Error; err != nil {
		HandleError(c, err, http.StatusInternalServerError, "Error al actualizar datos, intente nuevamente")
		return
	}

	rol_usuarioResponse = dto.Rol_usuarioResponse{
		ID:        rol_usuario.ID,
		CreatedAt: rol_usuario.CreatedAt,
		UpdatedAt: rol_usuario.UpdatedAt,
		RolID:     rol_usuario.RolID,
		UsuarioID: rol_usuario.UsuarioID,
	}

	// Mostrar/enviar mensaje de actualización exitosa y rol_response
	c.JSON(http.StatusOK, gin.H{
		"message":        "Datos actualizados con éxito",
		"roles_usuarios": rol_usuarioResponse,
	})

}

// Función para eliminar lógicamente una relación rol_usuario
func EliminarRol_usuario(c *gin.Context) {
	var rol_usuario models.Rol_usuario

	idParamRol := c.Param("rolID")
	rolID, err := strconv.ParseUint(idParamRol, 10, 64)
	if err != nil {
		HandleError(c, err, http.StatusBadRequest, "ID Rol inválido")
		return
	}

	idParamRol_usuarioID := c.Param("rol_usuarioID")
	rol_usuarioID, err := strconv.ParseUint(idParamRol_usuarioID, 10, 64)
	if err != nil {
		HandleError(c, err, http.StatusBadRequest, "ID Rol_usuario inválido")
		return
	}

	// Buscar en la base de datos
	if err := configs.DB.Where("deleted_at IS NULL").Where("rol_id = ? AND id = ?", rolID, rol_usuarioID).First(&rol_usuario).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			HandleError(c, nil, http.StatusNotFound, "Rol_usuario no encontrado")
			return
		}
		HandleError(c, err, http.StatusInternalServerError, "Error al tratar de obtener los datos")
		return
	}

	// Poner fecha y hora de la eliminación lógica
	now := time.Now()
	rol_usuario.DeletedAt = &now

	// Guardar en la base de datos la fecha y hora de la eliminación lógica
	if err := configs.DB.Save(&rol_usuario).Error; err != nil {
		HandleError(c, err, http.StatusInternalServerError, "Error al tratar de eliminar Rol_usuario")
		return
	}

	// Mostrar/enviar mensaje de eliminación éxitosa
	c.JSON(http.StatusOK, gin.H{
		"message": "Rol_usuario eliminado exitosamente",
	})

}
