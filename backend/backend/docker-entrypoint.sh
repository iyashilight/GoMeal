#!/bin/sh
set -e

# Wait for MySQL
echo "Waiting for MySQL..."
while ! nc -z mysql 3306 2>/dev/null; do
    sleep 1
done
echo "MySQL is ready."

# Wait for Redis
echo "Waiting for Redis..."
while ! nc -z redis 6379 2>/dev/null; do
    sleep 1
done
echo "Redis is ready."

exec ./server
