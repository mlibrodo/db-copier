package pgcommands

import (
	"github.com/mlibrodo/db-copier/postgres/conn"
	"github.com/mlibrodo/db-copier/postgres/pgcommands/common"
)

var (
	// PGDumpCmd is the path to the `pg_dump` executable
	PGIsReadyCmd = "pg_isready"
)

type IsReady struct {
	*conn.PGConnInfo
	// Seconds to wait when attempting connection. 0 disables (default: 3)
	Timeout int
}

func (x *IsReady) Exec() common.Result {
	execFn := common.PGCLIExecutor(PGIsReadyCmd, x.PGConnInfo, x.ParseArgs)

	return execFn()

}

func (x *IsReady) ParseArgs() []string {
	var args []string

	if y := x.DBHostAsCmdArg(); y != nil {
		args = append(args, *y)
	}

	if y := x.DBPortAsCmdArg(); y != nil {
		args = append(args, *y)
	}

	return args
}
