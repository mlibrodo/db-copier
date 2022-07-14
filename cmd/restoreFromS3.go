package cmd

import (
	"github.com/mlibrodo/db-copier/aws/s3"
	"github.com/mlibrodo/db-copier/config"
	"github.com/mlibrodo/db-copier/copier"
	"github.com/mlibrodo/db-copier/postgres/conn"
	"github.com/spf13/cobra"
)

// restoreCmd runs pg_restore using a file from S3
var restoreCmd = &cobra.Command{

	Use:   "restore",
	Short: "Restore a DB from an S3 bucket with the given s3 key",
	Run: func(cmd *cobra.Command, args []string) {
		var err error

		var dbHost string
		var dbPort int32
		var dbName string
		var dbUser string
		var dbPassword string
		var s3Key string

		flags := cmd.Flags()
		if dbHost, err = flags.GetString("dbHost"); err != nil {
			panic("dbHost not set")
		}
		if dbPort, err = flags.GetInt32("dbPort"); err != nil {
			panic("dbPort not set")
		}
		if dbName, err = flags.GetString("dbName"); err != nil {
			panic("dbName not set")
		}
		if dbUser, err = flags.GetString("dbUser"); err != nil {
			panic("dbUser not set")
		}
		if dbPassword, err = flags.GetString("dbPassword"); err != nil {
			panic("dbPassword not set")
		}
		if s3Key, err = flags.GetString("s3Key"); err != nil {
			panic("s3Key not set")
		}

		connInfo := &conn.DBConnInfo{
			DBHost:   dbHost,
			DBPort:   dbPort,
			DBName:   dbName,
			Username: dbUser,
			Password: dbPassword,
		}

		restore := copier.RestoreFromS3{
			S3: s3.S3Object{
				Bucket: config.GetConfig().Backup.S3Bucket,
				Key:    s3Key,
			},
		}

		restore.Exec(connInfo)
	},
}

func init() {
	var s3Key string

	rootCmd.AddCommand(restoreCmd)
	restoreCmd.Flags().StringVarP(&s3Key, "s3Key", "k", "", "")
}
