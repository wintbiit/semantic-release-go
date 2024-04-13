package semantic

import (
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/rs/zerolog/log"
	"github.com/wintbiit/semantic-release-go/analyze"
	"github.com/wintbiit/semantic-release-go/output"
	"github.com/wintbiit/semantic-release-go/types"

	"github.com/wintbiit/semantic-release-go/utils"
)

func Run(path, channel, season, analyzer, repo string) {
	r, err := git.PlainOpen(path)
	if err != nil {
		log.Fatal().Err(err).Msg("Not a git repository")
	}

	result := &types.Result{
		Season:  season,
		Channel: channel,
		Repo:    repo,
		Built:   time.Now(),
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
		log.Info().Msgf("Last release: %s", utils.HashShort(result.LatestRelease))
	}

	// get commits since last version
	commits, err := r.Log(&git.LogOptions{From: currentCommit.Hash()})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to get commits")
	}

	sinceCommits, err := utils.CommitsSince(commits, utils.HashShort(currentCommit.Hash()))
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to get commits since last version")
	}

	result.Commits = sinceCommits
	if len(result.Commits) == 0 {
		log.Info().Msg("No commits since last version")
		return
	}

	log.Info().Msg("Analyzing commits...")

	result.NextRelease = types.SemverTag{
		Reference: currentCommit,
		Version:   result.LatestRelease.Version,
	}

	// analyze commits
	if err = analyze.Analyze(result, analyzer); err != nil {
		log.Fatal().Err(err).Msg("Failed to analyze commits")
	}

	if result.NextRelease.Version.SameFrom(result.LatestRelease.Version) {
		log.Info().Msg("No new version to release")
		return
	}

	log.Info().Str("next_release", result.NextRelease.String()).Str("release_type", result.ReleaseType).Msg("New version to release")

	// output result
	if err = output.Output(result); err != nil {
		log.Fatal().Err(err).Msg("Failed to output result")
	}
}
