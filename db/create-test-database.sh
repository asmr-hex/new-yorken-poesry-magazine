#! /bin/bash


function print_green() {
  local GREEN='\033[1;32m'
  local NC='\033[0m' # No Color

  printf "${GREEN}$1${NC}"
}

# create a test database
function create_test_database() {
  print_green "    creating database '$TEST_POSTGRES_DB'..."
  # login to psql with dev/prod user/db and create a new test db
  psql -q -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
	    CREATE DATABASE $TEST_POSTGRES_DB;
EOSQL
  print_green "OK.\n"
}

# only create test database if we are running within the dev environment
if [ $DEV_ENV ]; then
  print_green "Initializing DB in DEV_ENV mode\n"
  create_test_database
fi
