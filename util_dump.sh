#! /bin/sh

set -e
set -o pipefail

export PGPASSWORD=$PSQL_PASS

PSQL_HOST_OPTS="-h $PSQL_HOST -p $PSQL_PORT -U $PSQL_USER"
DATE=$(date +%Y-%m-%d_%H-%M)
FILEPATH="backups/$PSQL_DB-$DATE.gz"

pg_dump $PSQL_HOST_OPTS $PSQL_DB | gzip > $FILEPATH

echo $FILEPATH