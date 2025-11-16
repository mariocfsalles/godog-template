package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"orders-api/api"   // ajuste o módulo
	"orders-api/events"
	"orders-api/store"
)

// helper simples de env
func getenv(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}

func main() {
	port := getenv("PORT", "3000")
	dsn := getenv("DB_DSN", "app:apppass@tcp(mysql:3306)/orders?parseTime=true&charset=utf8mb4&collation=utf8mb4_0900_ai_ci")

	// DB (reset=true porque é projeto de testes)
	db := store.MustMySQL(dsn, true)
	defer db.Close()

	// Kafka
	brokers := strings.Split(getenv("KAFKA_BROKERS", "kafka:9092"), ",")
	topic := getenv("KAFKA_TOPIC", "orders.events")
	clientID := getenv("KAFKA_CLIENT_ID", "orders-api")
	publisher := events.NewPublisher(brokers, topic, clientID)
	defer publisher.Close()

	// API HTTP
	apiServer := api.NewServer(db, publisher)

	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           apiServer,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Printf("API → http://0.0.0.0:%s", port)
		if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("http: %v", err)
		}
	}()

	// shutdown gracioso
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
}
