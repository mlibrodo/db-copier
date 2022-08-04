package copier

import (
	"github.com/mlibrodo/db-copier/postgres/conn"
)

type TrimAnonymizeAndBackup struct {
	RestoreFromS3 RestoreFromS3
	Trim          Trim
	Anonymize     Anonymize
	TrimmedBackup BackupToS3
	DropDB        bool
}

func (in TrimAnonymizeAndBackup) Exec(pgConnInfo *conn.DBConnInfo) error {

	var err error

	// restore a backup
	if err = in.RestoreFromS3.Exec(pgConnInfo); err != nil {
		return err
	}

	// trim the restore
	if err = in.Trim.Exec(pgConnInfo); err != nil {
		return err
	}

	// anonymize the restore
	if err = in.Anonymize.Exec(pgConnInfo); err != nil {
		return err
	}

	// backup the trimmed DB
	if err = in.TrimmedBackup.Exec(pgConnInfo); err != nil {
		return err
	}

	if in.DropDB {
		if err = dropDB(pgConnInfo); err != nil {
			return err
		}
	}

	return nil
}
