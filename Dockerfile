
FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# ✅ Install tzdata untuk timezone di build stage (optional, tapi tidak perlu jika tidak dijalankan di sini)
RUN apk add --no-cache tzdata

RUN CGO_ENABLED=0 GOOS=linux go build -a -o /app/main ./cmd/main.go


# --- TAHAP 2: FINAL ---
FROM alpine:latest

# ✅ Install tzdata agar timezone seperti Asia/Singapore dikenali
RUN apk add --no-cache tzdata

# ✅ Set default timezone (optional tapi baik jika diperlukan)
ENV TZ=Asia/Singapore

WORKDIR /app

COPY --from=builder /app/main .

EXPOSE 8080

ENTRYPOINT ["/app/main"]
