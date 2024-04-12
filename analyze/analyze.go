package analyze

type IAnalyzer interface{}

var registeredAnalyzers = make(map[string]IAnalyzer)

func RegisterAnalyzer(name string, analyzer IAnalyzer) {
	registeredAnalyzers[name] = analyzer
}

func GetAnalyzer(name string) IAnalyzer {
	return registeredAnalyzers[name]
}
