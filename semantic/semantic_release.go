package semantic

import (
	"log"

	"github.com/go-git/go-git/v5"

	"github.com/wintbiit/semantic-release-go/utils"
)

type Result struct {
	Season  string
	Channel string `env:"CHANNEL"`
}

func Run(path string) {
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

	scannedTags, err := utils.ValidTags(tags, "2024uc", "insider")
	if err != nil {
		log.Fatalf("Failed to scan tags: %v", err)
		return
	}

	utils.SortTags(scannedTags)

	currentCommit, err := r.Head()
	if err != nil {
		log.Fatalf("Failed to get current commit: %v", err)
		return
	}

	log.Printf("Current commit: %v", utils.HashShort(currentCommit.Hash()))
	log.Printf("Using channel: %s, season: %s", CHANNEL, SEASON)
}
