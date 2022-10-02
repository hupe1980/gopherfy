package cmd

import (
	"fmt"

	"github.com/hupe1980/gopherfy/pkg/http"
	"github.com/spf13/cobra"
)

type httpOptions struct {
	addr    string
	method  string
	version string
}

func NewHTTPCmd() *cobra.Command {
	opts := &httpOptions{}

	cmd := &cobra.Command{
		Use:           "http",
		Short:         "Genrate http gopher link",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			http := http.NewHTTP(func(o *http.Options) {
				o.Addr = opts.addr
			})

			fmt.Println(http.Payload())

			return nil
		},
	}

	cmd.Flags().StringVarP(&opts.addr, "addr", "a", http.DefaultAddr, "http address")
	cmd.Flags().StringVarP(&opts.method, "method", "X", http.DefaultMethod, "http method")
	cmd.Flags().StringVarP(&opts.version, "version", "V", http.DefaultMethod, "http version")

	return cmd
}
