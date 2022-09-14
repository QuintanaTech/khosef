package core

import (
	ver "github.com/hashicorp/go-version"
)

var version = "0.0+dev"

type Version struct {
	current string
}

func NewVersion() *Version {
	return &Version{current: version}
}

func (v *Version) GetCurrent() string {
	return v.current
}

func (v *Version) IsNewer(t string) bool {
	current, err := ver.NewVersion(v.current)
	if err != nil {
		panic(err)
	}

	test, err := ver.NewVersion(t)
	if err != nil {
		panic(err)
	}

	return test.GreaterThan(current)
}
