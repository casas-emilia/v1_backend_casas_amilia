package configs

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDB() {
	var err error

	// Obtener las variables de entorno para la base de datos
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// Verificar si todas las variables de entorno están configuradas
	if dbUser == "" || dbPassword == "" || dbHost == "" || dbPort == "" || dbName == "" {
		log.Fatal("Faltan variables de entorno para la conexión a la base de datos")
	}

	// Crear el DSN (Data Source Name) utilizando las variables de entorno
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPassword, dbHost, dbPort, dbName)

	// Conectar a la base de datos
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect DB:", err)
	}

	log.Println("Conexión a la base de datos exitosa")
}
