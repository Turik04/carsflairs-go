FROM golang:1.24

WORKDIR /app

# Копируем go.mod и go.sum из папки backend
COPY backend/go.mod ./ 
COPY backend/go.sum ./ 

# Загружаем зависимости
RUN go mod download

# Копируем весь исходный код из backend
COPY backend/ ./backend/

# Переходим в папку с кодом и собираем Go-приложение
WORKDIR /app/backend/cmd
RUN go build -o main .

# Команда по умолчанию
CMD ["./main"]
