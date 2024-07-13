// SPDX-FileCopyrightText: 2024 Andrew Pantuso <ajpantuso@gmail.com>
//
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/ajpantuso/dyn-radius/internal/config"
	"github.com/ajpantuso/dyn-radius/internal/server"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer stop()

	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("loading config: %v\n", err)
		os.Exit(1)
	}

	zlog, err := zap.NewDevelopment()
	if err != nil {
		fmt.Printf("initializing logger: %v\n", err)
		os.Exit(1)
	}

	logger := zapr.NewLogger(zlog)
	serverLogger := logger.WithName("server")

	srv := server.NewServer(
		server.WithBindAddress(cfg.BindAddr),
		server.WithHealthAddress(cfg.HealthAddr),
		server.WithLogger{
			Logger: serverLogger,
		},
		server.WithSecret(cfg.ClientSecret),
		server.WithHandler{
			Handler: server.NewHandler(
				server.WithAllowedClientSources{
					AllowedClientSources: cfg.AllowedClientSources,
				},
				server.WithAuthenticator{
					Authenticator: server.NewTOTPAuthenticator(cfg.TOTPSecret),
				},
				server.WithLogger{
					Logger: serverLogger.WithName("handler"),
				},
			),
		},
	)

	logger.Info("starting server", "bindAddr", cfg.BindAddr, "healthAddr", cfg.HealthAddr)
	if err := srv.Run(ctx); err != nil {
		logger.Error(err, "server exited unexpectedly")

		os.Exit(1)
	}
}
