// SPDX-FileCopyrightText: 2024 Andrew Pantuso <ajpantuso@gmail.com>
//
// SPDX-License-Identifier: MPL-2.0

package config

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func Load() (*Config, error) {
	if err := processFlags(); err != nil {
		return nil, fmt.Errorf("processing flags: %w", err)
	}
	if err := processEnvVars(); err != nil {
		return nil, fmt.Errorf("processing env vars: %w", err)
	}

	clientSecretPath := viper.GetString("client-secret-path")
	if clientSecretPath == "" {
		return nil, errors.New(`"--client-secret-path" is missing`)
	}

	clientSecret, err := getSecret(clientSecretPath)
	if err != nil {
		return nil, fmt.Errorf("getting client secret: %w", err)
	}

	totpSecretPath := viper.GetString("totp-secret-path")
	if totpSecretPath == "" {
		return nil, errors.New(`"--client-secret-path" is missing`)
	}

	totpSecret, err := getSecret(totpSecretPath)
	if err != nil {
		return nil, fmt.Errorf("getting TOTP secret: %w", err)
	}

	var allowedClientSources *net.IPNet
	if raw := viper.GetString("allowed-client-sources"); raw != "" {
		_, allowedClientSources, err = net.ParseCIDR(raw)
		if err != nil {
			return nil, fmt.Errorf("parsing allowed client sources: %w", err)
		}
	}

	return &Config{
		AllowedClientSources: allowedClientSources,
		BindAddr:             viper.GetString("bind-addr"),
		HealthAddr:           viper.GetString("health-addr"),
		ClientSecret:         clientSecret,
		TOTPSecret:           totpSecret,
	}, nil
}

func processFlags() error {
	pflag.String("allowed-client-sources", "", "network from which client requests are allowed")
	pflag.String("bind-addr", ":51812", "bind address for the RADIUS server")
	pflag.String("health-addr", ":8080", "bind address for the health endpoint")
	pflag.String("client-secret-path", "./config/client-secret", "path to file containing client secret")
	pflag.String("totp-secret-path", "./config/totp-secret", "path to file containing totp secret")
	pflag.StringSlice("valid-users", []string{}, "valid usernames for authentication")
	pflag.Parse()

	return viper.BindPFlags(pflag.CommandLine)
}

func processEnvVars() error {
	for _, vs := range [][]string{
		{"allowed-client-sources", "DYN_RADIUS_ALLOWED_CLIENT_SOURCES"},
		{"bind-addr", "DYN_RADIUS_BIND_ADDR"},
		{"health-addr", "DYN_RADIUS_HEALTH_ADDR"},
		{"client-secret-path", "DYN_RADIUS_CLIENT_SECRET_PATH"},
		{"valid-users", "DYN_RADIUS_VALID_USERS"},
	} {
		if err := viper.BindEnv(vs...); err != nil {
			return fmt.Errorf("binding env var: %w", err)
		}
	}

	return nil
}

func getSecret(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("reading secret file: %w", err)
	}

	return strings.TrimSpace(string(data)), nil
}

type Config struct {
	AllowedClientSources *net.IPNet
	BindAddr             string
	HealthAddr           string
	ClientSecret         string
	TOTPSecret           string
}
