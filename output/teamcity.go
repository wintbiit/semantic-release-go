//go:build output_teamcity

package output

import "github.com/wintbiit/semantic-release-go/semantic"

type TeamcityOutput struct{}

func (o *TeamcityOutput) Output(result *semantic.Result) {
}

func init() {
	RegisterOutput("teamcity", &TeamcityOutput{})
}
