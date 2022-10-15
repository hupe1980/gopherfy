package cmd

import (
	"github.com/hupe1980/gopherfy/pkg/smtp"
	"github.com/spf13/cobra"
)

type smtpOptions struct {
	addr    string
	server  string
	from    string
	to      string
	subject string
	msg     string
}

func newSMTPCmd(globalOpts *globalOptions) *cobra.Command {
	opts := &smtpOptions{}

	cmd := &cobra.Command{
		Use:           "smtp",
		Short:         "Generate smtp gopher link",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			smtp := smtp.NewSMTP(func(o *smtp.Options) {
				o.Addr = opts.addr
			})

			payload := encodePayload(globalOpts.encoder, smtp.Payload())

			return output(payload, globalOpts.send)
		},
	}

	cmd.Flags().StringVarP(&opts.addr, "addr", "a", smtp.DefaultAddr, "smtp address")
	cmd.Flags().StringVar(&opts.server, "server", smtp.DefaultServer, "smtp server")
	cmd.Flags().StringVarP(&opts.from, "from", "f", "", "smtp mail from")
	cmd.Flags().StringVarP(&opts.to, "to", "t", "", "smtp mail to")
	cmd.Flags().StringVarP(&opts.subject, "subject", "s", "", "smtp subject")
	cmd.Flags().StringVarP(&opts.msg, "message", "m", "", "smtp message")

	return cmd
}
