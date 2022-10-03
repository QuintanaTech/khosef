package secrets

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/spf13/cobra"
	"khosef/pkg/cmd"
	cnf "khosef/pkg/config"
	"os"
	"path"
)

type PullCommand struct {
	secretsClient *secretsmanager.Client
	contextPath   *string
}

func NewPullCommand(profile *string, contextPath *string) *cobra.Command {
	c := &PullCommand{contextPath: contextPath}

	return &cobra.Command{
		Use: "pull",
		RunE: cmd.NewRunFnDecorator(func(cmd *cobra.Command, args []string) error {
			client, err := getAwsSdk(*profile)
			if err != nil {
				return err
			}

			c.secretsClient = client

			return nil
		}, cmd.NewSimpleRunFn(c)),
	}
}

func getAwsSdk(profile string) (*secretsmanager.Client, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithSharedConfigProfile(profile))
	if err != nil {
		return nil, err
	}

	return secretsmanager.NewFromConfig(cfg), nil
}

func (p *PullCommand) Validate() error {
	return nil
}

func (p *PullCommand) Run() error {
	c, err := cnf.ReadConfig(*p.contextPath)
	if err != nil {
		return err
	}

	for _, s := range c.GetSecretDefinitions() {
		if err = p.fetchToFile(s); err != nil {
			return err
		}

		fmt.Println(s.OutputPath, s.SecretId)
	}

	return nil
}

func (c *PullCommand) fetchToFile(d *cnf.SecretDefinition) error {
	s, err := c.secretsClient.GetSecretValue(context.TODO(), &secretsmanager.GetSecretValueInput{
		SecretId: &d.SecretId,
	})

	if err != nil {
		return err
	}

	p := path.Join(*c.contextPath, d.OutputPath)
	if _, err = os.Stat(path.Dir(p)); os.IsNotExist(err) {
		if err = os.MkdirAll(path.Dir(p), 0755); err != nil {
			return err
		}
	}

	if len(s.SecretBinary) > 0 {
		err = os.WriteFile(p, s.SecretBinary, 0644)
	} else {
		err = os.WriteFile(p, []byte(aws.ToString(s.SecretString)), 0644)
	}

	if err != nil {
		return err
	}

	return nil
}
