# Goose

Goose is a package that provides some useful functions to interact with a PostgreSQL database.

### Add a migration using a pgx connection as parameter

This package provide a function to register a migration that use a context and a pgx connection as parameter: `AddPgxContextMigration`
It allows the use of an sqlc querier to interact with the database easily:

```go

import (
    ...
    "github.com/Scalingo/go-utils/postgresql/goose"
    ...
)

func init() {
	err := goose.AddPgxContextMigration(upCreateNode, downCreateNode)
	if err != nil {
		panic(err)
	}
}

func upMigration(ctx context.Context, conn *pgx.Conn) error {
    querier := db.New(conn)
    ...
}

func downMigration(ctx context.Context, conn *pgx.Conn) error {
    querier := db.New(conn)
    ...
}

```
