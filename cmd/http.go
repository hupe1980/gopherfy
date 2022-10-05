package cmd

import (
	"fmt"

	"github.com/hupe1980/gopherfy/pkg/http"
	"github.com/spf13/cobra"
)

type httpOptions struct {
	addr      string
	method    string
	version   string
	path      string
	userAgent string
	headers   map[string]string
}

func newHTTPCmd(encoder *string) *cobra.Command {
	opts := &httpOptions{}

	cmd := &cobra.Command{
		Use:           "http",
		Short:         "Generate http gopher link",
		Example:       `gopherfy http -a 169.254.169.254:80 -p /latest/api/token -X PUT -H X-aws-ec2-metadata-token-ttl-seconds=21600`,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			http := http.NewHTTP(func(o *http.Options) {
				o.Addr = opts.addr
				o.Method = opts.method
				o.Version = opts.version
				o.Path = opts.path
				o.UserAgent = opts.userAgent
				o.ExtraHeaders = opts.headers
			})

			payload := encodePayload(*encoder, http.Payload())

			fmt.Println(payload)

			return nil
		},
	}

	// -H --header
	// -d --data
	// -c --cookie
	cmd.Flags().StringVarP(&opts.addr, "addr", "a", http.DefaultAddr, "http address")
	cmd.Flags().StringVarP(&opts.method, "request", "X", http.DefaultMethod, "http request method")
	cmd.Flags().StringVarP(&opts.version, "version", "V", http.DefaultVersion, "http protocol version")
	cmd.Flags().StringVarP(&opts.path, "path", "p", http.DefaultPath, "http path")
	cmd.Flags().StringVarP(&opts.userAgent, "user-agent", "A", http.DefaultUserAgent, "http user agent")
	cmd.Flags().StringToStringVarP(&opts.headers, "header", "H", nil, "http header value (key=value)")

	return cmd
}
