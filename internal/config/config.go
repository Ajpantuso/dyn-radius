// SPDX-FileCopyrightText: 2024 Andrew Pantuso <ajpantuso@gmail.com>
//
// SPDX-License-Identifier: MPL-2.0

package config

import (
	"errors"
	"os"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func Load() (*Config, error) {
	if err := processFlags(); err != nil {
		return nil, err
	}
	if err := processEnvVars(); err != nil {
		return nil, err
	}

	secretPath := viper.GetString("client-secret-path")
	if secretPath == "" {
		return nil, errors.New(`"--client-secret-path" is missing`)
	}

	secret, err := getClientSecret(secretPath)
	if err != nil {
		return nil, err
	}

	return &Config{
		BindAddr:     viper.GetString("bind-addr"),
		ClientSecret: secret,
	}, nil
}

func processFlags() error {
	pflag.String("bind-addr", ":1812", "bind address for the RADIUS server")
	pflag.String("client-secret-path", "", "path to file containing client secret")
	pflag.Parse()

	return viper.BindPFlags(pflag.CommandLine)
}

func processEnvVars() error {
	for _, vs := range [][]string{
		{"bind-addr", "DYN_RADIUS_BIND_ADDR"},
		{"client-secret-path", "DYN_RADIUS_CLIENT_SECRET_PATH"},
	} {
		if err := viper.BindEnv(vs...); err != nil {
			return err
		}
	}

	return nil
}

func getClientSecret(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(data)), nil
}

type Config struct {
	BindAddr     string
	ClientSecret string
}
