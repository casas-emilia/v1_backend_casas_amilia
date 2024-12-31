package middlewares

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET"))

type Claims struct {
	UsuarioID uint     `json:"usuario_id"`
	Rol       []string `json:"roles"` // Cambiado a slice para múltiples roles
	jwt.RegisteredClaims
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No se proporcionó token"})
			c.Abort()
			return
		}

		tokenString := strings.Split(authHeader, " ")
		if len(tokenString) != 2 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token malformado"})
			c.Abort()
			return
		}

		token, err := jwt.ParseWithClaims(tokenString[1], &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
			c.Abort()
			return
		}

		// Verifica si el token es válido y tiene claims
		if claims, ok := token.Claims.(*Claims); ok && token.Valid {
			if claims.ExpiresAt.After(time.Now()) {
				// Almacena usuarioID y roles en el contexto
				c.Set("usuarioID", claims.UsuarioID)
				c.Set("roles", claims.Rol)                                             // Almacena el slice de roles en el contexto
				fmt.Println("ID de usuario almacenado en contexto:", claims.UsuarioID) // Log de éxito
				fmt.Println("Roles almacenados en contexto:", claims.Rol)              // Log de éxito
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Token expirado"})
				c.Abort()
				return
			}
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No se pudo obtener el ID de usuario"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// Verificar si el usuarioID y los roles están almacenados en el contexto
func VerificarContexto(c *gin.Context) {
	usuarioID, usuarioIDExists := c.Get("usuarioID")
	roles, rolesExists := c.Get("roles")

	// Verificar si el usuarioID está en el contexto
	if usuarioIDExists {
		fmt.Printf("ID de usuario almacenado en contexto: %v\n", usuarioID)
	} else {
		fmt.Println("ID de usuario no encontrado en el contexto")
	}

	// Verificar si los roles están en el contexto
	if rolesExists {
		fmt.Printf("Roles almacenados en contexto: %v\n", roles)
	} else {
		fmt.Println("Roles no encontrados en el contexto")
	}
}
