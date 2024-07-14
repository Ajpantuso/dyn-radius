// SPDX-FileCopyrightText: 2024 Andrew Pantuso <ajpantuso@gmail.com>
//
// SPDX-License-Identifier: MPL-2.0

package server

import (
	"slices"

	"github.com/xlzd/gotp"
	"layeh.com/radius/rfc2865"
)

func NewTOTPAuthenticator(secret string) *TOTPAuthenticator {
	return &TOTPAuthenticator{
		totp: *gotp.NewDefaultTOTP(secret),
	}
}

type TOTPAuthenticator struct {
	cfg  TOTPAuthenticatorConfig
	totp gotp.TOTP
}

func (a *TOTPAuthenticator) Authenticate(req Request) (Response, error) {
	var (
		username = rfc2865.UserName_GetString(req.Req.Packet)
		password = rfc2865.UserPassword_GetString(req.Req.Packet)
	)

	if slices.Contains(a.cfg.ValidUsers, username) {
		return Deny(ResponseReasonUnknownUser), nil
	}
	if !a.totp.VerifyTime(password, req.Timestamp) {
		return Deny(ResponseReasonInvalidPassword), nil
	}

	return Approve(), nil
}

type TOTPAuthenticatorConfig struct {
	ValidUsers []string
}

func (c *TOTPAuthenticatorConfig) Options(opts ...TOTPAuthenticatorOption) {
	for _, opt := range opts {
		opt.ConfigureTOTPAuthenticator(c)
	}
}

type TOTPAuthenticatorOption interface {
	ConfigureTOTPAuthenticator(*TOTPAuthenticatorConfig)
}
