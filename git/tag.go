package git

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"
)

type Tag struct {
	Name string
	*Commit
}

func (t *Tag) String() string {
	return fmt.Sprintf("%s %s", t.Name, t.Commit)
}

func (r *Repository) Tags() ([]*Tag, error) {
	out, err := r.execute(
		"tag",
		"-l",
		"--format=%(refname:short)|%(objectname)|%(subject)|%(creatordate)|%(creator)",
	)
	if err != nil {
		return nil, err
	}

	sp := splitLines(out)
	tags := make([]*Tag, 0, len(sp))
	var wg sync.WaitGroup
	var mu sync.Mutex
	wg.Add(len(sp))
	for _, line := range sp {
		go func(line string) {
			defer wg.Done()
			tag := r.parseTag(line)
			if tag == nil {
				return
			}

			mu.Lock()
			tags = append(tags, tag)
			mu.Unlock()
		}(line)
	}

	wg.Wait()

	sort.Slice(tags, func(i, j int) bool {
		return tags[i].Date.After(tags[j].Date)
	})

	return tags, nil
}

func (r *Repository) CreateTag(name, hash, message string) (*Tag, error) {
	_, err := r.execute("tag", "-a", name, hash, "-m", message)
	if err != nil {
		return nil, err
	}

	commit, err := r.Commit(hash)
	if err != nil {
		return nil, err
	}

	return &Tag{Name: name, Commit: commit}, nil
}

func (r *Repository) DeleteTag(name string) error {
	_, err := r.execute("tag", "-d", name)
	return err
}

func (r *Repository) parseTag(s string) *Tag {
	sp := split(s, "|")
	if len(sp) != 5 {
		return nil
	}

	var tag Tag
	date, err := time.Parse(gitDateLayout, sp[3])
	if err != nil {
		return nil
	}

	tag.Name = sp[0]
	tag.Commit = &Commit{
		Hash:    sp[1],
		Message: sp[2],
		Date:    date,
	}

	author := strings.Split(sp[4], " ")
	tag.Commit.Author = Author{
		Name:  author[0],
		Email: author[1],
	}

	tag.Commit.Author.Email = strings.Trim(tag.Commit.Author.Email, "<>")
	return &tag
}
