package secrets

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager/types"
	"github.com/spf13/cobra"
	"io/ioutil"
	"khosef/pkg/config"
	"path"
)

type AddCommand struct {
	secretsClient *secretsmanager.Client
	contextPath   *string
	filePath      string
	secretId      string
}

func NewAddCommand(profile *string, contextDir *string) *cobra.Command {
	return &cobra.Command{
		Use: "add [file path] [secret id]",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 2 {
				return errors.New("missing require arguments")
			}

			client, err := getAwsSdk(*profile)
			if err != nil {
				return err
			}

			c := &AddCommand{
				secretsClient: client,
				contextPath:   contextDir,
				filePath:      args[0],
				secretId:      args[1],
			}

			return c.Run()
		},
	}
}

func (a *AddCommand) Validate() error {
	return nil
}

func (a *AddCommand) Run() error {
	cnf, err := config.ReadConfig(*a.contextPath)
	if err != nil {
		return err
	}

	if _, existing := cnf.FindFile(a.filePath); existing == nil {
		cnf.Secrets = append(cnf.Secrets, fmt.Sprintf("%s:%s", a.filePath, a.secretId))
	}

	err = saveSecret(a.secretsClient, &config.SecretDefinition{
		SecretId:   a.secretId,
		OutputPath: a.filePath,
	}, *a.contextPath)
	if err != nil {
		return err
	}

	return cnf.Save()
}

func saveSecret(client *secretsmanager.Client, d *config.SecretDefinition, contextPath string) error {
	_, err := client.DescribeSecret(context.TODO(), &secretsmanager.DescribeSecretInput{
		SecretId: &d.SecretId,
	})

	if err != nil {
		var nfe *types.ResourceNotFoundException
		if errors.As(err, &nfe) {
			_, err = client.CreateSecret(context.TODO(), &secretsmanager.CreateSecretInput{
				Name: &d.SecretId,
			})

			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	b, err := ioutil.ReadFile(path.Join(contextPath, d.OutputPath))
	if err != nil {
		return err
	}

	_, err = client.PutSecretValue(context.TODO(), &secretsmanager.PutSecretValueInput{
		SecretId:     &d.SecretId,
		SecretBinary: b,
	})

	return err
}
