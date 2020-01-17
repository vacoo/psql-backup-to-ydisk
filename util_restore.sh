
if [ "$#" -ne 1 ]
then
    echo "You must pass the path of the backup file to restore"
fi

export PGPASSWORD=$PSQL_PASS
echo "=> Восстановление дампа из $1"
set -o pipefail
gunzip --stdout "$1" | psql -h -U "$PSQL_USER" -h "$PSQL_HOST" -d "$PSQL_DB" -f "$1"
echo "=> Успешное восстановление"