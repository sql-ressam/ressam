# Ressam [WIP]

Ressam (Turkish meaning artist) - is lightweight CLI database diagram explorer 
with relationships prediction (a.k.a. virtual foreign key).
Ressam can show database diagrams in the self-hosted Web app.
No telemetry, external calls, fully open.

Example usage:

```sh
export RESSAM_DRIVER=pg # aliases: postgres, postgresql
export RESSAM_DSN="postgresql://user:password@127.0.0.1:5432/ressam?sslmode=disable"
ressam draw
```

TODO:

* MySQL support
* MSSQL support
* Mongo support?
* Export to JSON, YAML, UML, XML
* Relationship prediction

# Contribute

## run unit tests

```bash
go test -run=. ./... 
```

## run integration tests

1. run test migrations
```bash
goose -dir=pg/testdata/migrations postgres "postgresql://postgres:postgres@localhost:5432/ressam?sslmode=disable" up
```

or

```bash
docker compose up -d
```

2. Run integration tests
```
export TEST_RESSAM_DSN_PG="postgresql://postgres:postgres@localhost:5432/ressam?sslmode=disable"
go test -tags integration -run=. ./... 
```