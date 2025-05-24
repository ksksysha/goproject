# Этап сборки
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Копируем файлы зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main.go

# Финальный этап
FROM alpine:latest

WORKDIR /app

# Копируем бинарный файл из этапа сборки
COPY --from=builder /app/main .
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/static ./static
COPY --from=builder /app/assets ./assets

# Создаем непривилегированного пользователя
RUN adduser -D -g '' appuser
USER appuser

# Переменные окружения будут передаваться при запуске контейнера
ENV APP_ENV=production

EXPOSE 8080

CMD ["./main"] 