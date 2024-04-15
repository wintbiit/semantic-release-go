package main

import (
	"flag"
	"os"
	"time"

	"github.com/wintbiit/semantic-release-go/utils"

	"github.com/wintbiit/semantic-release-go/types"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/wintbiit/semantic-release-go/semantic"
)

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	logF, err := os.OpenFile(".semantic-release.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to open log file")
	}

	log.Logger = log.Output(zerolog.MultiLevelWriter(
		zerolog.ConsoleWriter{Out: os.Stderr},
		zerolog.SyncWriter(logF),
	))
}

var version string = "v0.0.1"

func main() {
	start := time.Now()
	defer func() {
		log.Info().Str("elapsed", time.Since(start).String()).Msg("Semantic release done")
	}()

	log.Info().Str("version", version).Time("started", start).Msg("Starting semantic-release-go")

	var opt types.SemanticOptions
	var debug bool

	flag.StringVar(&opt.Path, "path", ".", "Path to git repository")
	flag.StringVar(&opt.Channel, "channel", os.Getenv("CHANNEL"), "Channel to release")
	flag.StringVar(&opt.Branch, "branch", os.Getenv("branch"), "Branch to release")
	flag.StringVar(&opt.Analyzer, "analyzer", os.Getenv("ANALYZER"), "Analyzer to use")
	flag.StringVar(&opt.Repo, "repo", os.Getenv("REPO"), "Repository to release")
	flag.BoolVar(&opt.Dry, "dry", false, "Dry run")
	flag.BoolVar(&opt.Tag, "tag", true, "Tag changes")
	flag.BoolVar(&opt.Push, "push", true, "Push changes")
	flag.StringVar(&opt.Changelog, "changelog", "Changelog.md", "Generate changelog")
	flag.BoolVar(&debug, "debug", false, "Debug mode")
	flag.Parse()

	if debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	if opt.Channel == "" {
		opt.Channel = types.ChannelInsider
		log.Warn().Msg("channel not set, using insider channel")
	}

	if !utils.ChannelValid(opt.Channel) {
		log.Fatal().Str("channel", opt.Channel).Msg("Invalid channel")
	}

	if opt.Branch == "" {
		log.Fatal().Msg("branch not set")
	}

	if opt.Analyzer == "" {
		opt.Analyzer = "angular"
		log.Warn().Msg("analyzer not set, using angular analyzer")
	}

	if opt.Repo == "" {
		log.Fatal().Msg("repo url not set")
	}

	semantic.Run(opt)
}
