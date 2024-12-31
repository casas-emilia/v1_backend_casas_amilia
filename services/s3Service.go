package services

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// loadAWSConfig carga la configuración de AWS desde variables de entorno
func loadAWSConfig() (aws.Config, error) {
	if os.Getenv("AWS_ACCESS_KEY_ID") == "" || os.Getenv("AWS_SECRET_ACCESS_KEY") == "" || os.Getenv("AWS_REGION") == "" {
		return aws.Config{}, fmt.Errorf("faltan variables de entorno requeridas para la configuración de AWS")
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(os.Getenv("AWS_REGION")),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			os.Getenv("AWS_ACCESS_KEY_ID"),
			os.Getenv("AWS_SECRET_ACCESS_KEY"),
			"",
		)),
	)
	if err != nil {
		return cfg, err
	}
	return cfg, nil
}

// UploadToS3 sube un archivo a S3 y devuelve la URL del archivo
func UploadToS3(file io.Reader, folder string, filename string) (string, error) {
	// Cargar la configuración de AWS
	cfg, err := loadAWSConfig()
	if err != nil {
		log.Printf("Error al cargar la configuración de AWS: %v", err)
		return "", err
	}

	// Crear el cliente de S3
	s3Client := s3.NewFromConfig(cfg)

	// Leer el archivo completo en un buffer
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, file)
	if err != nil {
		log.Printf("Error al leer el archivo: %v", err)
		return "", fmt.Errorf("no se pudo leer el archivo: %v", err)
	}

	// Obtener el tipo MIME desde los primeros 512 bytes del buffer
	contentType := http.DetectContentType(buf.Bytes()[:512])

	// Construir la clave del archivo (incluyendo la carpeta)
	key := filepath.Join(folder, filename)
	key = strings.ReplaceAll(key, "\\", "/")

	// Subir el archivo a S3
	_, err = s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String("bucket-casas-emilia"),
		Key:         aws.String(key),
		Body:        bytes.NewReader(buf.Bytes()),
		ContentType: aws.String(contentType),
	})
	if err != nil {
		log.Printf("Error al subir el archivo a S3: %v", err)
		return "", fmt.Errorf("no se pudo subir el archivo a S3: %v", err)
	}

	// Construir la URL del archivo subido
	url := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", "bucket-casas-emilia", os.Getenv("AWS_REGION"), key)
	return url, nil
}
