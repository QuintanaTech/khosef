package plugin

import (
	"fmt"
	"khosef/pkg/config"
	"testing"
)

func TestPlugin_Execute(t *testing.T) {
	cnf := config.NewConfig("/tmp/test")
	cnf.Provider = "foo"

	p := NewPlugin(cnf)
	p.executor = func(name string, args ...string) ([]byte, error) {
		return []byte(fmt.Sprintf("%s %s", name, args)), nil
	}
	p.silent = true
	args := []string{"arg1", "--opt1", "opt"}
	if err := p.Execute("test", args); err != nil {
		t.Error("should run successfully")
	}
}
