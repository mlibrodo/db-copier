package copier

import (
	"github.com/doug-martin/goqu"
	"github.com/mlibrodo/db-copier/log"
	"github.com/mlibrodo/db-copier/postgres/conn"
	"github.com/mlibrodo/db-copier/postgres/pgcommands"
)

type TrimDef struct {
	Table  string
	Filter string
}

func (in TrimDef) ToSQL() string {
	ds := goqu.From(in.Table).Where(goqu.Literal(in.Filter)).Delete()

	return ds.Sql
}

type Trim struct {
	TrimDefs []TrimDef
}

func (in Trim) Exec(pgConnInfo *conn.DBConnInfo) error {

	for _, d := range in.TrimDefs {
		log.WithFields(log.Fields{
			"table":  d.Table,
			"filter": d.Filter,
		}).Infof("Running %v", d.ToSQL())

		query := pgcommands.PSQLQuery{
			pgConnInfo,
			d.ToSQL(),
		}

		result := query.Exec()
		if result.Error != nil {
			log.WithFields(
				log.Fields{
					"Command": result.FullCommand,
					"Error":   result.Error.Err,
				},
			).Debug("Trim failed")

			log.Error(result.Output)

			return result.Error.Err
		}
	}

	return nil
}
