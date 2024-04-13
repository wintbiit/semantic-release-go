package semantic

import (
	"log"

	"github.com/go-git/go-git/v5/plumbing/object"

	"github.com/go-git/go-git/v5"

	"github.com/wintbiit/semantic-release-go/utils"
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
	NextRelease   utils.Version
	LatestRelease utils.Version
	Commits       []*object.Commit
}

func Run(path, channel, season, analyzer string) {
	r, err := git.PlainOpen(path)
	if err != nil {
		log.Fatalf("Not a git repository: %v", err)
		return
	}

	tags, err := r.Tags()
	if err != nil {
		log.Fatalf("Failed to get tags: %v", err)
		return
	}

	scannedTags, err := utils.History(tags, season, channel)
	if err != nil {
		log.Fatalf("Failed to scan tags: %v", err)
		return
	}

	currentCommit, err := r.Head()
	if err != nil {
		log.Fatalf("Failed to get current commit: %v", err)
		return
	}

	log.Printf("Current commit: %v", utils.HashShort(currentCommit.Hash()))
	log.Printf("Using channel: %s, season: %s", channel, season)

	var lastVersion utils.SemverTag

	if len(scannedTags) == 0 {
		log.Printf("No history of %s %s, will use vcs tree tail and release first version v1.0.0", season, channel)
	} else {
		lastVersion = scannedTags[len(scannedTags)-1]
		log.Printf("Last release: %s", lastVersion.Version.String())
	}

	// get commits since last version
	commits, err := r.Log(&git.LogOptions{From: currentCommit.Hash()})
	if err != nil {
		log.Fatalf("Failed to get commits: %v", err)
		return
	}

	sinceCommits, err := utils.CommitsSince(commits, lastVersion.Hash().String())
	if err != nil {
		log.Fatalf("Failed to get commits since last version: %v", err)
		return
	}

	log.Printf("Commits since last version: %d", len(sinceCommits))
}
