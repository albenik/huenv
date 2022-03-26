// Package version records versioning information about this module.
package version

import (
	"fmt"
)

const (
	Major      = 0
	Minor      = 1
	Patch      = 0
	PreRelease = ""
)

// String formats the version string for this module in semver format.
//
// Examples:
//	v1.20.1
//	v1.21.0-rc.1
func String() string {
	v := fmt.Sprintf("v%d.%d.%d", Major, Minor, Patch)
	if PreRelease != "" {
		v += "-" + PreRelease
	}
	return v
}
