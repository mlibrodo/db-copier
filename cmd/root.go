package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Short: "Manages copying, restoring, trimming data in a postgres database and storing it in S3Backup",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize()
	var dbHost string
	var dbPort int32
	var dbName string
	var dbUser string
	var dbPassword string

	rootCmd.PersistentFlags().StringVarP(&dbHost, "dbHost", "H", "", "")
	rootCmd.PersistentFlags().Int32VarP(&dbPort, "dbPort", "p", 0, "")
	rootCmd.PersistentFlags().StringVarP(&dbName, "dbName", "d", "", "")
	rootCmd.PersistentFlags().StringVarP(&dbUser, "dbUser", "u", "", "")
	rootCmd.PersistentFlags().StringVarP(&dbPassword, "dbPassword", "P", "", "")

}
