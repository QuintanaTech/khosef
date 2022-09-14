package plugin

import (
	"fmt"
	cnf "khosef/pkg/config"
	"os/exec"
)

type PluginExecutor func(name string, args ...string) ([]byte, error)

type Plugin struct {
	config   *cnf.Config
	executor PluginExecutor
	silent   bool
}

func NewPlugin(config *cnf.Config) *Plugin {
	return &Plugin{
		config:   config,
		executor: shellExecutor,
		silent:   false,
	}
}

func shellExecutor(name string, args ...string) ([]byte, error) {
	return exec.Command(name, args...).Output()
}

func (p *Plugin) Execute(cmd string, args []string) error {
	a := append([]string{cmd}, "--context", p.config.GetContextDir())
	a = append(a, args...)
	b, err := p.executor(fmt.Sprintf("kh-%s", p.config.Provider), a...)
	if err != nil {
		return err
	}

	if !p.silent {
		fmt.Println(string(b))
	}

	return nil
}
