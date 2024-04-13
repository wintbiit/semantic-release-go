package semantic

import (
	"github.com/go-git/go-git/v5"
	"github.com/rs/zerolog/log"
	"github.com/wintbiit/semantic-release-go/analyze"
	"github.com/wintbiit/semantic-release-go/types"

	"github.com/wintbiit/semantic-release-go/utils"
)

func Run(path, channel, season, analyzer string) {
	r, err := git.PlainOpen(path)
	if err != nil {
		log.Fatal().Err(err).Msg("Not a git repository")
	}

	result := &types.Result{
		Season:  season,
		Channel: channel,
	}

	tags, err := r.Tags()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to get tags")
	}

	scannedTags, err := utils.History(tags, season, channel)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to scan tags")
	}

	currentCommit, err := r.Head()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to get current commit")
	}

	log.Info().Msgf("Current commit: %v", utils.HashShort(currentCommit.Hash()))
	log.Info().Msgf("Using channel: %s, season: %s", channel, season)

	if len(scannedTags) == 0 {
		log.Info().Msgf("No history of %s %s, will use vcs tree tail and release first version v1.0.0", season, channel)
	} else {
		result.LatestRelease = scannedTags[len(scannedTags)-1]
		log.Info().Msgf("Last release: %s", result.LatestRelease.String())
	}

	// get commits since last version
	commits, err := r.Log(&git.LogOptions{From: currentCommit.Hash()})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to get commits")
	}

	sinceCommits, err := utils.CommitsSince(commits, result.LatestRelease.Hash().String())
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to get commits since last version")
	}

	result.Commits = sinceCommits
	log.Info().Msgf("Commits since last version: %d", len(sinceCommits))
	log.Info().Msg("Analyzing commits...")

	// analyze commits
	err = analyze.Analyze(result, analyzer)
}
