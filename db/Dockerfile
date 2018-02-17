FROM postgres

# add a script to the docker entrypoint which will be executed when the
# container is run. This script will optionally setup a test database if
# the container is not running in prod.
COPY create-test-database.sh /docker-entrypoint-initdb.d/
