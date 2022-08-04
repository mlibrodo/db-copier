package cmd

import (
	"fmt"
	"github.com/mlibrodo/db-copier/aws/s3"
	"github.com/mlibrodo/db-copier/config"
	"github.com/mlibrodo/db-copier/copier"
	"github.com/mlibrodo/db-copier/postgres/conn"
	"github.com/spf13/cobra"
	"path/filepath"
)

// trimAndBackupCmd trims a backup in S3 and then backs that up
var trimAndBackupCmd = &cobra.Command{

	Use:   "trimAndBackup",
	Short: "Takes an S3 backup, restores it, trims that and then backs that backup to S3. Uses the connection to restore the DB",
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

		restore := copier.RestoreFromS3{
			S3: s3.S3Object{
				Bucket: config.GetConfig().Backup.S3Bucket,
				Key:    s3Key,
			},
		}

		path, file := filepath.Split(s3Key)
		backup := copier.BackupToS3{
			S3Backup: s3.S3Object{
				Bucket: config.GetConfig().Backup.S3Bucket,
				Key:    fmt.Sprintf("%strimmed-%s", path, file),
			},
		}
		// Load from a json file
		trim := copier.Trim{
			[]copier.TrimDef{{"foo", "i > 3"}},
		}

		// Load from a json file
		anonymize := copier.Anonymize{
			[]copier.AnonymizeDef{{"foobar", "x_str", copier.ANONYMIZE_STRING}},
		}

		// Load from a json file
		trimAndAnonymize := copier.TrimAnonymizeAndBackup{
			RestoreFromS3: restore,
			Trim:          trim,
			Anonymize:     anonymize,
			TrimmedBackup: backup,
			DropDB:        false,
		}

		trimAndAnonymize.Exec(connInfo)
	},
}

func init() {
	var s3Key string
	rootCmd.AddCommand(trimAndBackupCmd)
	trimAndBackupCmd.Flags().StringVarP(&s3Key, "s3Key", "k", "", "")
}
