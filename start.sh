#!/bin/sh

set -e

echo "running database migrations"
source /usr/src/app/.env
/usr/src/app/migrate -path /usr/src/app/migrations -database "$DB_URL" -verbose up

echo "starting the app"
exec "$@"