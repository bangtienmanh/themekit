package kit

import (
	"crypto"
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/inconshreveable/go-update"
)

var (
	ThemeKitVersion = version{Major: 0, Minor: 4, Patch: 7}
)

type versionComparisonResult int

const (
	// VersionLessThan is a versionComparisonResult for a version that is less than
	VersionLessThan versionComparisonResult = -1
	// VersionEqual is a versionComparisonResult for a version that is equal
	VersionEqual = 0
	// VersionEqual is a versionComparisonResult for a version that is greater than
	VersionGreaterThan = 1
)

// LibraryInfo will return a string array with information about the library used
// for logging.
func LibraryInfo() []string {
	return []string{
		"ThemeKit - Shopify Theme Utilities",
		ThemeKitVersion.String(),
	}
}

type version struct {
	Major int
	Minor int
	Patch int
}

func (v version) String() string {
	return fmt.Sprintf("v%d.%d.%d", v.Major, v.Minor, v.Patch)
}

func (v version) toArray() [3]int {
	return [3]int{v.Major, v.Minor, v.Patch}
}

// Compare ... I often get confused by comparison, so comparison results are going
// to be the same as what <=> would return in Ruby.
// http://ruby-doc.org/core-1.9.3/Comparable.html
func (v version) Compare(o version) versionComparisonResult {
	vAry := v.toArray()
	oAry := o.toArray()
	for i := 0; i < len(vAry); i++ {
		diff := vAry[i] - oAry[i]
		if diff < 0 {
			return VersionLessThan
		} else if diff > 0 {
			return VersionGreaterThan
		}
	}
	return VersionEqual
}

// ParseVersionString will parse a string and convert it into a version for comparison.
func ParseVersionString(ver string) version {
	sanitizedVer := strings.Replace(ver, "v", "", 1)
	expandedVersionString := strings.Split(sanitizedVer, ".")
	major, _ := strconv.Atoi(expandedVersionString[0])
	minor, _ := strconv.Atoi(expandedVersionString[1])
	patch, _ := strconv.Atoi(expandedVersionString[2])
	return version{Major: major, Minor: minor, Patch: patch}
}

// ApplyUpdate will take a url and digest and download the update specified, then
// apply it to the current version.
func ApplyUpdate(updateURL, digest string) error {
	checksum, err := hex.DecodeString(digest)
	if err != nil {
		return err
	}

	updateFile, err := http.Get(updateURL)
	if err != nil {
		return err
	}
	defer updateFile.Body.Close()

	err = update.Apply(updateFile.Body, update.Options{
		Hash:     crypto.MD5,
		Checksum: checksum,
	})
	if err != nil {
		if rerr := update.RollbackError(err); rerr != nil {
			fmt.Printf("Failed to rollback from bad update: %v", rerr)
		}
	}
	return err
}
