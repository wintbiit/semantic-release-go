package types

import (
	"fmt"
	"time"

	"github.com/wintbiit/semantic-release-go/git"
)

const (
	ChannelInsider = "insider"
	ChannelAlpha   = "alpha"
	ChannelBeta    = "beta"
	ChannelRelease = "release"

	ReleaseTypeMajor = "major"
	ReleaseTypeMinor = "minor"
	ReleaseTypePatch = "patch"
)

type SemanticOptions struct {
	Path      string
	Channel   string
	Branch    string
	Analyzer  string
	Repo      string
	Dry       bool
	Tag       bool
	Push      bool
	Changelog string
}

type Result struct {
	Repo          string
	Branch        string
	Channel       string
	Built         time.Time
	NewRelease    bool
	NextRelease   SemverTag
	LatestRelease SemverTag
	ReleaseType   string
	Commits       []*git.Commit
	ReleaseNotes  map[string][]ReleaseNote
}

type ReleaseNote struct {
	Commit *git.Commit
	Scope  string
	Desc   string
}

func (n *ReleaseNote) Describe(repo string) string {
	commitUrl := fmt.Sprintf("%s/commit/%s", repo, n.Commit.Hash)

	if n.Scope != "" {
		return fmt.Sprintf("- [%s](%s) [**%s**]: %s @[%s](mailto://%s)\n", n.Commit.Hash[0:7], commitUrl, n.Scope, n.Desc, n.Commit.Author.Name, n.Commit.Author.Email)
	} else {
		return fmt.Sprintf("- [%s](%s) %s @[%s](mailto://%s)\n", n.Commit.Hash[0:7], commitUrl, n.Desc, n.Commit.Author.Name, n.Commit.Author.Email)
	}
}
