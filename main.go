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
		fmt.Println("loading config: %w", err)
		os.Exit(1)
	}

	zlog, err := zap.NewDevelopment()
	if err != nil {
		fmt.Println("initializing logger: %w", err)
		os.Exit(1)
	}

	logger := zapr.NewLogger(zlog)

	srv := server.NewServer(
		server.WithBindAddress(cfg.BindAddr),
		server.WithLogger{
			Logger: logger.WithName("server"),
		},
		server.WithSecret(cfg.ClientSecret),
	)

	logger.Info("starting server", "bindAddr", cfg.BindAddr)
	if err := srv.Run(ctx); err != nil {
		logger.Error(err, "server exited unexpectedly")

		os.Exit(1)
	}
}