package cmd

import (
	"github.com/hupe1980/gopherfy/pkg/fastcgi"
	"github.com/spf13/cobra"
)

type fastCGIOptions struct {
	addr string
	file string
	code string
}

func newFastCGICmd(globalOpts *globalOptions) *cobra.Command {
	opts := &fastCGIOptions{}

	cmd := &cobra.Command{
		Use:           "fastcgi",
		Short:         "Generate fastcgi gopher link",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			cgi := fastcgi.NewFastCGI(func(o *fastcgi.Options) {
				o.Addr = opts.addr
				o.File = opts.file
				o.Code = opts.code
			})

			payload := encodePayload(globalOpts.encoder, cgi.Payload())

			return output(payload, globalOpts.send)
		},
	}

	cmd.Flags().StringVarP(&opts.addr, "addr", "a", fastcgi.DefaultAddr, "fastcgi address")
	cmd.Flags().StringVarP(&opts.file, "file", "f", fastcgi.DefaultFile, "absolute php file path")
	cmd.Flags().StringVarP(&opts.code, "code", "c", fastcgi.DefaultCode, "code to execute")

	return cmd
}
