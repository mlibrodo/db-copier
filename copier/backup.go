package copier

import (
	"fmt"
	"github.com/mlibrodo/db-copier/aws/s3"
	"github.com/mlibrodo/db-copier/log"
	"github.com/mlibrodo/db-copier/postgres/conn"
	"github.com/mlibrodo/db-copier/postgres/pgcommands"
	"github.com/mlibrodo/db-copier/postgres/pgcommands/common"
	"os"
	"time"
)

type BackupToS3 struct {
	S3Backup s3.S3Object
}

func (in BackupToS3) Exec(db *conn.DBConnInfo) error {

	backupFile, err := backup(db, in.S3Backup.GetFileName())

	if err != nil {
		return err
	}

	s3Object := in.S3Backup
	if err = s3.Upload(s3Object, backupFile); err != nil {
		return err
	}

	log.WithFields(
		log.Fields{
			"Host":            db.DBHost,
			"Port":            db.DBPort,
			"DB":              db.DBName,
			"Backup File":     backupFile,
			"S3Backup Bucket": s3Object.Bucket,
			"S3Backup Key":    s3Object.Key,
		},
	).Debug("Backup success")

	return nil
}

func backup(db *conn.DBConnInfo, name string) (*string, error) {
	filePattern := fmt.Sprintf(`%v_%v-*.sql.tar.gz`, name, time.Now().Unix())
	tempFile, err := os.CreateTemp(common.TEMP_DIR, filePattern)

	if err != nil {
		return nil, err
	}

	dump := pgcommands.NewPGDump(db, tempFile.Name())
	dump.Verbose = true

	result := dump.Exec()

	if result.Error != nil {
		log.WithFields(
			log.Fields{
				"Command": result.FullCommand,
				"Error":   result.Error.Err,
			},
		).Debug("Backup failed")

		log.Error(result.Output)

		return nil, result.Error.Err
	}

	fullPath := tempFile.Name()

	log.WithFields(
		log.Fields{
			"PGDump Flags": result.FullCommand,
		},
	).Debug("Backup success")

	return &fullPath, nil
}
