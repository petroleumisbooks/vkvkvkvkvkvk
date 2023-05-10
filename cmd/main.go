package main

import (
	"github.com/joho/godotenv"
	h "internbot/internal/handler"
	"internbot/internal/usecase"
	_ "internbot/internal/usecase"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("файл .env не найден")
	}
}

func main() {
	token, exists := os.LookupEnv("TOKEN")
	if !exists {
		log.Print("в файле .env не указаны переменные окружения: TOKEN")
	}
	handler := h.Handler{}

	usecase := usecase.UseCase{
		Token:   token,
		Handler: handler,
	}

	go func() {
		usecase.Runing = true
		usecase.Run()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Print("приложение останавливается")

	usecase.Shutdown()
}
