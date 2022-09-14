package secrets

import "github.com/spf13/cobra"

func NewSecretsCommand(contextDir *string, verbose *bool) *cobra.Command {
	c := &cobra.Command{
		Use: "secrets",
	}

	c.AddCommand(NewPullCommand(contextDir, verbose))

	return c
}
