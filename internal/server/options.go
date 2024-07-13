// SPDX-FileCopyrightText: 2024 Andrew Pantuso <ajpantuso@gmail.com>
//
// SPDX-License-Identifier: MPL-2.0

package server

import (
	"net"

	"github.com/go-logr/logr"
	"layeh.com/radius"
)

type WithBindAddress string

func (w WithBindAddress) ConfigureServer(c *ServerConfig) {
	c.BindAddr = string(w)
}

type WithHealthAddress string

func (w WithHealthAddress) ConfigureServer(c *ServerConfig) {
	c.HealthAddr = string(w)
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

func (w WithLogger) ConfigureHandler(c *HandlerConfig) {
	c.Logger = w.Logger
}

type WithAllowedClientSources struct {
	AllowedClientSources *net.IPNet
}

func (w WithAllowedClientSources) ConfigureHandler(c *HandlerConfig) {
	c.AllowedClientSources = w.AllowedClientSources
}

type WithAuthenticator struct {
	Authenticator Authenticator
}

func (w WithAuthenticator) ConfigureHandler(c *HandlerConfig) {
	c.Authenticator = w.Authenticator
}

type WithValidUsers []string

func (w WithValidUsers) ConfigureTOTPAuthenticator(c *TOTPAuthenticatorConfig) {
	c.ValidUsers = append(c.ValidUsers, w...)
}
