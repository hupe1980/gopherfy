package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/hupe1980/gopherfy/internal"
	"github.com/hupe1980/gopherfy/pkg/gopher"
	"github.com/spf13/cobra"
)

func Execute(version string) {
	rootCmd := newRootCmd(version)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

type globalOptions struct {
	encoder string
	send    bool
}

func newRootCmd(version string) *cobra.Command {
	globalOpts := &globalOptions{}

	cmd := &cobra.Command{
		Use:           "gopherfy",
		Version:       version,
		Short:         "Tool to generate gopher links for exploiting SSRF",
		SilenceErrors: true,
	}

	cmd.PersistentFlags().StringVarP(&globalOpts.encoder, "encoder", "e", "none", `the encoder to use. allowed: "base64", "url" or "none"`)
	cmd.PersistentFlags().BoolVarP(&globalOpts.send, "send", "", false, "send the selector string")

	cmd.AddCommand(
		newFastCGICmd(globalOpts),
		newHTTPCmd(globalOpts),
		newMySQLCmd(globalOpts),
		newSMTPCmd(globalOpts),
		newPostgresCmd(globalOpts),
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

func send(payload string) error {
	resp, err := gopher.Get(payload)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Println(string(body))

	return nil
}

func output(payload string, sending bool) error {
	if sending {
		return send(payload)
	}

	fmt.Println(payload)

	return nil
}
