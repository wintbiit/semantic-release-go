package output

import "github.com/wintbiit/semantic-release-go/types"

type ChangeLogOutput struct{}

func (o *ChangeLogOutput) Output(result *types.Result) error {
	return nil
}

func init() {
	RegisterOutput("changelog", &ChangeLogOutput{})
}
