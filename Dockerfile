FROM golang:1.23-alpine AS builder

WORKDIR /app

# Salin file dependensi dulu untuk optimasi cache
COPY go.mod go.sum ./
RUN go mod download

# Salin sisa kode sumber proyek Anda
COPY . .

RUN apt-get update && apt-get install -y tzdata

# Compile aplikasi Go menjadi satu file program (binary)
# Ini adalah langkah inti dari tahap builder
RUN CGO_ENABLED=0 GOOS=linux go build -a -o /app/main ./cmd/main.go


# --- TAHAP 2: FINAL (Image Akhir yang Super Ramping) ---
# Kita mulai dari image kosong Alpine Linux yang sangat kecil
FROM alpine:latest

WORKDIR /app

# Salin HANYA file program yang sudah jadi dari tahap builder
COPY --from=builder /app/main .

# Beri tahu Docker port berapa yang digunakan aplikasi ini
EXPOSE 8080

# Perintah untuk menjalankan aplikasi saat container dimulai
ENTRYPOINT ["/app/main"]