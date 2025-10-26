
FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN apk add --no-cache tzdata

RUN CGO_ENABLED=0 GOOS=linux go build -a -o /app/main ./cmd/main.go


FROM alpine:latest

RUN apk add --no-cache tzdata

ENV TZ=Asia/Singapore

WORKDIR /app

COPY --from=builder /app/main .

EXPOSE 8080

ENTRYPOINT ["/app/main"]
