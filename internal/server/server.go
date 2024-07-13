// SPDX-FileCopyrightText: 2024 Andrew Pantuso <ajpantuso@gmail.com>
//
// SPDX-License-Identifier: MPL-2.0

package server

import (
	"context"
	"time"

	"github.com/go-logr/logr"
	"layeh.com/radius"
)

func NewServer(opts ...ServerOption) *Server {
	var cfg ServerConfig

	cfg.Options(opts...)
	cfg.Default()

	return &Server{
		log: cfg.Logger,
		srv: &radius.PacketServer{
			Addr:         cfg.BindAddr,
			Handler:      cfg.Handler,
			SecretSource: radius.StaticSecretSource([]byte(cfg.Secret)),
		},
	}
}

type Server struct {
	srv *radius.PacketServer
	log logr.Logger
}

func (s *Server) Run(ctx context.Context) error {
	errCh := make(chan error)

	go func() {
		errCh <- s.srv.ListenAndServe()
	}()

	for {
		select {
		case <-ctx.Done():
			s.log.Info("shutting down")

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			return s.srv.Shutdown(ctx)
		case err := <-errCh:
			return err
		}
	}
}

type ServerConfig struct {
	BindAddr string
	Secret   string
	Handler  radius.Handler
	Logger   logr.Logger
}

func (c *ServerConfig) Options(opts ...ServerOption) {
	for _, opt := range opts {
		opt.ConfigureServer(c)
	}
}

func (c *ServerConfig) Default() {
	if c.BindAddr == "" {
		c.BindAddr = ":1812"
	}

	if c.Handler == nil {
		c.Handler = NewHandler(WithLogger{
			Logger: c.Logger.WithName("handler"),
		})
	}

	if c.Logger.GetSink() == nil {
		c.Logger = logr.Discard()
	}
}

type ServerOption interface {
	ConfigureServer(*ServerConfig)
}
