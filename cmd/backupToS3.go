package cmd

import (
	"github.com/mlibrodo/db-copier/aws/s3"
	"github.com/mlibrodo/db-copier/config"
	"github.com/mlibrodo/db-copier/copier"
	"github.com/mlibrodo/db-copier/postgres/conn"
	"github.com/spf13/cobra"
)

// backupCmd runs pgdump and stores the output to S3Backup
var backupCmd = &cobra.Command{

	Use:   "backup",
	Short: "Backup a DB to S3Backup bucket with the given s3 key",
	Run: func(cmd *cobra.Command, args []string) {
		var err error

		var dbHost string
		var dbPort int32
		var dbName string
		var dbUser string
		var dbPassword string
		var s3Key string

		parentFlags := cmd.Parent().PersistentFlags()
		flags := cmd.Flags()

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

		backup := copier.BackupToS3{
			S3Backup: s3.S3Object{
				Bucket: config.GetConfig().Backup.S3Bucket,
				Key:    s3Key,
			},
		}

		backup.Exec(connInfo)
	},
}

func init() {
	var s3Key string
	rootCmd.AddCommand(backupCmd)

	backupCmd.Flags().StringVarP(&s3Key, "s3Key", "k", "", "")
}
