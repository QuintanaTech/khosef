package plugin

import (
	"github.com/spf13/cobra"
	"khosef/pkg/cmd"
	"khosef/pkg/config"
)

type PluginCommand struct {
	name       string
	contextDir *string
	verbose    *bool
	args       []string
}

func PluginRunFn(p *PluginCommand) cmd.RunFn {
	return cmd.NewRunFnDecorator(func(cmd *cobra.Command, args []string) error {
		p.args = args

		return nil
	}, cmd.NewSimpleRunFn(p))
}

func NewPluginCommand(name string, short string, contextDir *string, verbose *bool) *cobra.Command {
	p := &PluginCommand{
		name:       name,
		contextDir: contextDir,
		verbose:    verbose,
	}

	return &cobra.Command{
		Use:   p.name,
		Short: short,
		RunE:  PluginRunFn(p),
	}
}

func (p *PluginCommand) Validate() error {
	return nil
}

func (p *PluginCommand) Run() error {
	c, err := config.ReadConfig(*p.contextDir)
	if err != nil {
		return err
	}

	i := NewPlugin(c)

	return i.Execute(p.name, p.args)
}
