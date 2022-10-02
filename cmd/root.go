package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func Execute(version string) {
	rootCmd := newRootCmd(version)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func newRootCmd(version string) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "gopherfy",
		Version:       version,
		Short:         "Tool to generate gopher links for exploiting SSRF",
		SilenceErrors: true,
	}

	cmd.AddCommand(
		NewHTTPCmd(),
		NewMySQLCmd(),
		NewPostgreSQLCmd(),
		NewSMTPCmd(),
	)

	return cmd
}
