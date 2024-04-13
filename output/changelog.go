package output

import "github.com/wintbiit/semantic-release-go/semantic"

type ChangeLogOutput struct{}

func (o *ChangeLogOutput) Output(result *semantic.Result) error {
	return nil
}

func init() {
	RegisterOutput("changelog", &ChangeLogOutput{})
}
