package main

import (
	"banking-api/internal/config"
	"banking-api/internal/handler"
	"banking-api/internal/middleware"
	"banking-api/internal/repository"
	"banking-api/internal/service"
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

func main() {
	// Загрузка .env
	if err := godotenv.Load(); err != nil {
		log.Fatal("Ошибка загрузки .env файла")
	}

	// Инициализация логгера
	config.InitLogger()
	config.Log.Info("Логгер инициализирован")

	// Конфигурация
	cfg := config.LoadConfig()

	// Подключение к БД
	db, err := sql.Open("postgres", cfg.DBUrl)
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("БД недоступна: %v", err)
	}
	config.Log.Info("Подключение к базе данных успешно")

	// Репозитории
	userRepo := repository.NewUserRepository(db)
	accountRepo := repository.NewAccountRepository(db)

	// Email-сервис
	emailService := service.NewEmailService(
		cfg.SMTPHost,
		cfg.SMTPPort,
		cfg.SMTPUser,
		cfg.SMTPPass,
	)

	// Сервисы
	authService := service.NewAuthService(userRepo)
	accountService := service.NewAccountService(accountRepo, userRepo, emailService)

	// Хендлеры
	authHandler := handler.NewAuthHandler(authService)
	accountHandler := handler.NewAccountHandler(accountService)

	// Роутинг
	router := mux.NewRouter()

	// Публичные
	router.HandleFunc("/register", authHandler.Register).Methods("POST")
	router.HandleFunc("/login", authHandler.Login).Methods("POST")

	// Защищённые
	protected := router.PathPrefix("/").Subrouter()
	protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))

	protected.HandleFunc("/accounts", accountHandler.Create).Methods("POST")
	protected.HandleFunc("/accounts/topup", accountHandler.TopUp).Methods("POST")
	protected.HandleFunc("/transfer", accountHandler.Transfer).Methods("POST")
	protected.HandleFunc("/transfer/by-usernames", accountHandler.TransferByUsernames).Methods("POST")

	// Запуск сервера
	config.Log.Infof("Сервер запущен на порту %s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, router))
}
