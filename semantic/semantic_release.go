package semantic

import (
	"encoding/json"
	"fmt"
	"github.com/wintbiit/semantic-release-go/output"
	"time"

	"github.com/wintbiit/semantic-release-go/git"

	"github.com/rs/zerolog/log"
	"github.com/wintbiit/semantic-release-go/analyze"
	"github.com/wintbiit/semantic-release-go/types"

	"github.com/wintbiit/semantic-release-go/utils"
)

func Run(opt types.SemanticOptions) {
	r, err := git.Open(opt.Path)
	if err != nil {
		log.Fatal().Err(err).Msg("Not a git repository")
	}

	result := &types.Result{
		Branch:     opt.Branch,
		Channel:    opt.Channel,
		Repo:       opt.Repo,
		NewRelease: false,
		Built:      time.Now(),
	}

	defer func() {
		// output result
		outputs := output.Output(result, &opt)
		log.Info().Msg("Output done")
		j, _ := json.MarshalIndent(outputs, "", "  ")
		fmt.Println(string(j))
	}()

	tags, err := r.Tags()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to get tags")
	}

	scannedTags, err := utils.History(tags, opt.Branch, opt.Channel)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to scan tags")
	}

	currentCommit, err := r.LastCommit()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to get current commit")
	}

	log.Info().Msgf("Current commit: %v", utils.HashShort(currentCommit.Hash))
	log.Info().Msgf("Using channel: %s, branch: %s", opt.Channel, opt.Branch)

	initRelease := len(scannedTags) == 0
	if initRelease {
		log.Info().Msgf("No history of %s %s, will use vcs tree tail and release first version v1.0.0", opt.Branch, opt.Channel)
		result.Commits, err = r.Commits()
	} else {
		result.LatestRelease = scannedTags[0]
		log.Info().Msgf("Last release: %s", result.LatestRelease.String())
		result.Commits, err = r.CommitsSince(result.LatestRelease.Hash)
	}

	if err != nil {
		log.Fatal().Err(err).Msg("Failed to get commits since last version")
	}

	if len(result.Commits) == 0 {
		log.Info().Msg("No commits since last version")
		if !initRelease {
			log.Info().Msg("")
			result.NextRelease = result.LatestRelease
		}
		return
	}

	log.Info().Msg("Analyzing commits...")

	result.NextRelease = types.SemverTag{
		Commit:  currentCommit,
		Version: result.LatestRelease.Version,
	}
	result.NextRelease.Version.Channel = opt.Channel
	result.NextRelease.Version.Branch = opt.Branch

	// analyze commits
	if err = analyze.Analyze(result, &opt, opt.Analyzer); err != nil {
		log.Fatal().Err(err).Msg("Failed to analyze commits")
	}

	if !initRelease {
		if result.NextRelease.Version.SameFrom(result.LatestRelease.Version) {
			log.Info().Msg("No new version to release")
			if !initRelease {
				log.Info().Msg("")
				result.NextRelease = result.LatestRelease
			}
			return
		}
	} else {
		result.NextRelease.Version.Major = 1
		result.NextRelease.Version.Minor = 0
		result.NextRelease.Version.Patch = 0
	}

	result.NewRelease = true

	log.Info().Str("next_release", result.NextRelease.String()).Str("release_type", result.ReleaseType).Msg("New version to release")

	if !opt.Dry {
		if opt.Tag {
			if _, err = r.CreateTag(
				result.NextRelease.Version.Tag(),
				result.NextRelease.Hash,
				result.NextRelease.Version.String()); err != nil {
				log.Fatal().Err(err).Msg("Failed to create tag")
			}
		}

		if opt.Push {
			if err = r.PushCommit(result.NextRelease.Hash); err != nil {
				log.Fatal().Err(err).Msg("Failed to push tag")
			}
		}
	}
}
