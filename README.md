# go-cookbook

Go (golang) cookbook for common patterns and requirements e.g. PostgreSQL, logging, gRPC,
JWT auth etc.

Each folder contains a specific example for a given use case and is described below.

Some examples use [https://docs.docker.com/compose/install/](Docker compose) to run servers e.g. postgres,
so you _should_ have it installed.

** Absolutely no guarantee that these examples are bug free or production ready! **

## Cookbooks

### sql

Uses the standard library and postgres driver to read and write data, with a simple repository pattern.
The entity/table is a simple NHS Pharmacy object taken from public UK data.

Uses docker compose and PostgreSQL to create the tables.

See [https://pkg.go.dev/database/sql](https://pkg.go.dev/database/sql)


### sql-pgx

Uses the popular [https://github.com/jackc/pgx](pgx) library and postgres driver to read and write data, with a simple repository pattern.
The entity/table is a simple NHS Pharmacy object taken from public UK data.

Uses docker compose and PostgreSQL to create the tables.

### sql-pgx-pool

Same as [sql-pgx](#sql-pgx), but uses a custom connection pool instead of pgx.Connect() method

### logging

Use standard library logger to set various formatting options and also how to write to a file, and both std out and file
