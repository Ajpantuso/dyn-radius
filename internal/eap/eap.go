// SPDX-FileCopyrightText: 2024 Andrew Pantuso <ajpantuso@gmail.com>
//
// SPDX-License-Identifier: MPL-2.0

// source: https://github.com/theaaf/radius-server

package eap

import (
	"encoding"
	"encoding/binary"
	"errors"
	"fmt"
)

// https://datatracker.ietf.org/doc/html/rfc3748#section-4
type Packet struct {
	Code       Code
	Identifier byte
	Length     uint16
	Data       []byte
}

func (p *Packet) UnmarshalBinary(data []byte) error {
	buf := make([]byte, len(data))
	copy(buf, data)

	if len(buf) < 4 {
		return fmt.Errorf("%w: packet must be 4 or more bytes", ErrMalformedPacket)
	}

	code, err := CodeFromByte(buf[0])
	if err != nil {
		return fmt.Errorf("processing code byte: %w", err)
	}

	p.Code = code
	p.Identifier = data[1]
	p.Length = binary.BigEndian.Uint16(data[2:4])
	p.Data = data[4:p.Length]

	return nil
}

var ErrMalformedPacket = errors.New("malformed packet")

func (p *Packet) MarshalBinary() ([]byte, error) {
	buf := make([]byte, p.Length)

	buf[0] = p.Code.ToByte()
	buf[1] = p.Identifier
	binary.BigEndian.PutUint16(buf[2:4], p.Length)
	copy(buf[4:], p.Data)

	return buf, nil
}

func (p *Packet) DecodeData() (interface{}, error) {
	if len(p.Data) == 0 {
		return nil, nil
	}

	var v encoding.BinaryUnmarshaler

	t, err := RequestResponseTypeFromByte(p.Data[0])
	if err != nil {
		return nil, fmt.Errorf("parsing request-response type: %w", err)
	}

	switch t {
	case RequestResponseTypeIdentity:
		v = &Identity{}
	}
	if err := v.UnmarshalBinary(p.Data); err != nil {
		return nil, fmt.Errorf("unmarshalling packet data: %w", err)
	}

	return v, nil
}

func CodeFromByte(b byte) (Code, error) {
	switch Code(b) {
	case CodeRequest:
		return CodeRequest, nil
	case CodeResponse:
		return CodeResponse, nil
	case CodeSuccess:
		return CodeSuccess, nil
	case CodeFailure:
		return CodeFailure, nil
	default:
		return Code(0), fmt.Errorf("unknown code")
	}
}

type Code byte

func (c Code) ToByte() byte {
	return byte(c)
}

const (
	CodeRequest  Code = 1
	CodeResponse Code = 2
	CodeSuccess  Code = 3
	CodeFailure  Code = 4
)

func RequestResponseTypeFromByte(b byte) (RequestResponseType, error) {
	switch RequestResponseType(b) {
	case RequestResponseTypeIdentity:
		return RequestResponseTypeIdentity, nil
	case RequestResponseTypeNotification:
		return RequestResponseTypeNotification, nil
	default:
		return RequestResponseType(0), UnsupportedRequestResponseTypeError{}
	}
}

type RequestResponseType byte

func (t RequestResponseType) ToByte() byte {
	return byte(t)
}

const (
	RequestResponseTypeIdentity     RequestResponseType = 1
	RequestResponseTypeNotification RequestResponseType = 2
)

type UnsupportedRequestResponseTypeError struct {
	Type int
}

func (e UnsupportedRequestResponseTypeError) Error() string {
	return fmt.Sprintf("unsupported request-response type: %d", e.Type)
}

// https://datatracker.ietf.org/doc/html/rfc3748#section-5.1
type Identity struct {
	Identity string
}

func (i *Identity) UnmarshalBinary(data []byte) error {
	if rrt, _ := RequestResponseTypeFromByte(data[0]); rrt != RequestResponseTypeIdentity {
		return ErrRequestResponseTypeMismatch
	}
	i.Identity = string(data[1:])

	return nil
}

var ErrRequestResponseTypeMismatch = errors.New("request-response type mismatch")
