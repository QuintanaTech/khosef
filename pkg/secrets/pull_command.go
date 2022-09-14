package secrets

import (
	"fmt"
	"github.com/spf13/cobra"
	"khosef/pkg/cmd"
	"khosef/pkg/config"
	"os/exec"
)

type PullCommand struct {
	contextDir *string
	verbose    *bool
	args       []string
}

func NewPullCommand(contextDir *string, verbose *bool) *cobra.Command {
	c := &PullCommand{contextDir: contextDir, verbose: verbose}

	return &cobra.Command{
		Use:   "pull",
		Short: "Pull secrets defined in manifest",
		RunE: cmd.NewRunFnDecorator(func(cmd *cobra.Command, args []string) error {
			c.args = args

			return nil
		}, cmd.NewSimpleRunFn(c)),
	}
}

func (p *PullCommand) Validate() error {
	return nil
}

func (p *PullCommand) Run() error {
	c, err := config.ReadConfig(*p.contextDir)
	if err != nil {
		return err
	}

	args := append([]string{"pull"}, "--context", *p.contextDir)
	args = append(args, p.args...)
	b, err := exec.Command(fmt.Sprintf("kh-%s", c.Provider), args...).Output()
	if err != nil {
		return err
	}

	fmt.Println(string(b))

	return nil
}
