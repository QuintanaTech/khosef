package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"khosef/ext/aws/secrets"
	"os"
)

var (
	profile    string
	contextDir string
	rootCmd    = &cobra.Command{
		Use: "kh-aws",
	}
)

func init() {
	pwd, _ := os.Getwd()

	rootCmd.AddCommand(secrets.NewPullCommand(&profile, &contextDir))
	rootCmd.AddCommand(secrets.NewStoreCommand(&profile, &contextDir))
	rootCmd.AddCommand(secrets.NewRemoveCommand(&profile, &contextDir))
	rootCmd.AddCommand(secrets.NewAddCommand(&profile, &contextDir))

	rootCmd.PersistentFlags().StringVarP(&profile, "profile", "P", "", "Config profile")
	rootCmd.PersistentFlags().StringVarP(&contextDir, "context", "C", pwd, "Context directory")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
