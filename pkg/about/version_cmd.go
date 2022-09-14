package about

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

	return nil
}

func (v *VersionCmd) Validate() error {
	return nil
}
