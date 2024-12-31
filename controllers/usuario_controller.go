package controllers

import (
	"errors"
	"net/http"
	"strconv"
	"time"
	"v1_prefabricadas/configs"
	"v1_prefabricadas/dto"
	"v1_prefabricadas/models"
	"v1_prefabricadas/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Función para crear un Usuario
func CrearUsuario(c *gin.Context) {
	var request dto.CrearUsuarioRequest
	var usuario models.Usuario
	var usuarioResponse dto.UsuarioResponse

	idParamEmpresa := c.Param("empresaID")
	empresaID, err := strconv.ParseUint(idParamEmpresa, 10, 64)
	if err != nil {
		HandleError(c, nil, http.StatusBadRequest, "ID Empresa inválido")
		return
	}

	//Validamos el body
	// usar `c.ShouldBind()` en lugar de `c.ShouldBindJSON()`. `ShouldBind()`
	// es más flexible y puede manejar diferentes tipos de contenido, incluyendo formularios multipart.

	if err := c.ShouldBind(&request); err != nil {
		HandleError(c, err, http.StatusBadRequest, "Error en los datos del formulario")
		return
	}

	// Manejo de la imagen
	fileHeader, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Image is required", "details": err.Error()})
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open image", "details": err.Error()})
		return
	}
	defer file.Close()

	// Subir la imagen a AWS S3
	url, err := services.UploadToS3(file, "imagenes_usuarios", fileHeader.Filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload image to S3", "details": err.Error()})
		return
	}

	// Creamos los datos
	usuario.PrimerNombre = request.PrimerNombre
	usuario.SegundoNombre = request.SegundoNombre
	usuario.PrimerApellido = request.PrimerApellido
	usuario.SegundoApellido = request.SegundoApellido
	usuario.Image = url
	usuario.EmpresaID = uint(empresaID)

	// Guardar Usuario en la base de datos
	if err := configs.DB.Create(&usuario).Error; err != nil {
		HandleError(c, err, http.StatusInternalServerError, "Error, no se pudo crear el Usuario")
		return
	}

	usuarioResponse = dto.UsuarioResponse{
		ID:              usuario.ID,
		CreatedAt:       usuario.CreatedAt,
		UpdatedAt:       usuario.UpdatedAt,
		PrimerNombre:    usuario.PrimerNombre,
		SegundoNombre:   usuario.SegundoNombre,
		PrimerApellido:  usuario.PrimerApellido,
		SegundoApellido: usuario.SegundoApellido,
		Image:           usuario.Image,
		EmpresaID:       usuario.EmpresaID,
	}

	// Mostrar/enviar mensaje exitoso y el Usuario
	c.JSON(http.StatusOK, gin.H{
		"message": "Usuario creado con éxito",
		"usuario": usuarioResponse,
	})
}

// Función para obtener todos los usuarios de una empresa
func ObtenerUsuarios(c *gin.Context) {
	var usuarios []models.Usuario
	var usuariosResponse []dto.UsuarioResponse

	idParamEmpresa := c.Param("empresaID")
	empresaID, err := strconv.ParseUint(idParamEmpresa, 10, 64)
	if err != nil {
		HandleError(c, nil, http.StatusBadRequest, "ID Empresa inválido")
		return
	}

	// Buscar todos los usuarios de una empresa
	if err := configs.DB.Where("empresa_id = ?", empresaID).Where("deleted_at IS NULL").Find(&usuarios).Error; err != nil {
		HandleError(c, err, http.StatusInternalServerError, "Error, no se pudo obtener a los Usuarios")
		return
	}

	// Controllar si no se encuentran registros
	if len(usuarios) == 0 {
		HandleError(c, nil, http.StatusNotFound, "Sin Usuarios guardados")
		return
	}

	for _, usuario := range usuarios {
		usuariosResponse = append(usuariosResponse, dto.UsuarioResponse{
			ID:              usuario.ID,
			CreatedAt:       usuario.CreatedAt,
			UpdatedAt:       usuario.UpdatedAt,
			PrimerNombre:    usuario.PrimerNombre,
			SegundoNombre:   usuario.SegundoNombre,
			PrimerApellido:  usuario.PrimerApellido,
			SegundoApellido: usuario.SegundoApellido,
			Image:           usuario.Image,
			EmpresaID:       usuario.EmpresaID,
		})
	}

	// Mostrar/enviar los Usuarios
	c.JSON(http.StatusOK, gin.H{"usuarios": usuariosResponse})
}

// Función para obtener un Usuario de acuerdo a su ID
func ObtenerUsuario(c *gin.Context) {
	var usuario models.Usuario
	var usuarioResponse dto.UsuarioResponse

	idParamEmpresa := c.Param("empresaID")
	empresaID, err := strconv.ParseUint(idParamEmpresa, 10, 64)
	if err != nil {
		HandleError(c, nil, http.StatusBadRequest, "ID Empresa inválido")
		return
	}

	idParamUsuario := c.Param("usuarioID")
	usuarioID, err := strconv.ParseUint(idParamUsuario, 10, 64)
	if err != nil {
		HandleError(c, nil, http.StatusBadRequest, "ID Usuario inválido")
		return
	}

	// Buscar Usuario en la Base de datos de acuerdo a su ID
	if err := configs.DB.Where("empresa_id = ?", empresaID).Where("deleted_at IS NULL").First(&usuario, usuarioID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			HandleError(c, nil, http.StatusNotFound, "Usuario no encontrado")
			return
		}
		HandleError(c, err, http.StatusInternalServerError, "Error al obtener Usuario")
		return
	}

	usuarioResponse = dto.UsuarioResponse{
		ID:              usuario.ID,
		CreatedAt:       usuario.CreatedAt,
		UpdatedAt:       usuario.UpdatedAt,
		PrimerNombre:    usuario.PrimerNombre,
		SegundoNombre:   usuario.SegundoNombre,
		PrimerApellido:  usuario.PrimerApellido,
		SegundoApellido: usuario.SegundoApellido,
		Image:           usuario.Image,
		EmpresaID:       usuario.EmpresaID,
	}

	// Mostrar/enviar Usuario
	c.JSON(http.StatusOK, gin.H{"usuario": usuarioResponse})
}

