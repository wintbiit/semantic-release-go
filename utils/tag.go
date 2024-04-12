package utils

import (
	"fmt"
	"strings"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/storer"
)

type Version struct {
	season  string
	channel string
	major   int
	minor   int
	patch   int
}

type SemverTag struct {
	*plumbing.Reference
	Version
}

const VersionFormat = "%s %s v%d.%d.%d"

func ValidTags(tags storer.ReferenceIter, season string, channel string) ([]SemverTag, error) {
	var semverTags []SemverTag
	var ref *plumbing.Reference
	var tagStr, currSeason, currChannel string
	var major, minor, patch int
	var semverTag SemverTag
	var err error

	for {
		ref, err = tags.Next()
		if err != nil || !ref.Name().IsTag() {
			break
		}

		tagStr = ref.Name().Short()
		tagStr = strings.ReplaceAll(tagStr, "/", " ")
		// check if fits the format
		if _, err = fmt.Sscanf(tagStr, VersionFormat, &currSeason, &currChannel, &major, &minor, &patch); err != nil {
			continue
		}

		// check if the season and channel are the same
		if currSeason != season || currChannel != channel {
			continue
		}

		semverTag = SemverTag{
			Reference: ref,
			Version: Version{
				season:  currSeason,
				channel: currChannel,
				major:   major,
				minor:   minor,
				patch:   patch,
			},
		}

		semverTags = append(semverTags, semverTag)
	}

	return semverTags, nil
}

func SortTags(tags []SemverTag) {
	// bubble sort
	for i := 0; i < len(tags); i++ {
		for j := 0; j < len(tags)-i-1; j++ {
			if tags[j].major > tags[j+1].major {
				tags[j], tags[j+1] = tags[j+1], tags[j]
			} else if tags[j].major == tags[j+1].major {
				if tags[j].minor > tags[j+1].minor {
					tags[j], tags[j+1] = tags[j+1], tags[j]
				} else if tags[j].minor == tags[j+1].minor {
					if tags[j].patch > tags[j+1].patch {
						tags[j], tags[j+1] = tags[j+1], tags[j]
					}
				}
			}
		}
	}
}
