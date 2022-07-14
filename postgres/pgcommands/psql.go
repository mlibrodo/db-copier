package pgcommands

import (
	"fmt"
	"github.com/mlibrodo/db-copier/postgres/conn"
	"github.com/mlibrodo/db-copier/postgres/pgcommands/common"
)

var (
	// PGDumpCmd is the path to the `pg_dump` executable
	PSQLCmd = "psql"
)

type PSQLQuery struct {
	*conn.DBConnInfo
	Query string
}

func (x *PSQLQuery) Exec() common.Result {
	execFn := common.PGCLIExecutor(PSQLCmd, x.DBConnInfo, x.ParseArgs)

	return execFn()

}

func (x *PSQLQuery) ParseArgs() []string {
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

	args = append(args, fmt.Sprintf("--command=%s", x.Query))

	return args
}
