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
	S3 s3.S3Object
}

func (in BackupToS3) Exec(db *conn.DBConnInfo) error {

	backupFile, err := backup(db, in.S3.GetFileName())

	if err != nil {
		return err
	}

	s3Object := in.S3
	if err = s3.Upload(s3Object, backupFile); err != nil {
		return err
	}

	log.WithFields(
		log.Fields{
			"Host":        db.DBHost,
			"Port":        db.DBPort,
			"DB":          db.DBName,
			"Backup File": backupFile,
			"S3 Bucket":   s3Object.Bucket,
			"S3 Key":      s3Object.Key,
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

	dumpExec := dump.Exec()

	if dumpExec.Error != nil {
		log.WithFields(
			log.Fields{
				"Command": dumpExec.FullCommand,
				"Error":   dumpExec.Error.Err,
			},
		).Debug("Backup failed")

		log.Error(dumpExec.Output)

		return nil, dumpExec.Error.Err

	}
	fullPath := tempFile.Name()

	log.WithFields(
		log.Fields{
			"PGDump Flags": dumpExec.FullCommand,
		},
	).Debug("Backup success")

	return &fullPath, nil
}
