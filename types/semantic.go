package types

import (
	"github.com/go-git/go-git/v5/plumbing/object"
)

const (
	ChannelInsider = "insider"
	ChannelAlpha   = "alpha"
	ChannelBeta    = "beta"
	ChannelRelease = "release"
)

type Result struct {
	Season        string
	Channel       string
	NextRelease   SemverTag
	LatestRelease SemverTag
	Commits       []*object.Commit
}
