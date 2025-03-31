#!/bin/sh

set -e

echo "running database migrations"
/usr/src/app/migrate -path /usr/src/app/migrations -database "$DB_URL" -verbose up

echo "starting the app"
exec "$@"