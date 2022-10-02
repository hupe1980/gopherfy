package cmd

import (
	"fmt"

	"github.com/hupe1980/gopherfy/pkg/mysql"
	"github.com/spf13/cobra"
)

type mySQLOptions struct {
	addr  string
	user  string
	query string
}

func NewMySQLCmd() *cobra.Command {
	opts := &mySQLOptions{}

	cmd := &cobra.Command{
		Use:           "mysql",
		Short:         "Genrate mysql gopher link",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			my := mysql.NewMySQL(func(o *mysql.Options) {
				o.Addr = opts.addr
				o.User = opts.user
				o.Query = opts.query
			})

			fmt.Println(my.Payload())

			return nil
		},
	}

	cmd.Flags().StringVarP(&opts.addr, "addr", "a", mysql.DefaultAddr, "mysql address")
	cmd.Flags().StringVarP(&opts.user, "user", "u", mysql.DefaultUser, "mysql username")
	cmd.Flags().StringVarP(&opts.query, "query", "q", "", "mysql query")

	return cmd
}
