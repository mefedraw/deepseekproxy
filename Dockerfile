# Используем официальный образ Go
FROM golang:1.21-alpine

# Создаем рабочую директорию
WORKDIR /app

# Копируем файлы go.mod и go.sum (если есть)
COPY go.mod ./
# COPY go.sum ./

# Скачиваем зависимости
RUN go mod download

# Копируем исходный код
COPY *.go ./

# Собираем приложение
RUN go build -o /proxy-dipsy

# Указываем порт, который будет использоваться
EXPOSE 3232

# Запускаем приложение
CMD ["/proxy-dipsy"]