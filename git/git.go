package git

import "os/exec"

type Repository struct {
	Path string
}

func Open(path string) (*Repository, error) {
	if path == "" {
		path = "."
	}

	if _, err := exec.LookPath("git"); err != nil {
		return nil, err
	}

	c := exec.Command("git", "rev-parse", "--show-toplevel")
	c.Dir = path
	if err := c.Run(); err != nil {
		return nil, err
	}

	return &Repository{Path: path}, nil
}

func (r *Repository) execute(subcommand string, args ...string) (string, error) {
	c := exec.Command("git", append([]string{subcommand}, args...)...)
	c.Dir = r.Path

	out, err := c.CombinedOutput()
	if err != nil {
		return "", err
	}

	return string(out), nil
}
