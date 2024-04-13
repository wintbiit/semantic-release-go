package analyze

import (
	"errors"
	"fmt"

	"github.com/wintbiit/semantic-release-go/types"
)

type IAnalyzer interface {
	Analyze(result *types.Result) error
}

var registeredAnalyzers = make(map[string]IAnalyzer)

func RegisterAnalyzer(name string, analyzer IAnalyzer) {
	registeredAnalyzers[name] = analyzer
}

func GetAnalyzer(name string) IAnalyzer {
	return registeredAnalyzers[name]
}

func Analyze(result *types.Result, analyzer string) error {
	a := GetAnalyzer(analyzer)
	if a == nil {
		return errors.New(fmt.Sprintf("Analyzer %s not found", analyzer))
	}

	return a.Analyze(result)
}
