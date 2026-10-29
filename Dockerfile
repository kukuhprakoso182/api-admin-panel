# ── Stage 1: Builder ───────────────────────────────────────────────────────────
FROM golang:1.21-alpine AS builder

# Install dependency OS yang dibutuhkan
RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

# Copy dependency dulu agar layer cache optimal
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary — static agar bisa jalan di scratch image.
# TARGETARCH diisi otomatis oleh BuildKit (amd64/arm64); default amd64.
ARG TARGETOS=linux
ARG TARGETARCH=amd64
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build -ldflags="-w -s" -o bin/api-admin-panel ./cmd/main.go

# Direktori log milik user non-root (scratch tidak punya mkdir/chown)
RUN mkdir -p /out/logs

# ── Stage 2: Runner ────────────────────────────────────────────────────────────
FROM scratch

# Salin file yang dibutuhkan dari builder
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/bin/api-admin-panel /api-admin-panel
COPY --from=builder --chown=65532:65532 /out/logs /logs

# Jangan jalankan sebagai root
USER 65532:65532

# Expose port aplikasi
EXPOSE 8080

# Jalankan binary
ENTRYPOINT ["/api-admin-panel"]
