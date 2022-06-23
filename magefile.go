//go:build mage
// +build mage

package main

import (
	"os"

	"get.porter.sh/magefiles/mixins"
	"get.porter.sh/magefiles/releases"
)

const (
	mixinName    = "xfer"
	mixinPackage = "github.com/richbai90/xfer/mixin/xfer"
	mixinBin     = "bin/mixins/" + mixinName
)

var magefile = mixins.NewMagefile(mixinPackage, mixinName, mixinBin)

// Build the mixin
func Build() {
	magefile.Build()
}

// Cross-compile the mixin before a release
func XBuildAll() {
	magefile.XBuildAll()
}

// Run unit tests
func TestUnit() {
	magefile.TestUnit()
}

func Test() {
	magefile.Test()
}

// Publish the mixin to github
func Publish() {
	// You can test out publishing locally by overriding PORTER_RELEASE_REPOSITORY and PORTER_PACKAGES_REMOTE
	if _, overridden := os.LookupEnv(releases.ReleaseRepository); !overridden {
		os.Setenv(releases.ReleaseRepository, "github.com/richbai90/xfer")
	}
	magefile.PublishBinaries()

	// TODO: uncomment out the lines below to publish a mixin feed
	// Set PORTER_PACKAGES_REMOTE to a repository that will contain your mixin feed, similar to github.com/getporter/packages
	//if _, overridden := os.LookupEnv(releases.PackagesRemote); !overridden {
	//	os.Setenv("PORTER_PACKAGES_REMOTE", "git@github.com:YOURNAME/YOUR_PACKAGES_REPOSITORY")
	//}
	//magefile.PublishMixinFeed()
}

// Install the mixin
func Install() {
	magefile.Install()
}

// Remove generated build files
func Clean() {
	magefile.Clean()
}
