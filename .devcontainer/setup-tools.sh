#!/bin/bash
set -e

go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
go install github.com/google/wire/cmd/wire@latest
