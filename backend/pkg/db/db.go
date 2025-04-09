package db

import (
    "database/sql"
    "fmt"
    "log"
    "os"

    _ "github.com/lib/pq"
    "github.com/joho/godotenv"
)

var DB *sql.DB // Глобальная переменная для базы данных

func ConnectDB() {
    // Загружаем переменные окружения
    err := godotenv.Load("../../.env")
    if err != nil {
        log.Fatal("Ошибка загрузки .env файла")
    }

    // Создаём строку подключения
    connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
        os.Getenv("DB_HOST"),
        os.Getenv("DB_PORT"),
        os.Getenv("DB_USER"),
        os.Getenv("DB_PASSWORD"),
        os.Getenv("DB_NAME"),
        os.Getenv("SSL_MODE"),
    )

    var errDB error
    DB, errDB = sql.Open("postgres", connStr) // Используем глобальную DB
    if errDB != nil {
        log.Fatal("Ошибка подключения к базе:", errDB)
    }

    // Проверяем подключение
    errDB = DB.Ping()
    if errDB != nil {
        log.Fatal("База данных недоступна:", errDB)
    }

    fmt.Println("Подключение к базе успешно!")
}

func CloseDB() {
    if DB != nil {
        DB.Close()
        fmt.Println("Подключение к базе закрыто")
    }
}
