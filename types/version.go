package types

import (
	"fmt"

	"github.com/go-git/go-git/v5/plumbing"
)

type Version struct {
	Branch  string
	Channel string
	Major   int
	Minor   int
	Patch   int
}

type SemverTag struct {
	*plumbing.Reference
	Version
}

func (v *Version) String() string {
	return fmt.Sprintf("%s %s %s", v.Branch, v.Channel, v.ShortString())
}

func (v *Version) ShortString() string {
	return fmt.Sprintf("v%d.%d.%d", v.Major, v.Minor, v.Patch)
}

func (v *Version) Tag() string {
	return fmt.Sprintf("%s/%s/v%d.%d.%d", v.Branch, v.Channel, v.Major, v.Minor, v.Patch)
}

func (v *Version) SameFrom(s Version) bool {
	return v.Branch == s.Branch && v.Channel == s.Channel && v.Major == s.Major && v.Minor == s.Minor && v.Patch == s.Patch
}

func (s SemverTag) String() string {
	return s.Version.String() + " " + s.Hash().String()[0:7]
}
