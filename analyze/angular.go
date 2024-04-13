//go:build analyzer_angular

package analyze

import "github.com/wintbiit/semantic-release-go/semantic"

type AngularAnalyzer struct{}

func (a *AngularAnalyzer) Analyze(result *semantic.Result) error {
	return nil
}

func init() {
	RegisterAnalyzer("angular", &AngularAnalyzer{})
}
