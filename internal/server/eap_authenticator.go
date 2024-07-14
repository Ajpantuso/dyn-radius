// SPDX-FileCopyrightText: 2024 Andrew Pantuso <ajpantuso@gmail.com>
//
// SPDX-License-Identifier: MPL-2.0

package server

import (
	"errors"
	"fmt"

	"github.com/ajpantuso/dyn-radius/internal/eap"
)

func NewEAPAuthenticator(secret string) *EAPAuthenticator {
	return &EAPAuthenticator{}
}

type EAPAuthenticator struct {}

// TODO: Implement PEAP detection and protocol
// TODO: Implement MSCHAPV2 detection and protocol
// TODO: Connect that to an identity store
func (a *EAPAuthenticator) Authenticate(req Request) (Response, error) {
	eapMsg, found := req.Req.Attributes.Lookup(79)
	if !found {
		return Deny(), ErrMissingEAPMessage
	}

	var packet eap.Packet

	if err := packet.UnmarshalBinary(eapMsg); err != nil {
		return Deny(), fmt.Errorf("unmarshalling packet: %w", err)
	}

	// After unmarshalling packet the data must be decoded
	// as a PEAP message (since that is all we will implement).
	// Within the PEAP message an MSCHAPV2 message should be
	// present (again this is all we will implement).

	return Approve(), nil
}

var ErrMissingEAPMessage = errors.New("missing EAP message")
