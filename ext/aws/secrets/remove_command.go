package secrets

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/spf13/cobra"
	"khosef/pkg/config"
)

type RemoveCommand struct {
	secretsClient *secretsmanager.Client
	contextPath   *string
	filePath      string
}

func NewRemoveCommand(profile *string, contextDir *string) *cobra.Command {
	return &cobra.Command{
		Use: "remove",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("missing required arguments")
			}

			client, err := getAwsSdk(*profile)
			if err != nil {
				return err
			}

			c := &RemoveCommand{
				secretsClient: client,
				contextPath:   contextDir,
				filePath:      args[0],
			}

			return c.Run()
		},
	}
}

func (r *RemoveCommand) Validate() error {
	return nil
}

func (r *RemoveCommand) Run() error {
	cnf, err := config.ReadConfig(*r.contextPath)
	if err != nil {
		return err
	}

	if i, existing := cnf.FindFile(r.filePath); existing != nil {
		cnf.Secrets = append(cnf.Secrets[:i], cnf.Secrets[i+1:]...)
		_, err = r.secretsClient.DeleteSecret(context.TODO(), &secretsmanager.DeleteSecretInput{
			SecretId: &existing.SecretId,
		})

		return cnf.Save()
	}

	return nil
}
