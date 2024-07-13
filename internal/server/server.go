// SPDX-FileCopyrightText: 2024 Andrew Pantuso <ajpantuso@gmail.com>
//
// SPDX-License-Identifier: MPL-2.0

package server

import (
	"context"
	"time"

	"layeh.com/radius"
	"layeh.com/radius/rfc2865"
	"github.com/go-logr/logr"
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

// TODO: Need to also serve REST API for approvals
func (s *Server) Run(ctx context.Context) error {
	errCh := make(chan error)

	go func() {
		errCh <- s.srv.ListenAndServe()
	}()

	for {
		select {
		case <-ctx.Done():
			s.log.Info("received shutdown signal")

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
		c.Handler = newDefaultHandler(WithLogger{
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

func newDefaultHandler(opts ...defaultHandlerOption) *defaultHandler {
	var cfg defaultHandlerConfig
	cfg.Options(opts...)

	return &defaultHandler{
		cfg: cfg,
	}
}

type defaultHandler struct{
	cfg defaultHandlerConfig
}

// TODO: Actually implement dynamic approval logic here
// TODO: Restrict client by source IP
func (h *defaultHandler) ServeRADIUS(w radius.ResponseWriter, r *radius.Request) {
	username := rfc2865.UserName_GetString(r.Packet)
	password := rfc2865.UserPassword_GetString(r.Packet)

	var code radius.Code
	if username == "test-user" && password == "password" {
		code = radius.CodeAccessAccept
	} else {
		code = radius.CodeAccessReject
	}
	h.cfg.Logger.Info("replying to request", "responseCode", code, "remoteAddress", r.RemoteAddr)
	_ = w.Write(r.Response(code))
}

type defaultHandlerConfig struct {
	Logger logr.Logger
}

func (c *defaultHandlerConfig) Options(opts ...defaultHandlerOption) {
	for _, opt := range opts {
		opt.ConfigureDefaultHandler(c)
	}
}

func (c *defaultHandlerConfig) Default() {
	if c.Logger.GetSink() == nil {
		c.Logger = logr.Discard()
	}
}

type defaultHandlerOption interface {
	ConfigureDefaultHandler(*defaultHandlerConfig)
}
