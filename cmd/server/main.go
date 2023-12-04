package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/okutsen/PasswordManager/config"
	"github.com/okutsen/PasswordManager/internal/api"
	"github.com/okutsen/PasswordManager/internal/controller"
	"github.com/okutsen/PasswordManager/internal/log"
	"github.com/okutsen/PasswordManager/internal/repo"
)

func main() {
	logger := log.New()
	config, err := config.New()
	if err != nil {
		logger.Fatalf("failed to initialize config: %s", err.Error())
	}

	db, err := repo.OpenConnection(&repo.Config{
		Host:     config.DB.Host,
		Port:     config.DB.Port,
		DBName:   config.DB.DBName,
		Username: config.DB.Username,
		SSLMode:  config.DB.SSLMode,
		Password: config.DB.Password,
	})
	if err != nil {
		logger.Fatalf("failed to open DB connection: %s", err.Error())
	}
	logger.Info("DB connected")

	userRepo, err := repo.NewUserRepository(db)
	if err != nil {
		logger.Fatalf("failed to init userRepo: %s", err.Error())
	}

	credentialRecordRepo, err := repo.NewCredentialRecordRepository(db)
	if err != nil {
		logger.Fatalf("failed to init credentialRecordRepo: %s", err.Error())
	}

	ctrl, err := controller.New(logger, userRepo, credentialRecordRepo)
	if err != nil {
		logger.Fatalf("failed to init ctrl: %s", err.Error())
	}

	serviceAPI := api.New(&api.Config{Port: config.API.Port}, ctrl, logger)

	go func() {
		logger.Infof("Staring server, listening on %s", config.API.Port)

		err = serviceAPI.Start()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Errorf("failed to start server %s", err.Error())
			return
		}
	}()

	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, syscall.SIGINT, syscall.SIGTERM)

	osCall := <-osSignals
	logger.Infof("system call: %v", osCall)

	err = repo.CloseConnection(db)
	if err != nil {
		logger.Warnf("failed to close DB: %s", err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.API.ShutdownTimeout)
	defer cancel()

	err = serviceAPI.Stop(ctx)
	if err != nil {
		logger.Fatalf("failed to stop application %s", err.Error())
	}

}
