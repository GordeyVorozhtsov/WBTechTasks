# Dockerfile в корне проекта CLEAN_WEBSERVICE

FROM golang:1.23.5

WORKDIR /app

# Копируем модули и скачиваем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь исходный код
COPY . .

# Собираем бинарник
RUN go build -o main ./main.go

# Открываем порт 8081 (HTTP сервер)
EXPOSE 8081

# Запускаем приложение
CMD ["./main"]

