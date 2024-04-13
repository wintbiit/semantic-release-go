package utils

import (
	"fmt"
	"strings"

	"github.com/wintbiit/semantic-release-go/types"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/storer"
)

const VersionFormat = "%s %s v%d.%d.%d"

func History(tags storer.ReferenceIter, season string, channel string) ([]types.SemverTag, error) {
	var semverTags []types.SemverTag
	var ref *plumbing.Reference
	var tagStr, currSeason, currChannel string
	var major, minor, patch int
	var semverTag types.SemverTag
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

		semverTag = types.SemverTag{
			Reference: ref,
			Version: types.Version{
				Branch:  currSeason,
				Channel: currChannel,
				Major:   major,
				Minor:   minor,
				Patch:   patch,
			},
		}

		semverTags = append(semverTags, semverTag)
	}

	return semverTags, nil
}

func SortTags(tags []types.SemverTag) {
	// bubble sort
	for i := 0; i < len(tags); i++ {
		for j := 0; j < len(tags)-i-1; j++ {
			if tags[j].Major > tags[j+1].Major {
				tags[j], tags[j+1] = tags[j+1], tags[j]
			} else if tags[j].Major == tags[j+1].Major {
				if tags[j].Minor > tags[j+1].Minor {
					tags[j], tags[j+1] = tags[j+1], tags[j]
				} else if tags[j].Minor == tags[j+1].Minor {
					if tags[j].Patch > tags[j+1].Patch {
						tags[j], tags[j+1] = tags[j+1], tags[j]
					}
				}
			}
		}
	}
}
