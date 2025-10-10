#!/bin/sh
# Espera até que o Redis esteja pronto
until nc -z -v -w30 redis 6379
do
  echo "Aguardando Redis..."
  sleep 1
done

echo "Redis está pronto. Iniciando a aplicação..."
exec "$@"

