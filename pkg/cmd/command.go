package cmd

import "github.com/spf13/cobra"

type Cmd interface {
	Run() error
	Validate() error
}

type RunFn func(cmd *cobra.Command, args []string) error

func NewRunFnDecorator(decorator RunFn, fn RunFn) RunFn {
	return func(cmd *cobra.Command, args []string) error {
		if err := decorator(cmd, args); err != nil {
			return err
		}

		return fn(cmd, args)
	}
}

func NewSimpleRunFn(c Cmd) RunFn {
	return func(cmd *cobra.Command, args []string) error {
		if err := c.Validate(); err != nil {
			return err
		}

		return c.Run()
	}
}
