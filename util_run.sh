#! /bin/sh
set -e
export PGPASSWORD=$PSQL_PASS

if [ "${SCHEDULE}" = "**None**" ]; then
  ./app backup
else
  exec go-cron "$SCHEDULE" ./app backup
fi
