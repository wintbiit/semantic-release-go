//go:build output_teamcity

package output

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/wintbiit/semantic-release-go/utils"

	"github.com/wintbiit/semantic-release-go/types"
)

type TeamcityOutput struct{}

func (o *TeamcityOutput) Output(result *types.Result, _ *types.SemanticOptions) error {
	log.Info().Msg("Output to teamcity")

	teamCityTag(result.Channel)
	teamCityTag(result.Season)
	teamCityTag(result.ReleaseType)

	teamCityParam("channel", result.Channel)
	teamCityParam("release.type", result.ReleaseType)
	teamCityParam("season", result.Season)
	teamCityParam("repo", result.Repo)
	teamCityParam("built", result.Built.Format("2006-01-02 15:04:05"))
	teamCityParam("release.next.version", result.NextRelease.Version.ShortString())
	teamCityParam("release.next.hash", utils.HashShort(result.NextRelease.Hash()))
	teamCityParam("release.next.major", fmt.Sprintf("%d", result.NextRelease.Major))
	teamCityParam("release.next.minor", fmt.Sprintf("%d", result.NextRelease.Minor))
	teamCityParam("release.next.patch", fmt.Sprintf("%d", result.NextRelease.Patch))
	if result.LatestRelease.Reference != nil {
		teamCityParam("release.latest.version", result.LatestRelease.Version.ShortString())
		teamCityParam("release.latest.hash", utils.HashShort(result.LatestRelease.Hash()))
		teamCityParam("release.latest.major", fmt.Sprintf("%d", result.LatestRelease.Major))
		teamCityParam("release.latest.minor", fmt.Sprintf("%d", result.LatestRelease.Minor))
		teamCityParam("release.latest.patch", fmt.Sprintf("%d", result.LatestRelease.Patch))
	}

	log.Info().Msg("Teamcity output done")
	return nil
}

func init() {
	RegisterOutput("teamcity", &TeamcityOutput{})
}

func teamCityParam(key, value string) {
	fmt.Printf("##teamcity[setParameter name='%s' value='%s']\n", key, value)
}

func teamCityTag(tag string) {
	fmt.Printf("##teamcity[addBuildTag name='%s']\n", tag)
}
