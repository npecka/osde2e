package installselectors

import (
	"fmt"
	"log"

	"github.com/Masterminds/semver"
	"github.com/openshift/osde2e/pkg/common/config"
	"github.com/openshift/osde2e/pkg/common/spi"
	"github.com/openshift/osde2e/pkg/common/state"
)

func init() {
	registerSelector(previousReleaseFromDefault{})
}

// previousReleaseFromDefault is the selector which will
type previousReleaseFromDefault struct{}

func (p previousReleaseFromDefault) ShouldUse() bool {
	return config.Instance.Cluster.PreviousReleaseFromDefault > 0
}

func (p previousReleaseFromDefault) Priority() int {
	return 50
}

func (p previousReleaseFromDefault) SelectVersion(versionList *spi.VersionList) (*semver.Version, string, error) {
	availableVersions := versionList.AvailableVersions()
	defaultIndex := findDefaultVersionIndex(availableVersions)
	numReleasesFromDefault := config.Instance.Cluster.PreviousReleaseFromDefault
	versionType := fmt.Sprintf("version %d releases prior to the default", numReleasesFromDefault)

	if defaultIndex < 0 {
		log.Printf("unable to find default version in avaialable version list")
		state.Instance.Cluster.PreviousVersionFromDefaultFound = false
	}

	targetIndex := defaultIndex - numReleasesFromDefault

	if targetIndex < 0 {
		log.Printf("not enough enabled versions to go back %d releases", numReleasesFromDefault)
		state.Instance.Cluster.PreviousVersionFromDefaultFound = false
		return nil, versionType, nil
	}

	return availableVersions[targetIndex].Version(), versionType, nil
}

func findDefaultVersionIndex(versions []*spi.Version) int {
	for index, version := range versions {
		if version.Default() {
			return index
		}
	}

	return -1
}
