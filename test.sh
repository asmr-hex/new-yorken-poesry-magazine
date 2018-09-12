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
server_exit_code=$?

# run npm tests
echo -e "\n\n================ CLIENT TESTS ================"
docker-compose exec dev_client env CI=true npm test
client_exit_code=$?

if [ $server_exit_code -eq 0 -a $client_exit_code -eq 0 ]
then
  echo -e "\n\n================ TESTS PASSED ================"
  exit 0
else
  echo -e "\n\n================ TESTS FAILED ================"
  exit 1
fi

