package controllers

import (
	"net/http"
	"os"
	"time"
	"v1_prefabricadas/configs"
	"v1_prefabricadas/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET")) // Lee la clave secreta desde variables de entorno

// Estructura para las claims del token JWT
type Claims struct {
	UsuarioID uint     `json:"usuario_id"`
	Roles     []string `json:"roles"` // Campo para los roles
	jwt.RegisteredClaims
}

// Login: Verificar credenciales y generar un JWT con múltiples roles
func Login(c *gin.Context) {
	var credencialRequest models.Credencial

	// Parsear las credenciales desde el request
	if err := c.ShouldBindJSON(&credencialRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var usuario models.Usuario
	// Buscar el usuario por email y precargar las relaciones de Credencial y Rol_usuario con Rol
	if err := configs.DB.
		Joins("JOIN credenciales ON credenciales.usuario_id = usuarios.id").
		Where("credenciales.email = ?", credencialRequest.Email).
		Preload("Credencial").      // Precargar Credencial
		Preload("Rol_usuario.Rol"). // Precargar Rol_usuario y Rol
		First(&usuario).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciales inválidas"})
		return
	}

	// Verificar si la contraseña es correcta
	err := bcrypt.CompareHashAndPassword([]byte(usuario.Credencial.Password), []byte(credencialRequest.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciales inválidas"})
		return
	}

	// Extraer los nombres de los roles en un slice
	var roles []string
	for _, rolUsuario := range usuario.Rol_usuario {
		roles = append(roles, rolUsuario.Rol.NombreRol)
	}

	// Generar el token JWT con usuarioID y roles
	tokenString, err := generarJWT(usuario.ID, roles)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo generar el token"})
		return
	}

	// Devolver el token al cliente
	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

// Generar un JWT con una duración de 24 horas
func generarJWT(usuarioID uint, roles []string) (string, error) {
	claims := &Claims{
		UsuarioID: usuarioID,
		Roles:     roles, // Añadir roles a las claims
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			Issuer:    "miApp",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
