package goose

import (
	"context"
	"database/sql"
	"runtime"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"

	"github.com/Scalingo/go-utils/errors/v2"
)

type PgxContextMigration func(context.Context, *pgx.Conn) error

// AddPgxContextMigration adds a migration that uses a pgx.Conn instead of a sql.DB
// It can only be used with a pgx driver
func AddPgxContextMigration(upMigration PgxContextMigration, downMigration PgxContextMigration) error {
	// runtime.Caller returns the filename of the migration file that calls this function
	_, filename, _, ok := runtime.Caller(1)
	if !ok {
		return errors.New(context.Background(), "could not get caller filename")
	}
	goose.AddNamedMigrationNoTxContext(
		filename,
		wrapNoTxContextMigrationToPgxContextMigration(upMigration),
		wrapNoTxContextMigrationToPgxContextMigration(downMigration),
	)
	return nil
}

func wrapNoTxContextMigrationToPgxContextMigration(pgxContextMigration PgxContextMigration) goose.GoMigrationNoTxContext {
	return func(ctx context.Context, db *sql.DB) error {
		// sqlConn is a single connection to the database coming from the sql.DB pool
		sqlConn, err := db.Conn(ctx)
		if err != nil {
			return errors.Wrap(ctx, err, "could not get a connection from the connection pool")
		}
		// Raw is a function that allows to run a function with the underlying driver connection as argument
		err = sqlConn.Raw(func(driverConn interface{}) error {
			// pgxConn is the underlying driver connection that casted to a pgx.Conn
			pgxConn := driverConn.(*stdlib.Conn).Conn()
			return pgxContextMigration(ctx, pgxConn)
		})
		if err != nil {
			return errors.Wrap(ctx, err, "could not run the pgx migration")
		}
		return nil
	}
}
