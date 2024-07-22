// Package main предоставляет основной функционал для запуска API сервера, который отправляет сообщения в Kafka.
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	handlers "github.com/vova4o/messagio/internal/api"
	"github.com/vova4o/messagio/internal/config"
	"github.com/vova4o/messagio/internal/service"
	"github.com/vova4o/messagio/internal/storage"
)

// ServerConfig содержит конфигурацию сервера, включая порт, обработчики и сервисы.
type ServerConfig struct {
	Port    string             // Порт, на котором запускается сервер
	Handler *gin.Engine        // Обработчик запросов Gin
	Handle  *handlers.Handlers // Обработчики API
	Service *service.Service   // Сервис для бизнес-логики
}

// @title Api to Send Messages to Kafka
// @version 1.0
// @description Api server to send messages to Kafka broker

// @host localhost:8080
// @BasePath /

// NewApp инициализирует и настраивает приложение, соединяя различные компоненты системы.
// Возвращает сконфигурированный экземпляр ServerConfig.
func NewApp() *ServerConfig {
	port := config.Address()
	kafka := config.Kafka()
	kafkaGroup := config.KafkaGroup()

	dbConn := storage.NewConn(kafka, kafkaGroup)

	gin.SetMode(gin.ReleaseMode)

	handler := gin.Default()

	serv := service.NewService(dbConn)
	// web := web.Templates()
	handl := handlers.NewHandler(serv, handler)

	handl.SetupRoutes()

	// Start checking consumed messages and getting there id and mark as produces in DB
	serv.StartConsumeMessages()

	return &ServerConfig{
		Port:    port,
		Handler: handler,
		Handle:  handl,
		Service: serv,
	}
}

// NewServer creates an instance of http.Server
func (c *ServerConfig) NewServer() *http.Server {
	return &http.Server{
		Addr:    c.Port,
		Handler: c.Handler,
	}
}

// StartServer does what it called as
func (c *ServerConfig) StartServer() {
	go func() {
		log.Printf("Starting server on %s\n", c.Port)
		if err := c.NewServer().ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()
}

// ShutdownServer this function provides setup for graceful shutdown when signal is recived
func (c *ServerConfig) ShutdownServer(srv *http.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	err := c.Service.CloseDB()
	if err != nil {
		log.Println(err)
	}

	log.Println("Server exiting")
}

func main() {
	app := NewApp()

	app.StartServer()

	app.ShutdownServer(app.NewServer())
}
