package utils

import "github.com/go-git/go-git/v5/plumbing/object"

func CommitsSince(commits object.CommitIter, hash string) ([]*object.Commit, error) {
	var sinceCommits []*object.Commit
	var commit *object.Commit
	var err error

	for {
		commit, err = commits.Next()
		if err != nil {
			break
		}

		if commit.Hash.String() == hash {
			break
		}

		sinceCommits = append(sinceCommits, commit)
	}

	return sinceCommits, nil
}
