# Etapa 1: Construcción de la aplicación Go
FROM golang:1.23 AS builder

WORKDIR /app

# Copia el código fuente y descarga las dependencias
COPY . .

# Compila el binario
RUN go mod tidy
RUN go build -o main .

# Etapa 2: Crear la imagen final basada en Alpine
FROM alpine:latest

WORKDIR /root/

# Copia los certificados SSL (si es necesario)
RUN apk --no-cache add ca-certificates

# Copia el binario compilado desde la etapa anterior
COPY --from=builder /app/main .

# Expón el puerto en el que corre la aplicación
EXPOSE 8080

# Comando para ejecutar la aplicación
CMD ["./main"]
