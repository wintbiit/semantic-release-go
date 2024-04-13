//go:build output_teamcity

package output

import (
	"fmt"

	"github.com/wintbiit/semantic-release-go/types"
)

type TeamcityOutput struct{}

func (o *TeamcityOutput) Output(result *types.Result) error {
	return nil
}

func init() {
	RegisterOutput("teamcity", &TeamcityOutput{})
}

func teamCityParam(key, value string) {
	fmt.Printf("##teamcity[setParameter name='%s' value='%s']", key, value)
}
