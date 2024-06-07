package utils

import (
	"context"
	"database/sql"
	"runtime"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

type PgxContextMigration func(context.Context, *pgx.Conn) error

func AddPgxContextMigration(upMigration PgxContextMigration, downMigration PgxContextMigration) {
	// runtime.Caller returns the filename of the migration file that calls this function
	_, filename, _, _ := runtime.Caller(1)
	goose.AddNamedMigrationNoTxContext(
		filename,
		wrapNoTxContextMigrationToPgxContextMigration(upMigration),
		wrapNoTxContextMigrationToPgxContextMigration(downMigration),
	)
}

func wrapNoTxContextMigrationToPgxContextMigration(pgxContextMigration PgxContextMigration) goose.GoMigrationNoTxContext {
	return func(ctx context.Context, db *sql.DB) error {
		rawSQLConn, err := db.Conn(ctx)
		if err != nil {
			return err
		}
		err = rawSQLConn.Raw(func(driverConn interface{}) error {
			pgxConn := driverConn.(*stdlib.Conn).Conn()
			return pgxContextMigration(ctx, pgxConn)
		})
		if err != nil {
			return err
		}
		return nil
	}
}
