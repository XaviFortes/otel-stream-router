# Capa 1: Compilación (Ya lo haces local, pero en CI/CD se hace así)
FROM golang:1.26.4-alpine AS builder
RUN apk add --no-cache git
# Instalamos el builder de OTel
RUN go install go.opentelemetry.io/collector/cmd/builder@v0.154.0
WORKDIR /app
COPY . .
# Ejecutamos el builder dentro del contenedor
RUN env PATH="/go/bin:${PATH}" builder --config=builder-config.yaml

# Capa 2: Imagen limpia para producción
FROM alpine:3.24
RUN apk add --no-cache ca-certificates
WORKDIR /
# Copiamos solo el binario final compilado en la capa anterior
COPY --from=builder /app/dist/otel-collector-prosegur /otel-collector

# Exponemos los puertos típicos de OTel si hicieran falta
EXPOSE 4317 4318 55679

ENTRYPOINT ["/otel-collector"]