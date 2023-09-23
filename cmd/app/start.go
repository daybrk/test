package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"test-task/internal/adapters/db"
	"test-task/internal/adapters/db/postgresdb"
	"test-task/internal/adapters/web"
	http2 "test-task/internal/controller/http"
	"test-task/internal/controller/kafka"
	"test-task/internal/domain/user"
	"time"
)

func NewEnrichment() (kafka.Kafka, http2.Handler) {
	l := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true})).WithGroup("enrichment_domain")

	storage := postgresdb.NewUserStorage(db.Connection, l)
	router := web.NewRouter(l)

	enrichmentService := user.NewUserService(storage, router, l)
	enrichmentUseCase := user.NewUserUseCase(enrichmentService, l)

	kafkaHandler := kafka.NewUserKafka(enrichmentUseCase, l)
	httpHandler := http2.NewUserHandler(enrichmentUseCase, l)

	return kafkaHandler, httpHandler
}

func Run(addr string) <-chan struct{} {
	mux := http.NewServeMux()

	kafkaEntry, handlerEntry := NewEnrichment()
	go kafkaEntry.Start()
	handlerEntry.Register(mux)

	server := &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  3 * time.Minute,
		WriteTimeout: 3 * time.Minute,
		IdleTimeout:  3 * time.Minute,
	}
	done := make(chan struct{})
	go func() {
		if err := server.ListenAndServe(); err != nil {
			fmt.Println(err)
		}
		close(done)
	}()
	return done
}