// Función para actualizar datos de Usuario
func ActualizarUsuario(c *gin.Context) {
	var request dto.ActualizarUsuarioRequest
	var usuario models.Usuario
	var usuarioResponse dto.UsuarioResponse

	idParamEmpresa := c.Param("empresaID")
	empresaID, err := strconv.ParseUint(idParamEmpresa, 10, 64)
	if err != nil {
		HandleError(c, nil, http.StatusBadRequest, "ID Empresa inválido")
		return
	}

	idParamUsuario := c.Param("usuarioID")
	usuarioID, err := strconv.ParseUint(idParamUsuario, 10, 64)
	if err != nil {
		HandleError(c, nil, http.StatusBadRequest, "ID Usuario inválido")
		return
	}

	// Bind JSON del request a la estructura de dto
	if err := c.ShouldBind(&request); err != nil {
		HandleError(c, err, http.StatusBadRequest, "Error en los datos del formulario")
		return
	}

	// Buscar datos de Usuario en la base de datos de acuerdo a su ID
	if err := configs.DB.Where("empresa_id = ?", empresaID).Where("deleted_at IS NULL").First(&usuario, usuarioID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			HandleError(c, nil, http.StatusNotFound, "Usuario no encontrado")
			return
		}
		HandleError(c, err, http.StatusInternalServerError, "Error al obtener Usuario")
		return
	}

	// Si se proporciona una nueva imagen, subirla a S3
	if fileHeader, err := c.FormFile("image"); err == nil {
		// Si hay una nueva imagen, subirla a S3
		file, err := fileHeader.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open image", "details": err.Error()})
			return
		}
		defer file.Close()

		// Subir la nueva imagen a S3
		url, err := services.UploadToS3(file, "imagenes_usuarios", fileHeader.Filename)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload image to S3", "details": err.Error()})
			return
		}
		// Actualizar la URL de la imagen
		usuario.Image = url
	}

	// Actulizar datos de Usuario
	usuario.PrimerNombre = request.PrimerNombre
	usuario.SegundoNombre = request.SegundoNombre
	usuario.PrimerApellido = request.PrimerApellido
	usuario.SegundoApellido = request.SegundoApellido

	// Guradar en la Base de datos
	if err := configs.DB.Save(&usuario).Error; err != nil {
		HandleError(c, err, http.StatusInternalServerError, "Error, no se pudo actualizar los datos de Usuario")
		return
	}

	usuarioResponse = dto.UsuarioResponse{
		ID:              usuario.ID,
		CreatedAt:       usuario.CreatedAt,
		UpdatedAt:       usuario.UpdatedAt,
		PrimerNombre:    usuario.PrimerNombre,
		SegundoNombre:   usuario.SegundoNombre,
		PrimerApellido:  usuario.PrimerApellido,
		SegundoApellido: usuario.SegundoApellido,
		Image:           usuario.Image,
		EmpresaID:       usuario.EmpresaID,
	}

	// Mostrar/enviar mensaje exitoso y el usuarioResponse
	c.JSON(http.StatusOK, gin.H{
		"message": "Datos de Usuario actualizados con éxito",
		"usuario": usuarioResponse,
	})
}

// Función para eliminar lógicamente un Usuario
func EliminarUsuario(c *gin.Context) {
	var usuario models.Usuario

	idParamEmpresa := c.Param("empresaID")
	empresaID, err := strconv.ParseUint(idParamEmpresa, 10, 64)
	if err != nil {
		HandleError(c, nil, http.StatusBadRequest, "ID Empresa inválido")
		return
	}

	idParamUsuario := c.Param("usuarioID")
	usuarioID, err := strconv.ParseUint(idParamUsuario, 10, 64)
	if err != nil {
		HandleError(c, nil, http.StatusBadRequest, "ID Usuario inválido")
		return
	}

	// Buscar Usuario en la base de datos de acuerdo a su ID
	if err := configs.DB.Where("empresa_id = ?", empresaID).Where("deleted_at IS NULL").Unscoped().First(&usuario, usuarioID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			HandleError(c, nil, http.StatusNotFound, "Usuario no encontrado")
			return
		}
		HandleError(c, err, http.StatusInternalServerError, "Error al Obtener Usuario")
		return
	}

	// Verificar si el Usuario ya se encuentra eliminado
	if usuario.DeletedAt != nil && !usuario.DeletedAt.IsZero() {
		HandleError(c, nil, http.StatusBadRequest, "El Usuario ya se encuentra eliminado")
		return
	}

	// Poner fecha y hora de la eliminación lógica
	now := time.Now()
	usuario.DeletedAt = &now

	// Guardar en la base de datos la fecha y hora de la eliminación lógica
	if err := configs.DB.Save(&usuario).Error; err != nil {
		HandleError(c, err, http.StatusInternalServerError, "Error, no se pudo eliminar al usuario")
		return
	}

	// Mostrar/enviar un mensaje de eliminación exitosa
	c.JSON(http.StatusOK, gin.H{
		"message": "Usuario eliminado exitosamente",
	})
}
