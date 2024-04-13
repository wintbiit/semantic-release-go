package output

import "github.com/wintbiit/semantic-release-go/semantic"

type EnvOutput struct{}

func (o *EnvOutput) Output(result *semantic.Result) error {
	return nil
}

func init() {
	RegisterOutput("env", &EnvOutput{})
}
