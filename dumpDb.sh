#!/bin/sh

PGHOST="$db_host"
PGDATABASE="$POSTGRES_DB"
PGUSER="$POSTGRES_USER"
PGPASSWORD="$POSTGRES_PASSWORD"

DUMP_FILE_NAME="backupOn`date +%Y-%m-%d-%H-%M`.dump"
echo "Creating dump: $DUMP_FILE_NAME"

cd pg_backup

pg_dump -w --format=c --blobs > $DUMP_FILE_NAME

if [ $? -ne 0 ]; then
  rm $DUMP_FILE_NAME
  echo "Back up not created, check db connection settings"
  exit 1
fi

echo 'Successfully Backed Up'
exit 0