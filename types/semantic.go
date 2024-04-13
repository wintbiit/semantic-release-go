package types

import (
	"fmt"
	"time"

	"github.com/go-git/go-git/v5/plumbing/object"
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
	Season    string
	Analyzer  string
	Repo      string
	Dry       bool
	Tag       bool
	Push      bool
	Changelog bool
}

type Result struct {
	Repo          string
	Season        string
	Channel       string
	Built         time.Time
	NextRelease   SemverTag
	LatestRelease SemverTag
	ReleaseType   string
	Commits       []*object.Commit
	ReleaseNotes  map[string][]ReleaseNote
}

type ReleaseNote struct {
	Commit *object.Commit
	Scope  string
	Desc   string
}

func (n *ReleaseNote) Describe(repo string) string {
	commitUrl := fmt.Sprintf("%s/commit/%s", repo, n.Commit.Hash.String())

	if n.Scope != "" {
		return fmt.Sprintf("- [%s](%s) [**%s**]: %s @[%s](mailto://%s)\n", n.Commit.Hash.String()[0:7], commitUrl, n.Scope, n.Desc, n.Commit.Author.Name, n.Commit.Author.Email)
	} else {
		return fmt.Sprintf("- [%s](%s) %s @[%s](mailto://%s)\n", n.Commit.Hash.String()[0:7], commitUrl, n.Desc, n.Commit.Author.Name, n.Commit.Author.Email)
	}
}
