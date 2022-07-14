package pgcommands

import (
	"github.com/mlibrodo/db-copier/postgres/conn"
	"github.com/mlibrodo/db-copier/postgres/pgcommands/common"
)

var (
	// PGCreateDBCmd is the path to the `createdb` executable
	PGCreateDBCmd = "createdb"
)

type CreateDB struct {
	*conn.DBConnInfo
}

func NewCreateDB(pgConnInfo *conn.DBConnInfo) *CreateDB {
	return &CreateDB{DBConnInfo: pgConnInfo}
}

// Exec `createdb` for specified DB
func (x *CreateDB) Exec() common.Result {
	execFn := common.PGCLIExecutor(PGCreateDBCmd, x.DBConnInfo, x.ParseArgs)
	return execFn()
}

func (x *CreateDB) ParseArgs() []string {

	var args []string

	if y := x.DBHostAsCmdArg(); y != nil {
		args = append(args, *y)
	}

	if y := x.DBPortAsCmdArg(); y != nil {
		args = append(args, *y)
	}

	if y := x.UsernameAsCmdArg(); y != nil {
		args = append(args, *y)
	}

	if y := x.DBNameAsCmdArg(nil); y != nil {
		args = append(args, *y)
	}

	return args
}
