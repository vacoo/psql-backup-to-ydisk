#! /bin/sh

set -e

export PGPASSWORD=$PSQL_PASS

# if [ "${SCHEDULE}" = "**None**" ]; then
#   sh backup.sh
# else
#   exec go-cron "$SCHEDULE" /bin/sh backup.sh
# fi

if [ "${SCHEDULE}" = "**None**" ]; then
  ./backup 
else
  exec go-cron "$SCHEDULE" ./backup 
fi