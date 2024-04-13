package semantic

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-git/go-git/v5/config"

	"github.com/go-git/go-git/v5"
	"github.com/rs/zerolog/log"
	"github.com/wintbiit/semantic-release-go/analyze"
	"github.com/wintbiit/semantic-release-go/output"
	"github.com/wintbiit/semantic-release-go/types"

	"github.com/wintbiit/semantic-release-go/utils"
)

func Run(opt types.SemanticOptions) {
	r, err := git.PlainOpen(opt.Path)
	if err != nil {
		log.Fatal().Err(err).Msg("Not a git repository")
	}

	result := &types.Result{
		Season:  opt.Season,
		Channel: opt.Channel,
		Repo:    opt.Repo,
		Built:   time.Now(),
	}

	tags, err := r.Tags()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to get tags")
	}

	scannedTags, err := utils.History(tags, opt.Season, opt.Channel)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to scan tags")
	}

	currentCommit, err := r.Head()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to get current commit")
	}

	log.Info().Msgf("Current commit: %v", utils.HashShort(currentCommit.Hash()))
	log.Info().Msgf("Using channel: %s, season: %s", opt.Channel, opt.Season)

	if len(scannedTags) == 0 {
		log.Info().Msgf("No history of %s %s, will use vcs tree tail and release first version v1.0.0", opt.Season, opt.Channel)
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
	result.NextRelease.Version.Channel = opt.Channel
	result.NextRelease.Version.Season = opt.Season

	// analyze commits
	if err = analyze.Analyze(result, &opt, opt.Analyzer); err != nil {
		log.Fatal().Err(err).Msg("Failed to analyze commits")
	}

	if result.NextRelease.Version.SameFrom(result.LatestRelease.Version) {
		log.Info().Msg("No new version to release")
		return
	}

	log.Info().Str("next_release", result.NextRelease.String()).Str("release_type", result.ReleaseType).Msg("New version to release")

	// output result
	outputs := output.Output(result, &opt)
	log.Info().Msg("Output done")
	j, _ := json.MarshalIndent(outputs, "", "  ")
	fmt.Println(string(j))

	if !opt.Dry {
		if opt.Tag {
			// create tag
			if _, err = r.CreateTag(result.NextRelease.Version.Tag(), result.NextRelease.Hash(), &git.CreateTagOptions{
				Message: result.NextRelease.Version.String(),
			}); err != nil {
				log.Fatal().Err(err).Msg("Failed to create tag")
			}
		}

		if opt.Push {
			// push tag
			err = r.Push(&git.PushOptions{
				RemoteName: "origin",
				RefSpecs: []config.RefSpec{
					config.RefSpec("refs/tags/" + result.NextRelease.Version.Tag() + ":refs/tags/" + result.NextRelease.Version.Tag()),
				},
			})

			if err != nil {
				log.Fatal().Err(err).Msg("Failed to push tag")
			}
		}
	}
}
