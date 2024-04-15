package output

import (
	"fmt"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/wintbiit/semantic-release-go/types"
)

type ChangeLogOutput struct{}

func (o *ChangeLogOutput) Output(result *types.Result, opt *types.SemanticOptions) error {
	if opt.Dry || opt.Changelog == "" {
		log.Info().Msg("Changelog output disabled")
		return nil
	}
	log.Info().Msg("Outputting changelog")

	flag := os.O_CREATE | os.O_WRONLY
	if opt.ChangelogAppend {
		flag |= os.O_APPEND
	} else {
		flag |= os.O_TRUNC
	}
	f, err := os.OpenFile(opt.Changelog, flag, 0o666)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(fmt.Sprintf("# Changelog %s\n\n", result.LatestRelease.Version.ShortString()))
	if err != nil {
		return err
	}

	_, err = f.WriteString(fmt.Sprintf("> `%s` `%s` at %s\n", result.Branch, result.Channel, result.Built.Format("2006-01-02 15:04:05")))

	for title, notes := range result.ReleaseNotes {
		log.Info().Msgf("Writing %s", title)
		_, err = f.WriteString(fmt.Sprintf("## %s\n\n", title))
		if err != nil {
			return err
		}

		for _, note := range notes {
			_, err = f.WriteString(note.Describe(result.Repo))
			if err != nil {
				return err
			}
		}
		_, err = f.WriteString("\n")
		if err != nil {
			return err
		}
	}

	_, err = f.WriteString("\n")

	log.Info().Msg("Changelog written")
	return nil
}

func init() {
	RegisterOutput("changelog", &ChangeLogOutput{})
}
