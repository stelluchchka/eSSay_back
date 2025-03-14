# Используем официальный образ Golang
FROM golang:1.23-alpine AS build

RUN apk update && apk add bash ca-certificates git gcc g++ libc-dev librdkafka-dev pkgconf

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем файлы проекта
COPY . .

# Загружаем зависимости
RUN go mod tidy

RUN go build -tags dynamic -o main src/cmd/app/main.go

# Запускаем приложение
CMD ["./main"]
