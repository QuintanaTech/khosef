package core

import (
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"khosef/pkg/cmd"
	"khosef/pkg/config"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type InitializerCmd struct {
	contextDir *string
}

func NewInitializerCommand(contextDir *string) *cobra.Command {
	c := &InitializerCmd{
		contextDir: contextDir,
	}

	return &cobra.Command{
		Use:   "init",
		Short: "Initialize a new project",
		RunE:  cmd.NewSimpleRunFn(c),
	}
}

func (i *InitializerCmd) Run() error {
	_, err := config.ReadConfig(*i.contextDir)
	if err == nil {
		fmt.Println("Project is already initialized")

		return nil
	}

	if _, err := os.Stat(*i.contextDir); os.IsNotExist(err) {
		fmt.Println("Context", *i.contextDir, "does not exist")

		return nil
	}

	if err := initializeConfig(*i.contextDir); err != nil {
		return err
	}

	return initializeGitignore(*i.contextDir)
}

func (i *InitializerCmd) Validate() error {
	return nil
}

func initializeConfig(contextDir string) error {
	c := config.NewConfig(contextDir)
	c.Provider = "aws"
	c.Secrets = []string{"example.secret:example/secret"}

	return c.Save()
}

func initializeGitignore(contextDir string) error {
	fpath, err := filepath.Abs(path.Join(contextDir, ".gitignore"))
	if err != nil {
		return err
	}

	lines := []string{
		"##### Khosef Secrets Ignores #####",
		"*.secret",
		"*.secret.yaml",
		"*.secret.json",
	}

	f, err := os.OpenFile(fpath, os.O_RDONLY, 0644)
	if err != nil {
		if os.IsNotExist(err) {
			return updateGitIgnore(fpath, lines)
		}

		return err
	}

	scanner := bufio.NewScanner(f)
	var data []string
	var line string
	var pos int

	for scanner.Scan() {
		bs := scanner.Bytes()
		data = append(data, string(bs))

		if strings.HasPrefix(strings.ToLower(string(bs)), strings.ToLower(lines[0])) {
			pos = len(data) - 1
		}
	}

	if err = f.Close(); err != nil {
		return err
	}

	data = append(data, "")
	for _, line = range lines {
		if pos < len(data) {
			data[pos] = line
		} else {
			data = append(data, line)
		}

		pos += 1
	}

	return updateGitIgnore(fpath, data)
}

func updateGitIgnore(fpath string, data []string) error {
	f, err := os.OpenFile(fpath, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	writer := bufio.NewWriter(f)
	defer f.Close()

	for _, line := range data {
		if _, err := writer.WriteString(line + "\n"); err != nil {
			return err
		}
	}

	return writer.Flush()
}
