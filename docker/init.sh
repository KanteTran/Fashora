#!/bin/bash

# Start PostgreSQL in the background
docker-entrypoint.sh postgres &

# Wait for PostgreSQL to be ready
until pg_isready -U "$POSTGRES_USER" -d "$POSTGRES_DB"; do
  echo "Waiting for PostgreSQL to be ready..."
  sleep 2
done

# Run the SQL file
psql -U "$POSTGRES_USER" -d "$POSTGRES_DB" -f /postgresql/init_db.sql

# Keep the container running
wait