package cmd

import (
	"github.com/mlibrodo/db-copier/copier"
	"github.com/mlibrodo/db-copier/postgres/conn"
	"github.com/spf13/cobra"
)

// trimCmd trims a DB with a trim definition file
var trimCmd = &cobra.Command{

	Use:   "trim",
	Short: "trim a DB",
	Run: func(cmd *cobra.Command, args []string) {
		var err error

		var dbHost string
		var dbPort int32
		var dbName string
		var dbUser string
		var dbPassword string

		parentFlags := cmd.Parent().PersistentFlags()

		if dbHost, err = parentFlags.GetString("dbHost"); err != nil {
			panic("dbHost not set")
		}
		if dbPort, err = parentFlags.GetInt32("dbPort"); err != nil {
			panic("dbPort not set")
		}
		if dbName, err = parentFlags.GetString("dbName"); err != nil {
			panic("dbName not set")
		}
		if dbUser, err = parentFlags.GetString("dbUser"); err != nil {
			panic("dbUser not set")
		}
		if dbPassword, err = parentFlags.GetString("dbPassword"); err != nil {
			panic("dbPassword not set")
		}

		connInfo := &conn.DBConnInfo{
			DBHost:   dbHost,
			DBPort:   dbPort,
			DBName:   dbName,
			Username: dbUser,
			Password: dbPassword,
		}

		// Load from a json?
		trim := copier.TrimDefList{
			[]copier.TrimDef{{"foo", "i > 1"}},
		}

		trim.Exec(connInfo)
	},
}

func init() {
	rootCmd.AddCommand(trimCmd)
}
