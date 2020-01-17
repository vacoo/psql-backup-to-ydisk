#! /bin/sh
export PGPASSWORD=$PSQL_PASS

gunzip --stdout "$1" | psql -p "$PSQL_PORT" -U "$PSQL_USER" -h "$PSQL_HOST" -d "$2"