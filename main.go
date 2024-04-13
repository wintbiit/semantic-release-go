package main

import (
	"log/slog"
	"os"

	"github.com/wintbiit/semantic-release-go/semantic"
)

var (
	SEASON   = os.Getenv("SEASON")
	CHANNEL  = os.Getenv("CHANNEL")
	ANALYZER = os.Getenv("ANALYZER")
)

func init() {
	if SEASON == "" {
		slog.Error("SEASON is required")
		os.Exit(1)
	}

	if CHANNEL == "" {
		CHANNEL = semantic.ChannelInsider
		slog.Warn("CHANNEL is not set, use insider by default")
	}

	if ANALYZER == "" {
		ANALYZER = "angular"
	}
}

func main() {
	semantic.Run(".", "insider", SEASON, ANALYZER)
}
