package core

import (
	"fmt"
	"github.com/spf13/cobra"
	"khosef/pkg/cmd"
)

type VersionCmd struct {
	version *Version
}

func NewVersionCommand() *cobra.Command {
	c := &VersionCmd{version: NewVersion()}

	return &cobra.Command{
		Use:  "version",
		RunE: cmd.NewSimpleRunFn(c),
	}
}

func (v *VersionCmd) Run() error {
	fmt.Println("Using kh version:", v.version.GetCurrent())

	release, _, err := getLatestRelease()
	if err != nil {
		return err
	}

	currentVer := NewVersion()
	if currentVer.IsNewer(*release.TagName) {
		fmt.Println("There is a newer version available:", *release.TagName)
	}

	return nil
}

func (v *VersionCmd) Validate() error {
	return nil
}
