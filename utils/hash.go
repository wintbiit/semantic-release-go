package utils

import "fmt"

// HashShort convert long hash to short hash.
func HashShort(s fmt.Stringer) string {
	return s.String()[:8]
}
