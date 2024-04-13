//go:build analyzer_angular

package analyze

import (
	"github.com/wintbiit/semantic-release-go/types"
)

type AngularAnalyzer struct{}

func (a *AngularAnalyzer) Analyze(result *types.Result) error {
	return nil
}

func init() {
	RegisterAnalyzer("angular", &AngularAnalyzer{})
}
