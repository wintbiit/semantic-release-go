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
	teamCityTag(result.Branch)
	teamCityTag(result.ReleaseType)

	teamCityParam("semantic.release.channel", result.Channel)
	teamCityParam("semantic.release.type", result.ReleaseType)
	teamCityParam("semantic.release.branch", result.Branch)
	teamCityParam("semantic.release.repo", result.Repo)
	teamCityParam("semantic.release.built", result.Built.Format("2006-01-02 15:04:05"))
	teamCityParam("semantic.release.next.version", result.NextRelease.Version.ShortString())
	teamCityParam("semantic.release.next.hash", utils.HashShort(result.NextRelease.Hash()))
	teamCityParam("semantic.release.next.major", fmt.Sprintf("%d", result.NextRelease.Major))
	teamCityParam("semantic.release.next.minor", fmt.Sprintf("%d", result.NextRelease.Minor))
	teamCityParam("semantic.release.next.patch", fmt.Sprintf("%d", result.NextRelease.Patch))
	teamCityParam("semantic.release.next.tag", result.NextRelease.Tag())
	if result.LatestRelease.Reference != nil {
		teamCityParam("semantic.release.latest.version", result.LatestRelease.Version.ShortString())
		teamCityParam("semantic.release.latest.hash", utils.HashShort(result.LatestRelease.Hash()))
		teamCityParam("semantic.release.latest.major", fmt.Sprintf("%d", result.LatestRelease.Major))
		teamCityParam("semantic.release.latest.minor", fmt.Sprintf("%d", result.LatestRelease.Minor))
		teamCityParam("semantic.release.latest.patch", fmt.Sprintf("%d", result.LatestRelease.Patch))
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
	fmt.Printf("##teamcity[addBuildTag '%s']\n", tag)
}
