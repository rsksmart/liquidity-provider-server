package main

import (
	"testing"

	"github.com/rsksmart/liquidity-provider-server/internal/configuration/environment"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
	"github.com/stretchr/testify/require"
)

func TestBuildLogConfig(t *testing.T) {
	t.Run("uses release tag when set", func(t *testing.T) {
		liquidity_provider.BuildVersion = "v2.6.0"
		BuildVersion = "abc123"
		t.Cleanup(func() {
			liquidity_provider.BuildVersion = ""
			BuildVersion = ""
		})

		env := environment.Environment{
			LpsStage:  "testnet",
			LogFormat: "json",
			LogLevel:  "info",
			LogFile:   "/tmp/lps.log",
		}

		cfg := buildLogConfig(env)

		require.Equal(t, logServiceName, cfg.Service)
		require.Equal(t, "testnet", cfg.Environment)
		require.Equal(t, "v2.6.0", cfg.Version)
		require.Equal(t, "json", cfg.Format)
		require.Equal(t, "info", cfg.Level)
		require.Equal(t, "/tmp/lps.log", cfg.File)
	})

	t.Run("falls back to commit when release tag empty", func(t *testing.T) {
		liquidity_provider.BuildVersion = ""
		BuildVersion = "abc123"
		t.Cleanup(func() {
			liquidity_provider.BuildVersion = ""
			BuildVersion = ""
		})

		env := environment.Environment{
			LpsStage:  "regtest",
			LogFormat: "logfmt",
			LogLevel:  "debug",
		}

		cfg := buildLogConfig(env)

		require.Equal(t, "abc123", cfg.Version)
		require.Equal(t, "logfmt", cfg.Format)
	})
}

func TestApplyLogConfig(t *testing.T) {
	liquidity_provider.BuildVersion = "v2.6.0"
	t.Cleanup(func() {
		liquidity_provider.BuildVersion = ""
		logConfig = LogConfig{}
	})

	env := environment.Environment{
		LpsStage:  "mainnet",
		LogFormat: "json",
		LogLevel:  "warn",
		LogFile:   "/var/log/lps.log",
	}

	cfg := applyLogConfig(env)

	require.Equal(t, cfg, logConfig)
	require.Equal(t, logServiceName, logConfig.Service)
	require.Equal(t, "mainnet", logConfig.Environment)
}
