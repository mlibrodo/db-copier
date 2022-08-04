package copier

import (
	"github.com/doug-martin/goqu/v9"
	"github.com/mlibrodo/db-copier/log"
	"github.com/mlibrodo/db-copier/postgres/conn"
	"github.com/mlibrodo/db-copier/postgres/pgcommands"
)

type AnonymizeFn uint

func (a AnonymizeFn) name() string {
	return anonymize_fn_names[a]
}

func (a AnonymizeFn) ordinal() int {
	return int(a)
}

func (a AnonymizeFn) values() *[]string {
	return &anonymize_fn_names
}

const (
	ANONYMIZE_STRING AnonymizeFn = iota
	ANONYMIZE_INT
)

var anonymize_fn_names = []string{
	"ANONYMIZE_STRING",
	"ANONYMIZE_INT",
}

type AnonymizeDef struct {
	Table       string
	Column      string
	AnonymizeFn AnonymizeFn
}

func (in AnonymizeDef) ToSQL() string {
	anonymizeFn := goqu.Literal(anonymize_postgres_map[in.AnonymizeFn])
	ds := goqu.Update(in.Table).Set(goqu.Record{in.Column: anonymizeFn})
	sql, _, _ := ds.ToSQL()
	return sql
}

type Anonymize struct {
	AnonymizeDefs []AnonymizeDef
}

func (in Anonymize) Exec(pgConnInfo *conn.DBConnInfo) error {

	for _, d := range in.AnonymizeDefs {
		log.WithFields(log.Fields{
			"table":     d.Table,
			"column":    d.Column,
			"anonymize": d.AnonymizeFn.name(),
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
			).Debug("Anonymize failed")

			log.Error(result.Output)

			return result.Error.Err
		}
	}

	return nil
}
