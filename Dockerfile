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

# Build binary — static agar bisa jalan di scratch image
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-w -s" -o bin/api-admin-panel ./cmd/main.go

# ── Stage 2: Runner ────────────────────────────────────────────────────────────
FROM scratch

# Salin file yang dibutuhkan dari builder
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/bin/api-admin-panel /api-admin-panel

# Expose port aplikasi
EXPOSE 8080

# Jalankan binary
ENTRYPOINT ["/api-admin-panel"]
