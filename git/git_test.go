package git_test

import (
	"testing"

	"github.com/wintbiit/semantic-release-go/git"
)

func TestGit(t *testing.T) {
	repo, err := git.Open(".")
	if err != nil {
		t.Error(err)
	}

	since, err := repo.CommitsSince("53243a30f327b2685df88564b08fcf8a3eac84d1")
	if err != nil {
		t.Error(err)
		return
	}

	for _, commit := range since {
		t.Log(commit.String())
	}

	tags, err := repo.Tags()
	if err != nil {
		t.Error(err)
		return
	}

	for _, tag := range tags {
		t.Log(tag.String())
	}
}
