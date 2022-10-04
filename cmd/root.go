package cmd

import (
	"fmt"
	"os"

	"github.com/hupe1980/gopherfy/internal"
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
	var (
		encoder string
	)

	cmd := &cobra.Command{
		Use:           "gopherfy",
		Version:       version,
		Short:         "Tool to generate gopher links for exploiting SSRF",
		SilenceErrors: true,
	}

	cmd.PersistentFlags().StringVarP(&encoder, "encoder", "e", "none", `the encoder to use. allowed: "base64", "url" or "none"`)

	cmd.AddCommand(
		NewHTTPCmd(&encoder),
		NewMySQLCmd(&encoder),
		NewSMTPCmd(&encoder),
	)

	return cmd
}

func encodePayload(encoder, payload string) string {
	switch encoder {
	case "base64":
		return internal.Base64UrlSafeEncode(payload)
	case "url":
		return internal.URLEncode(payload)
	default:
		return payload
	}
}
