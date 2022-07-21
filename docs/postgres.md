# Postgres support in broln

With the introduction of the `kvdb` interface, broln can support multiple database
backends. One of the supported backends is Postgres. This document
describes how it can be configured.

## Building broln with postgres support

To build broln with postgres support, include the following build tag:

```shell
⛰  make tags="kvdb_postgres"
```

## Configuring Postgres for broln

In order for broln to run on Postgres, an empty database should already exist. A
database can be created via the usual ways (psql, pgadmin, etc). A user with
access to this database is also required.

Creation of a schema and the tables is handled by broln automatically.

## Configuring broln for Postgres

broln is configured for Postgres through the following configuration options:

* `db.backend=postgres` to select the Postgres backend.
* `db.postgres.dsn=...` to set the database connection string that includes
  database, user and password.
* `db.postgres.timeout=...` to set the connection timeout. If not set, no
  timeout applies.
