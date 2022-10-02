package cmd

import (
	"fmt"

	"github.com/hupe1980/gopherfy/pkg/postgresql"
	"github.com/spf13/cobra"
)

type postgreSQLOptions struct {
	addr  string
	user  string
	db    string
	query string
}

func NewPostgreSQLCmd() *cobra.Command {
	opts := &postgreSQLOptions{}

	cmd := &cobra.Command{
		Use:           "postgresql",
		Short:         "Genrate postgresql gopher link",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			p := postgresql.NewPostgreSQL(func(o *postgresql.Options) {
				o.Addr = opts.addr
				o.User = opts.user
				o.DB = opts.db
				o.Query = opts.query
			})

			fmt.Println(p.Payload())

			return nil
		},
	}

	cmd.Flags().StringVarP(&opts.addr, "addr", "a", postgresql.DefaultAddr, "postgresql address")
	cmd.Flags().StringVarP(&opts.user, "user", "u", postgresql.DefaultUser, "postgresql username")
	cmd.Flags().StringVarP(&opts.db, "db", "d", "", "postgresql database")
	cmd.Flags().StringVarP(&opts.query, "query", "q", "", "postgresql query")

	return cmd
}
