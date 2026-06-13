# Layer 1: Build the OTEL Collector using the builder
FROM golang:1.26.4-alpine AS builder
RUN apk add --no-cache git
# Install OTEL Collector Builder
RUN go install go.opentelemetry.io/collector/cmd/builder@v0.154.0
WORKDIR /app
COPY . .
# Build the OTEL Collector with the provided configuration
RUN env PATH="/go/bin:${PATH}" builder --config=builder-config.yaml

# Layer 2: Create a minimal image with the compiled binary
FROM alpine:3.24
RUN apk add --no-cache ca-certificates
WORKDIR /
# Copy the compiled OTEL Collector binary from the builder stage
COPY --from=builder /app/dist/otel-collector-prosegur /otel-collector

# Exponemos los puertos típicos de OTel si hicieran falta
EXPOSE 4317 4318 55679

ENTRYPOINT ["/otel-collector"]