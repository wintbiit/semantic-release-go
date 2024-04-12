package semantic

import (
	"log/slog"
	"os"
)

var (
	SEASON  = os.Getenv("SEASON")
	CHANNEL = os.Getenv("CHANNEL")
)

const (
	ChannelInsider = "insider"
	ChannelAlpha   = "alpha"
	ChannelBeta    = "beta"
	ChannelRelease = "release"
)

func init() {
	if SEASON == "" {
		slog.Error("SEASON is required")
		os.Exit(1)
	}
	if CHANNEL == "" {
		CHANNEL = ChannelInsider
		slog.Warn("CHANNEL is not set, use insider by default")
	}
}
