package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"khosef/pkg/about"
	"khosef/pkg/core"
	"khosef/pkg/secrets"
	"os"
)

var (
	contextDir string
	verbose    bool
	rootCmd    = &cobra.Command{
		Use: "kh",
	}
)

func init() {
	pwd, _ := os.Getwd()
	rootCmd.AddCommand(about.NewVersionCommand())
	rootCmd.AddCommand(core.NewInitializerCommand(&contextDir))
	rootCmd.AddCommand(secrets.NewSecretsCommand(&contextDir, &verbose))

	rootCmd.PersistentFlags().StringVarP(&contextDir, "context", "C", pwd, "Context directory to run in")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
