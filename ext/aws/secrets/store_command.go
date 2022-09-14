package secrets

import (
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/spf13/cobra"
	"khosef/pkg/cmd"
	"khosef/pkg/config"
)

type StoreCommmand struct {
	secretsClient *secretsmanager.Client
	contextPath   *string
}

func NewStoreCommand(profile *string, contextPath *string) *cobra.Command {
	c := &StoreCommmand{contextPath: contextPath}

	return &cobra.Command{
		Use: "store",
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

func (s *StoreCommmand) Validate() error {
	return nil
}

func (s *StoreCommmand) Run() error {
	cnf, err := config.ReadConfig(*s.contextPath)
	if err != nil {
		return err
	}

	for _, sec := range cnf.GetSecretDefinitions() {
		if err = saveSecret(s.secretsClient, sec, *s.contextPath); err != nil {
			return err
		}
	}

	return nil
}
