package output

import "github.com/wintbiit/semantic-release-go/semantic"

type IOutput interface {
	Output(result *semantic.Result)
}

var registeredOutputs = make(map[string]IOutput)

func RegisterOutput(name string, output IOutput) {
	registeredOutputs[name] = output
}

func GetOutput(name string) IOutput {
	return registeredOutputs[name]
}
