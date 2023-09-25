package main

import (
	"fmt"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"log/slog"
	"net/http"
	"os"
	"test-task/internal/adapters/db"
	"test-task/internal/adapters/db/postgresdb"
	"test-task/internal/adapters/web"
	"test-task/internal/controller/graphQL"
	http2 "test-task/internal/controller/http"
	"test-task/internal/controller/kafka"
	"test-task/internal/domain/user"
	"time"
)

func NewUser(mux *http.ServeMux) {
	l := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))
	postgres, err := db.ConnectToPostgres(l)
	if err != nil {
		l.Error("Ошибка при подключении к базе данных", slog.String("err", err.Error()))

		return
	}

	storage := postgresdb.NewUserStorage(postgres, l)
	router := web.NewRouter(l)

	userService := user.NewUserService(storage, router, l)
	userUseCase := user.NewUserUseCase(userService, l)

	srv := handler.NewDefaultServer(
		graphQL.NewExecutableSchema(graphQL.Config{Resolvers: graphQL.NewUserResolver(userUseCase, l)}))
	kafkaHandler := kafka.NewUserKafka(userUseCase, l)
	httpHandler := http2.NewUserHandler(userUseCase, l)

	go kafkaHandler.Start()
	httpHandler.Register(mux)
	mux.Handle("/", playground.Handler("GraphQL playground", "/query"))
	mux.Handle("/query", srv)

	l.Info("connect to http://localhost:8082/ for GraphQL playground")
}

func Run(addr string) <-chan struct{} {
	mux := http.NewServeMux()

	NewUser(mux)

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
