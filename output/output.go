package output

import (
	"sync"

	"github.com/rs/zerolog/log"

	"github.com/wintbiit/semantic-release-go/types"
)

type IOutput interface {
	Output(result *types.Result, opt *types.SemanticOptions) error
}

var registeredOutputs = make(map[string]IOutput)

func RegisterOutput(name string, output IOutput) {
	registeredOutputs[name] = output
}

func GetOutput(name string) IOutput {
	return registeredOutputs[name]
}

func Output(result *types.Result, opt *types.SemanticOptions) map[string]bool {
	ok := make(map[string]bool, len(registeredOutputs))
	var wg sync.WaitGroup
	wg.Add(len(registeredOutputs))
	for name, output := range registeredOutputs {
		go func(name string, output IOutput) {
			defer wg.Done()
			log.Info().Msgf("Outputting result using %s", name)
			err := output.Output(result, opt)
			ok[name] = true
			if err != nil {
				log.Error().Err(err).Msgf("Failed to output using %s", name)
				ok[name] = false
			}
		}(name, output)
	}

	wg.Wait()

	return ok
}
