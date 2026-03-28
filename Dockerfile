FROM golang:1.22-alpine AS builder
WORKDIR /app

# Устанавливаем необходимые пакеты для сборки с CGO (нужен для sqlite3, хотя мы скоро перейдём на Postgres)
RUN apk add --no-cache gcc musl-dev

COPY go.mod go.sum ./
RUN go mod download

# Копируем ВСЕ файлы проекта
COPY . .

# Собираем проект правильно (все .go файлы в пакете main)
RUN CGO_ENABLED=1 GOOS=linux go build -o main .

FROM alpine:latest
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app

COPY --from=builder /app/main .
EXPOSE 8080

CMD ["./main"]