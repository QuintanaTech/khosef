package about

var version = "unknown"

type Version struct {
	current string
}

func NewVersion() *Version {
	return &Version{current: version}
}

func (v *Version) GetCurrent() string {
	return v.current
}
