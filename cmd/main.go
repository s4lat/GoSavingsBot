package main

import (
	"context"
	"fmt"
	"github.com/s4lat/gosavingsbot/config"
	"github.com/s4lat/gosavingsbot/internal/common/postgres"
	v1 "github.com/s4lat/gosavingsbot/internal/controller/http/v1"
	postgres2 "github.com/s4lat/gosavingsbot/internal/infrastructure/postgres"
	"github.com/s4lat/gosavingsbot/internal/log"
	"github.com/s4lat/gosavingsbot/internal/usecase"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx := context.Background()
	cfg := config.GetConfig()
	defer func() { _ = log.Close() }()

	pg, err := postgres.NewPostgresClient(ctx, cfg.PgUrl, cfg.PgMaxConns, cfg.PgConnTimeout)
	if err != nil {
		log.Sugar().Fatalf("can't initialize postgres client: %s", err)
	}

	userRepo := postgres2.NewUserRepo(pg)
	expenseRepo := postgres2.NewExpenseRepo(pg)
	useCase := usecase.NewUseCase(userRepo, expenseRepo)

	router := v1.NewRouter(useCase, cfg.BotToken, false)

	errorChan := make(chan error)

	httpServer := serveHTTPInBackground(errorChan, router, fmt.Sprintf(":%s", cfg.HTTPPort))

	// For graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-quit:
		log.Sugar().Infof("got signal: %s", sig)
	case err := <-errorChan:
		log.Sugar().Errorf("received err from errorChan: %s", err)
	}

	log.Sugar().Info("shutdown httpServer...")

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Sugar().Fatalf("err on httpServer.Shutdown: %s", err)
	}

	<-ctx.Done()
	log.Sugar().Info("httpServer exited!")
}

func serveHTTPInBackground(errorChan chan<- error, handler http.Handler, addr string) *http.Server {
	srv := &http.Server{
		Addr:              addr,
		Handler:           handler,
		ReadHeaderTimeout: time.Second * 60,
	}

	go func() {
		log.Sugar().Info("serving http")
		errorChan <- fmt.Errorf("httpServer: %w", srv.ListenAndServe())
	}()

	return srv
}
