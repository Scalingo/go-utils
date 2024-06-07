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

// AddPgxContextMigration registers a migration in the goose migrations list using a pgx.Conn instead of a sql.DB
// It can only be used with a pgx driver
func AddPgxContextMigration(upMigration PgxContextMigration, downMigration PgxContextMigration) error {
	// runtime.Caller returns the filename of the migration file that calls this function
	// it is used to save the migration in the goose migrations list
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
		// Raw is a function that allows to run a function with the raw driver connection as argument
		err = sqlConn.Raw(func(driverConn interface{}) error {
			// stdlibConn is the raw driver connection casted to a pgx stdlib connection
			// stdlib is a package that wraps the pgx driver to make it compatible with the sql.DB interface
			stdlibConn, ok := driverConn.(*stdlib.Conn)
			if !ok {
				return errors.New(ctx, "could not cast the driver connection to a pgx.Conn")
			}
			// pgxConn is the pgx connection extracted from the stdlib connection
			pgxConn := stdlibConn.Conn()
			return pgxContextMigration(ctx, pgxConn)
		})
		if err != nil {
			return errors.Wrap(ctx, err, "could not run the pgx migration")
		}
		return nil
	}
}
