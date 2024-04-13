package output

import "github.com/wintbiit/semantic-release-go/types"

type EnvOutput struct{}

func (o *EnvOutput) Output(result *types.Result) error {
	return nil
}

func init() {
	RegisterOutput("env", &EnvOutput{})
}
