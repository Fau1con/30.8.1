package main

import (
	"log"
	"net/http"

	"30.8.1/internal/config"
	"30.8.1/internal/core"
	"30.8.1/internal/repository"
	"30.8.1/internal/server"
)

func main() {
	// Загрузка конфигурации
	cfg, err := config.LoadConfig("dbconfig.yaml")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Инициализация хранилища
	repo, err := repository.NewPostgresRepository(cfg.DBConnectionString())
	if err != nil {
		log.Fatalf("Failed to initialize repository: %v", err)
	}
	defer repo.Close()

	// Инициализация сервиса задач
	taskService := core.NewTaskService(repo)

	// Инициализация HTTP-сервера
	httpServer := server.NewHTTPServer(taskService)

	// Запуск сервера
	log.Println("Starting server on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", httpServer.Router()))
}
