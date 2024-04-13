package git

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"
)

type Commit struct {
	Hash    string
	Message string
	Author  Author
	Date    time.Time
}

var gitDateLayout = "Mon Jan 2 15:04:05 2006 -0700"

func (c *Commit) String() string {
	return fmt.Sprintf("%s %s @ %s - %s", c.Hash, c.Author, c.Date, c.Message)
}

type Author struct {
	Name  string
	Email string
}

func (a *Author) String() string {
	return fmt.Sprintf("%s <%s>", a.Name, a.Email)
}

const gitLogCommitFormat = "--pretty=format:%H|%an|%ae|%ad|%s"

func (r *Repository) CommitsBetween(expr string) ([]*Commit, error) {
	out, err := r.execute("log", gitLogCommitFormat, expr)
	if err != nil {
		return nil, err
	}

	sp := splitLines(out)
	commits := make([]*Commit, 0, len(sp))
	var wg sync.WaitGroup
	var mu sync.Mutex
	wg.Add(len(sp))
	for _, line := range sp {
		go func(line string) {
			defer wg.Done()
			commit := parseCommit(line)
			if commit == nil {
				return
			}

			mu.Lock()
			commits = append(commits, commit)
			mu.Unlock()
		}(line)
	}

	wg.Wait()

	sort.Slice(commits, func(i, j int) bool {
		return commits[i].Date.After(commits[j].Date)
	})

	return commits, nil
}

func (r *Repository) CommitsSince(hash string) ([]*Commit, error) {
	return r.CommitsBetween(hash + "..HEAD")
}

func (r *Repository) Commits() ([]*Commit, error) {
	return r.CommitsBetween("--all")
}

func (r *Repository) LastCommit() (*Commit, error) {
	out, err := r.execute("log", "-1", gitLogCommitFormat)
	if err != nil {
		return nil, err
	}

	return parseCommit(splitLines(out)[0]), nil
}

func (r *Repository) Commit(hash string) (*Commit, error) {
	out, err := r.execute("log", "-1", gitLogCommitFormat, hash)
	if err != nil {
		return nil, err
	}

	return parseCommit(splitLines(out)[0]), nil
}

func (r *Repository) PushCommit(hash string) error {
	_, err := r.execute("push", "origin", hash)
	return err
}

func split(s, sep string) []string {
	return strings.Split(s, sep)
}

func splitLines(s string) []string {
	return split(s, "\n")
}

func parseCommit(s string) *Commit {
	parts := split(s, "|")
	if len(parts) != 5 {
		return nil
	}

	date, err := time.Parse(gitDateLayout, parts[3])
	if err != nil {
		return nil
	}

	return &Commit{
		Hash:    parts[0],
		Author:  Author{Name: parts[1], Email: parts[2]},
		Date:    date,
		Message: parts[4],
	}
}
