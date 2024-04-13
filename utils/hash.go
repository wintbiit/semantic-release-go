package utils

import (
	"fmt"

	"github.com/wintbiit/semantic-release-go/types"
)

// HashShort convert long hash to short hash.
func HashShort(s fmt.Stringer) string {
	return s.String()[:8]
}

func ChannelValid(channel string) bool {
	switch channel {
	case types.ChannelInsider, types.ChannelAlpha, types.ChannelBeta, types.ChannelRelease:
		return true
	default:
		return false
	}
}
