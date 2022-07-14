package copier

import (
	"github.com/mlibrodo/db-copier/aws/s3"
	"github.com/mlibrodo/db-copier/log"
	"github.com/mlibrodo/db-copier/postgres/conn"
	"github.com/mlibrodo/db-copier/postgres/pgcommands"
	"github.com/mlibrodo/db-copier/postgres/pgcommands/common"
	"os"
)

type RestoreFromS3 struct {
	S3 s3.S3Object
}

func (in RestoreFromS3) Exec(pgConnInfo *conn.DBConnInfo) error {
	var err error
	file, _ := os.CreateTemp(common.TEMP_DIR, "pg_dump-*.sql.tar.gz")

	defer func(file *os.File) {
		tmpErr := file.Close()

		if err != nil {
			err = tmpErr
		}

	}(file)

	if err = s3.Download(in.S3, file); err != nil {
		return err
	}

	// TODO check if we really want to delete the db
	if err = dropDB(pgConnInfo); err != nil {
		log.WithFields(
			log.Fields{
				"Database": pgConnInfo.DBName,
			},
		).Warning("Database does not exist, continuing to restore")
	}

	if err = createDB(pgConnInfo); err != nil {
		return err
	}

	return restoreFromFile(file, pgConnInfo)
}

func dropDB(pgConnInfo *conn.DBConnInfo) error {
	dropDB := pgcommands.NewDropDB(pgConnInfo)
	result := dropDB.Exec()

	if result.Error != nil {
		log.WithFields(
			log.Fields{
				"Command": result.FullCommand,
				"Error":   result.Error.Err,
			},
		).Error("DropDB failed")

		log.Error(result.Output)

		return result.Error.Err

	}
	log.WithFields(
		log.Fields{
			"Command": result.FullCommand,
		},
	).Debug("DropDB success")

	log.Debug(result.Output)

	return nil
}

func createDB(pgConnInfo *conn.DBConnInfo) error {
	createDB := pgcommands.NewCreateDB(pgConnInfo)
	result := createDB.Exec()

	if result.Error != nil {
		log.WithFields(
			log.Fields{
				"Command": result.FullCommand,
				"Error":   result.Error.Err,
			},
		).Error("CreateDB failed")

		log.Error(result.Output)

		return result.Error.Err

	}
	log.WithFields(
		log.Fields{
			"Command": result.FullCommand,
		},
	).Debug("CreateDB success")

	log.Debug(result.Output)

	return nil
}

func restoreFromFile(file *os.File, pgConnInfo *conn.DBConnInfo) error {

	pgRestore := pgcommands.NewPGRestore(pgConnInfo, file.Name())

	restoreExec := pgRestore.Exec()

	if restoreExec.Error != nil {
		log.WithFields(
			log.Fields{
				"Command": restoreExec.FullCommand,
				"Error":   restoreExec.Error.Err,
			},
		).Error("Restore failed")

		log.Error(restoreExec.Output)

		return restoreExec.Error.Err

	}
	log.WithFields(
		log.Fields{
			"Command": restoreExec.FullCommand,
		},
	).Debug("Restore success")

	log.Debug(restoreExec.Output)

	return nil
}
