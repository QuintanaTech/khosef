package secrets

import (
	"github.com/spf13/cobra"
	"khosef/pkg/plugin"
)

func NewSecretsCommand(contextDir *string, verbose *bool) *cobra.Command {
	c := &cobra.Command{
		Use: "secrets",
	}

	c.AddCommand(plugin.NewPluginCommand("add", "Add and persist secret to manifest", contextDir, verbose))
	c.AddCommand(plugin.NewPluginCommand("pull", "Pull secrets defined in manifest", contextDir, verbose))
	c.AddCommand(plugin.NewPluginCommand("remove", "Remove secret", contextDir, verbose))
	c.AddCommand(plugin.NewPluginCommand("store", "Persist all secrets to storage", contextDir, verbose))

	return c
}
