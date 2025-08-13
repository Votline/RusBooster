#! /bin/bash
set -e

if [ -f .env ]; then
	export $(grep -v "^#" .env | xargs)
fi

DUMP_FILE="./rsdb_dump.sql"

echo ">>> Создаю дамп базы данных перед остановкой..."
docker exec rusbooster_db pg_dump -U "$POSTGRES_USER" -d "$POSTGRES_DB" > "$DUMP_FILE"

if [ $? -eq 0 ]; then
	echo ">>> Дамп успешно сохранён в $DUMP_FILE"
else
	echo ">>> Ошибка при создании дампа " >&2
	exit 1
fi

echo ">>> Останавливаю контейнеры..."
docker compose down
