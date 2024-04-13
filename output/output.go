package output

import (
	"log"

	"github.com/wintbiit/semantic-release-go/semantic"
)

type IOutput interface {
	Output(result *semantic.Result) error
}

var registeredOutputs = make(map[string]IOutput)

func RegisterOutput(name string, output IOutput) {
	registeredOutputs[name] = output
}

func GetOutput(name string) IOutput {
	return registeredOutputs[name]
}

func Output(result *semantic.Result) error {
	for name, output := range registeredOutputs {
		log.Printf("Outputting result using %s", name)
		err := output.Output(result)
		if err != nil {
			return err
		}
	}

	return nil
}
