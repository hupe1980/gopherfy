package cmd

import (
	"fmt"

	"github.com/hupe1980/gopherfy/pkg/postgres"
	"github.com/spf13/cobra"
)

type postgresOptions struct {
	addr  string
	user  string
	db    string
	query string
}

func newPostgresCmd(encoder *string) *cobra.Command {
	opts := &postgresOptions{}

	cmd := &cobra.Command{
		Use:           "postgres",
		Short:         "Generate postgres gopher link",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			postgres := postgres.NewPostgres(func(o *postgres.Options) {
				o.Addr = opts.addr
				o.User = opts.user
				o.DB = opts.db
				o.Query = opts.query
			})

			payload := encodePayload(*encoder, postgres.Payload())

			fmt.Println(payload)

			return nil
		},
	}

	cmd.Flags().StringVarP(&opts.addr, "addr", "a", postgres.DefaultAddr, "postgres address")
	cmd.Flags().StringVarP(&opts.user, "user", "u", postgres.DefaultUser, "postgres username")
	cmd.Flags().StringVarP(&opts.db, "db", "d", "", "postgres database name")
	cmd.Flags().StringVarP(&opts.query, "query", "q", "", "postgres query")

	return cmd
}
