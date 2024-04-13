package output

import (
	"errors"
	"sync"

	"github.com/rs/zerolog/log"

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
	errorM := make(map[string]error, len(registeredOutputs))
	var wg sync.WaitGroup
	wg.Add(len(registeredOutputs))
	for name, output := range registeredOutputs {
		go func(name string, output IOutput) {
			defer wg.Done()
			log.Info().Msgf("Outputting result using %s", name)
			err := output.Output(result)
			if err != nil {
				errorM[name] = err
			}
		}(name, output)
	}

	wg.Wait()

	if len(errorM) > 0 {
		err := "Failed to output using: "
		for name, e := range errorM {
			err += name + ": " + e.Error() + ", "
		}

		return errors.New(err)
	}

	return nil
}
