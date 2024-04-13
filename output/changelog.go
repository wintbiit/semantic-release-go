package output

import (
	"fmt"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/wintbiit/semantic-release-go/types"
)

type ChangeLogOutput struct{}

var fileName = "CHANGELOG.md"

func (o *ChangeLogOutput) Output(result *types.Result) error {
	log.Info().Msg("Outputting changelog")
	f, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(fmt.Sprintf("# Changelog v%d.%d.%d\n\n", result.LatestRelease.Major, result.LatestRelease.Minor, result.LatestRelease.Patch))
	if err != nil {
		return err
	}

	_, err = f.WriteString(fmt.Sprintf("> `%s` `%s` at %s\n", result.Season, result.Channel, result.Built.Format("2006-01-02 15:04:05")))

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

	if os.Getenv("CHANGELOG_FILE") != "" {
		fileName = os.Getenv("CHANGELOG_FILE")
	}
}
