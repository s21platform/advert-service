#!/bin/bash

# Применим миграции
# ВАЖНО! Для локального запуска вызвать из корня проекта, предварительно подгрузив переменные окружения в терминал командой  `set -a; source <путь к .env файлу>; set +a`
# Установка goose: go install github.com/pressly/goose/v3/cmd/goose@latest
goose -dir ./migrations postgres "user=$ADVERT_SERVICE__POSTGRES_USER password=$ADVERT_SERVICE__POSTGRES_PASSWORD dbname=$ADVERT_SERVICE__POSTGRES_DB host=$ADVERT_SERVICE__POSTGRES_HOST port=$ADVERT_SERVICE__POSTGRES_PORT sslmode=disable" up-by-one