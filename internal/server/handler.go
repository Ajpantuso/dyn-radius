// SPDX-FileCopyrightText: 2024 Andrew Pantuso <ajpantuso@gmail.com>
//
// SPDX-License-Identifier: MPL-2.0

package server

import (
	"fmt"
	"net"
	"net/netip"
	"time"

	"github.com/go-logr/logr"
	"layeh.com/radius"
)

func NewHandler(opts ...HandlerOption) *Handler {
	var cfg HandlerConfig
	cfg.Options(opts...)

	return &Handler{
		cfg: cfg,
	}
}

type Handler struct {
	cfg HandlerConfig
}

func (h *Handler) ServeRADIUS(w radius.ResponseWriter, r *radius.Request) {
	debug := h.cfg.Logger.V(1)
	debug.Info("received request", "request", r)

	log := h.cfg.Logger.WithValues("remoteAddress", r.RemoteAddr)

	log.Info("processing request")

	validSource, err := h.isSourceValid(r)
	if err != nil {
		log.Error(err, "validating source IP")

		if err := w.Write(r.Response(radius.CodeAccessReject)); err != nil {
			log.Error(err, "writing response")
		}

		return
	}

	if !validSource {
		log.Info("rejecting request")
		debug.Info("invalid source IP", "allowedClientSources", h.cfg.AllowedClientSources, "source", r.RemoteAddr)

		if err := w.Write(r.Response(radius.CodeAccessReject)); err != nil {
			log.Error(err, "writing response")
		}

		return
	}

	res, err := h.cfg.Authenticator.Authenticate(Request{
		Req:       r,
		Timestamp: time.Now(),
	})
	if err != nil {
		log.Error(err, "authenticating request")

		if err := w.Write(r.Response(radius.CodeAccessReject)); err != nil {
			log.Error(err, "writing response")
		}

		return
	}

	if !res.approved {
		log.Info("rejecting request")
		debug.Info("rejecting request", "reasons", res.reasons)

		if err := w.Write(r.Response(radius.CodeAccessReject)); err != nil {
			log.Error(err, "writing response")
		}

		return
	}

	log.Info("accepting request")

	_ = w.Write(r.Response(radius.CodeAccessAccept))
}

func (h *Handler) isSourceValid(r *radius.Request) (bool, error) {
	if h.cfg.AllowedClientSources == nil {
		return true, nil
	}

	addrPort, err := netip.ParseAddrPort(r.RemoteAddr.String())
	if err != nil {
		return false, fmt.Errorf("parsing remote address: %w", err)
	}

	return h.cfg.AllowedClientSources.Contains(addrPort.Addr().AsSlice()), nil
}

type HandlerConfig struct {
	AllowedClientSources *net.IPNet
	Authenticator        Authenticator
	Logger               logr.Logger
}

func (c *HandlerConfig) Options(opts ...HandlerOption) {
	for _, opt := range opts {
		opt.ConfigureHandler(c)
	}
}

func (c *HandlerConfig) Default() {
	if c.Logger.GetSink() == nil {
		c.Logger = logr.Discard()
	}
}

type HandlerOption interface {
	ConfigureHandler(*HandlerConfig)
}

type Authenticator interface {
	Authenticate(Request) (Response, error)
}

type Request struct {
	Req       *radius.Request
	Timestamp time.Time
}

func Approve() Response {
	return Response{
		approved: true,
		reasons:  []ResponseReason{ResponseReasonValid},
	}
}

func Deny(reasons ...ResponseReason) Response {
	return Response{
		approved: false,
		reasons:  reasons,
	}
}

type Response struct {
	approved bool
	reasons  []ResponseReason
}

type ResponseReason string

const (
	ResponseReasonValid           ResponseReason = "valid request"
	ResponseReasonUnknownUser     ResponseReason = "unknown user"
	ResponseReasonInvalidPassword ResponseReason = "invalid password"
)
