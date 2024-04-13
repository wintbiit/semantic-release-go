package output

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/wintbiit/semantic-release-go/types"
	"os"
)

type GithubOutput struct{}

var githubState = os.Getenv("GITHUB_STATE")
var githubOutput = os.Getenv("GITHUB_OUTPUT")
var stateFile *os.File
var outputFile *os.File

func (o *GithubOutput) Output(result *types.Result) error {
	log.Info().Msg("Outputting to github")

	setState("channel", result.Channel)
	setState("release.type", result.ReleaseType)
	setState("season", result.Season)
	setState("repo", result.Repo)
	setOutput("built", result.Built.Format("2006-01-02 15:04:05"))
	setOutput("release.next.version", result.NextRelease.Version.ShortString())
	setOutput("release.next.hash", result.NextRelease.Hash().String())
	setOutput("release.next.major", fmt.Sprintf("%d", result.NextRelease.Major))
	setOutput("release.next.minor", fmt.Sprintf("%d", result.NextRelease.Minor))
	setOutput("release.next.patch", fmt.Sprintf("%d", result.NextRelease.Patch))
	if result.LatestRelease.Reference != nil {
		setOutput("release.latest.version", result.LatestRelease.Version.ShortString())
		setOutput("release.latest.hash", result.LatestRelease.Hash().String())
		setOutput("release.latest.major", fmt.Sprintf("%d", result.LatestRelease.Major))
		setOutput("release.latest.minor", fmt.Sprintf("%d", result.LatestRelease.Minor))
		setOutput("release.latest.patch", fmt.Sprintf("%d", result.LatestRelease.Patch))
	}

	return nil
}

func init() {
	if githubState == "" {
		log.Warn().Msg("GITHUB_STATE not set, github output will be disabled")
		return
	}

	if githubOutput == "" {
		log.Warn().Msg("GITHUB_OUTPUT not set, github output will be disabled")
		return
	}

	var err error
	stateFile, err = os.OpenFile(githubState, os.O_WRONLY|os.O_APPEND, 0o666)
	if err != nil {
		log.Error().Err(err).Msg("Failed to open GITHUB_STATE file")
		return
	}

	outputFile, err = os.OpenFile(githubOutput, os.O_WRONLY|os.O_APPEND, 0o666)
	if err != nil {
		log.Error().Err(err).Msg("Failed to open GITHUB_OUTPUT file")
		return
	}

	RegisterOutput("github", &GithubOutput{})
}

func setState(key, value string) {
	if stateFile == nil {
		return
	}

	_, err := stateFile.WriteString(key + "=" + value + "\n")
	if err != nil {
		log.Error().Err(err).Msg("Failed to write to GITHUB_STATE file")
	}
}

func setOutput(key, value string) {
	if outputFile == nil {
		return
	}

	_, err := outputFile.WriteString(key + "=" + value + "\n")
	if err != nil {
		log.Error().Err(err).Msg("Failed to write to GITHUB_OUTPUT file")
	}
}
