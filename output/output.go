package output

import (
	"log"

	"github.com/wintbiit/semantic-release-go/types"
)

type IOutput interface {
	Output(result *types.Result) error
}

var registeredOutputs = make(map[string]IOutput)

func RegisterOutput(name string, output IOutput) {
	registeredOutputs[name] = output
}

func GetOutput(name string) IOutput {
	return registeredOutputs[name]
}

func Output(result *types.Result) error {
	for name, output := range registeredOutputs {
		log.Printf("Outputting result using %s", name)
		err := output.Output(result)
		if err != nil {
			return err
		}
	}

	return nil
}
