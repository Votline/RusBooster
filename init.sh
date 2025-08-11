#!/bin/bash
set -e

PG_USER="${POSTGRES_USER:-postgres}"
PG_PASSWORD="${POSTGRES_PASSWORD}"
PG_DB="${POSTGRES_DB:-postgres}"
DUMP_FILE="/tmp/rsdb_dump.sql"

base_psql() {
    PGPASSWORD="$PG_PASSWORD" psql -v ON_ERROR_STOP=1 -U "$PG_USER" -d "$1" -c "$2"
}

echo ">>> Starting database initialization..."

echo ">>> Ensuring database exists: $PG_DB"
base_psql "postgres" "DO \$\$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_database WHERE datname = '$PG_DB') THEN
        CREATE DATABASE $PG_DB OWNER $PG_USER;
    END IF;
END
\$\$;"

if [ -f "$DUMP_FILE" ]; then
    echo ">>> Restoring database from dump..."
    TMP_DUMP="/tmp/rsdb_dump_fixed.sql"
    sed "s/OWNER TO postgres/OWNER TO ${PG_USER}/g" "$DUMP_FILE" > "$TMP_DUMP"
    PGPASSWORD="$PG_PASSWORD" psql -v ON_ERROR_STOP=1 -U "$PG_USER" -d "$PG_DB" -f "$TMP_DUMP"
    rm "$TMP_DUMP"
else
    echo ">>> WARNING: Dump file $DUMP_FILE not found!"
fi

echo ">>> Database initialization complete!"

