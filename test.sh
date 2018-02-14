#! /bin/bash

# run all tests within dockerized dev env
#
# NOTE: you must have first run docker-compose up and
# allowed the stack to spin up before executing this script.
#
# ALSO NOTE: this script only runs the tests once, not in watch-mode.

# run golang tests
echo "================ SERVER TESTS ================"
docker-compose exec dev_server gotest -v ./...

# run npm tests
echo -e "\n\n================ CLIENT TESTS ================"
docker-compose exec dev_client env CI=true npm test


