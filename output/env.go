package output

import (
	"fmt"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/wintbiit/semantic-release-go/types"
)

type EnvOutput struct{}

func (o *EnvOutput) Output(result *types.Result, _ *types.SemanticOptions) error {
	log.Info().Msg("Outputting to env")

	os.Setenv("RELEASE_CHANNEL", result.Channel)
	os.Setenv("RELEASE_TYPE", result.ReleaseType)
	os.Setenv("RELEASE_BRANCH", result.Branch)
	os.Setenv("RELEASE_REPO", result.Repo)
	os.Setenv("RELEASE_BUILT", result.Built.Format("2006-01-02 15:04:05"))
	os.Setenv("RELEASE_NEXT_VERSION", result.NextRelease.Version.ShortString())
	os.Setenv("RELEASE_NEXT_HASH", result.NextRelease.Hash().String())
	os.Setenv("RELEASE_NEXT_MAJOR", fmt.Sprintf("%d", result.NextRelease.Major))
	os.Setenv("RELEASE_NEXT_MINOR", fmt.Sprintf("%d", result.NextRelease.Minor))
	os.Setenv("RELEASE_NEXT_PATCH", fmt.Sprintf("%d", result.NextRelease.Patch))
	if result.LatestRelease.Reference != nil {
		os.Setenv("RELEASE_LATEST_VERSION", result.LatestRelease.Version.ShortString())
		os.Setenv("RELEASE_LATEST_HASH", result.LatestRelease.Hash().String())
		os.Setenv("RELEASE_LATEST_MAJOR", fmt.Sprintf("%d", result.LatestRelease.Major))
		os.Setenv("RELEASE_LATEST_MINOR", fmt.Sprintf("%d", result.LatestRelease.Minor))
		os.Setenv("RELEASE_LATEST_PATCH", fmt.Sprintf("%d", result.LatestRelease.Patch))
	}

	log.Info().Msg("Env output done")

	return nil
}

func init() {
	RegisterOutput("env", &EnvOutput{})
}
