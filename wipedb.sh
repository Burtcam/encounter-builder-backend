## used to wipe the local dev db for testing purposes and set it up again in case of

#!/bin/bash
set -e

# 1. Delete the database data and shut down containers.
echo "Taking down containers and deleting volumes..."
# The -v flag ensures that attached volumes (like pgdata) are removed.
docker-compose down -v

# 2. Bring the containers back up.
echo "Starting containers..."
docker-compose up -d

# 3. Wait for PostgreSQL to start up.
echo "Waiting for PostgreSQL to be ready..."
# A simple sleep can work, but you might want to use a more robust health-check.
sleep 10

# 4. Load the schema into the new, empty database.
echo "Loading schema from ./schema/schema.sql ..."
cat ./schema/schema.sql | docker exec -i encounter-builder-postgres psql -U user -d encounterBuilder

echo "Database has been reinitialized with the new schema."