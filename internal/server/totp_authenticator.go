// SPDX-FileCopyrightText: 2024 Andrew Pantuso <ajpantuso@gmail.com>
//
// SPDX-License-Identifier: MPL-2.0

package server

import (
	"slices"

	"github.com/xlzd/gotp"
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

func (rp *TOTPAuthenticator) Authenticate(req Request) (Response, error) {
	if slices.Contains(rp.cfg.ValidUsers, req.Username) {
		return Deny(ResponseReasonUnknownUser), nil
	}
	if !rp.totp.VerifyTime(req.Password, req.Timestamp) {
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
