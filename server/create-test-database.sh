#! /bin/bash


function create_database() {
  local database=$POSTGRES_TEST_DB
  echo "  Creating database '$database'"
  psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB"<<-EOSQL
	    CREATE DATABASE $database;
EOSQL
}

create_database $db

# if [ -n "$POSTGRES_MULTIPLE_DATABASES" ]; then
#   echo "Multiple database creation requested: $POSTGRES_MULTIPLE_DATABASES"
#   for db in $(echo $POSTGRES_MULTIPLE_DATABASES | tr ',' ' '); do
#     create_database $db
#   done
#   echo "test database created"
# fi
