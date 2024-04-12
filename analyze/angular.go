//go:build analyzer_angular

package analyze

type AngularAnalyzer struct{}

func init() {
	RegisterAnalyzer("angular", &AngularAnalyzer{})
}
