FROM golang:1.24.5-alpine AS builder

# Instalar dependencias del sistema si son necesarias
RUN apk --no-cache add ca-certificates git

# Crear directorio de trabajo
WORKDIR /app

# Copiar go.mod y go.sum primero (mejor cache de layers)
COPY go.mod go.sum ./

# Descargar dependencias
RUN go mod download

# Copiar código fuente
COPY . .

# Compilar la aplicación
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main cmd/api/main.go

# Stage 2: Runtime (imagen final más pequeña)
FROM alpine:latest

# Instalar ca-certificates para HTTPS
RUN apk --no-cache add ca-certificates tzdata

# Crear usuario no-root por seguridad
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Crear directorio de trabajo
WORKDIR /app

# Copiar binario desde builder
COPY --from=builder /app/main .

# Cambiar ownership al usuario no-root
RUN chown -R appuser:appgroup /app

# Usar usuario no-root
USER appuser

# Exponer puerto
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Comando por defecto
CMD ["./main"]