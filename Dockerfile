FROM golang:1.22-alpine AS builder

# Menetapkan direktori kerja di dalam container.
WORKDIR /app

# `go mod download` hanya akan berjalan lagi jika file go.mod/go.sum berubah.
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Hasilnya adalah sebuah file bernama 'main' di dalam folder /app.
RUN CGO_ENABLED=0 GOOS=linux go build -v -o /app/main .


# Memulai dari image Alpine Linux yang sangat kecil untuk produksi.
FROM alpine:latest

# Menetapkan direktori kerja di dalam image akhir.
WORKDIR /root/

# Menyalin HANYA file binary yang sudah di-compile dari tahap 'builder'.
# Ini membuat image akhir sangat kecil dan aman.
COPY --from=builder /app/main .

# Memberi tahu Docker bahwa container akan mendengarkan pada port 10000.
# Platform seperti Render akan menggunakan informasi ini.
EXPOSE 10000

# Perintah default untuk menjalankan aplikasi Anda saat container dimulai.
CMD ["./main"]