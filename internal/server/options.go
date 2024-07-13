// SPDX-FileCopyrightText: 2024 Andrew Pantuso <ajpantuso@gmail.com>
//
// SPDX-License-Identifier: MPL-2.0

package server

import (
	"github.com/go-logr/logr"
	"layeh.com/radius"
)

type WithBindAddress string

func (w WithBindAddress) ConfigureServer(c *ServerConfig) {
	c.BindAddr = string(w)
}

type WithSecret string

func (w WithSecret) ConfigureServer(c *ServerConfig) {
	c.Secret = string(w)
}

type WithHandler struct {
	Handler radius.Handler
}

func (w WithHandler) ConfigureServer(c *ServerConfig) {
	c.Handler = w.Handler
}

type WithLogger struct {
	Logger logr.Logger
}

func (w WithLogger) ConfigureServer(c *ServerConfig) {
	c.Logger = w.Logger
}

func (w WithLogger) ConfigureDefaultHandler(c *defaultHandlerConfig) {
	c.Logger = w.Logger
}
