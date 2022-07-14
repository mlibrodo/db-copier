package pgcommands

import (
	"github.com/mlibrodo/db-copier/postgres/conn"
	"github.com/mlibrodo/db-copier/postgres/pgcommands/common"
)

var (
	// PGDropDBCmd is the path to the `dropdb` executable
	PGDropDBCmd = "dropdb"
)

type DropDB struct {
	*conn.DBConnInfo
}

func NewDropDB(pgConnInfo *conn.DBConnInfo) *DropDB {
	return &DropDB{DBConnInfo: pgConnInfo}
}

// Exec `Dropdb` for specified DB
func (x *DropDB) Exec() common.Result {
	execFn := common.PGCLIExecutor(PGDropDBCmd, x.DBConnInfo, x.ParseArgs)
	return execFn()
}

func (x *DropDB) ParseArgs() []string {

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
