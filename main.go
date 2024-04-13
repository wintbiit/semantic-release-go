package main

import (
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/wintbiit/semantic-release-go/semantic"

	"github.com/wintbiit/semantic-release-go/types"
)

var (
	SEASON   = os.Getenv("SEASON")
	CHANNEL  = os.Getenv("CHANNEL")
	ANALYZER = os.Getenv("ANALYZER")
	REPO     = os.Getenv("REPO")
)

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	log.Info().Str("version", version).Msg("Starting semantic-release-go")

	if SEASON == "" {
		log.Fatal().Msg("SEASON is required")
	}
	SEASON = strings.ToLower(SEASON)

	if CHANNEL == "" {
		CHANNEL = types.ChannelInsider
		log.Warn().Msg("CHANNEL is not set, use insider by default")
	}
	CHANNEL = strings.ToLower(CHANNEL)

	if ANALYZER == "" {
		ANALYZER = "angular"
	}
	ANALYZER = strings.ToLower(ANALYZER)

	if REPO == "" {
		log.Fatal().Msg("REPO is required")
	}

	logF, err := os.OpenFile(".semantic-release.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to open log file")
	}

	log.Logger = log.Output(zerolog.MultiLevelWriter(
		zerolog.ConsoleWriter{Out: os.Stderr},
		zerolog.SyncWriter(logF),
	))

	if os.Getenv("DEBUG") == "true" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
}

var version string = "v0.0.1"

func main() {
	start := time.Now()
	semantic.Run(".", CHANNEL, SEASON, ANALYZER, REPO)
	log.Info().Str("duration", time.Since(start).String()).Msg("Semantic release done")
}
