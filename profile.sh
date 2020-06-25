#!/bin/bash

VEHICLE_API_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"


vehicle () {
    # if no command given force help page
    local OPTION
	if [[ "$1" != "" ]]; then
        OPTION=$1
    else
        OPTION="help"
    fi
	# handle input options
    case "${OPTION}" in
        'help')
echo "Usage: $ ${FUNCNAME} [option]

Options:
- test: stop bili's docker container(s)
- integrations: run integration tests
- bench: bench test server setup
- mock: mock any services
- protoc: run protoc against any proto files"
        ;;
        'test')
          for directory in $( find . -type d | grep '/test$' ); do
            go test "${directory}" -coverpkg=./...
          done
        ;;
        'mock')
          for directory in $( find . -type d | grep '_service$' ); do
            folderName=${directory##*/}
            mockgen \
            -source="${VEHICLE_API_DIR}"/pkg/services/"${folderName}"/interface.go \
            -destination="${VEHICLE_API_DIR}"/pkg/services/"${folderName}"/"${folderName}"_mock.go \
            -package="${folderName}" \
            -mock_names Service=Mock"${folderName}"
          done
        ;;
        'protoc')
          protoc --go_out=plugins=grpc:. protos/*.proto
        ;;
        'integrations')
          VehicleIntegrationTests
        ;;
        'bench')
          go test -bench=.
        ;;
        *)
          echo -e "ERROR: invalid option. Try..\n$ ${FUNCNAME} help"
        ;;
    esac
}

VehicleIntegrationTests () {
  # *****************************************
  # CREATE a new vehicle
  TEST_VIN=$(cat /dev/urandom | env LC_CTYPE=C tr -dc a-zA-Z0-9 | head -c 16; echo)
  VEHICLE_ID=$(_vehicleIntegrationCREATEtest | jq -r ".id")
  # Read the vehicle and make sure it's there
  READ_CREATION_VIN=$(_vehicleIntegrationREADtest | jq -r ".vin")
  # Assert create resulted in the expected vin
  if [[ "${TEST_VIN}" != "${READ_CREATION_VIN}" ]]; then
    echo "assertion failed during insert"
    echo "expected '${TEST_VIN}', got '${READ_VIN}'"
    READ_CREATION_VIN=""
    return
  fi

  # *****************************************
  # UPDATE the existing vehicle with a new vin
  TEST_VIN=$(cat /dev/urandom | env LC_CTYPE=C tr -dc a-zA-Z0-9 | head -c 16; echo)
  VEHICLE_UPDATE_COUNT=$(_vehicleIntegrationUPDATEtest | jq -r ".count_updated")
  # Assert update count != 0
  if [[ "${VEHICLE_UPDATE_COUNT}" == "0" ]]; then
    echo "assertion failed during update"
    echo "no records updated"
    return
  fi
  # READ new vin
  READ_UPDATE_VIN=$(_vehicleIntegrationREADtest | jq -r ".vin")
  # Assert vin updated
  if [[ "${TEST_VIN}" != "${READ_UPDATE_VIN}" ]]; then
    echo "assertion failed during update"
    echo "expected '${TEST_VIN}', got '${READ_UPDATE_VIN}'"
    return
  fi

  # *****************************************
  # DELETE the test vehicle
  VEHICLE_DELETE_COUNT=$(_vehicleIntegrationDELETEtest | jq -r ".count_deleted")
  # Assert delete count != 0
  if [[ "${VEHICLE_DELETE_COUNT}" == "0" ]]; then
    echo "assertion failed during deletion"
    echo "no records deleted"
    return
  fi
  # Read the deleted vehicle's vin
  READ_VIN=$(_vehicleIntegrationREADtest | jq -r ".vin")
  # Assert deleted vehicle is not there
  if [[ "${READ_VIN}" != "" ]]; then
    echo "assertion failed during deletion"
    echo "expected '', got '${READ_VIN}'"
    return
  fi

  echo "ALL TESTS PASSED"
}


_vehicleIntegrationCREATEtest () {
  curl -s -d '{"make":"Make", "model":"Model", "vin":"'"${TEST_VIN}"'"}' \
    -H "Content-Type: application/json" \
    -X POST 'localhost:8080/api/vehicle'
}

_vehicleIntegrationREADtest () {
  curl -s \
    -d id="${VEHICLE_ID}" \
    -X GET -G 'localhost:8080/api/vehicle'
}

_vehicleIntegrationUPDATEtest () {
  curl -s -d '{"id":'"${VEHICLE_ID}"',"make":"Make", "model":"Model", "vin":"'"${TEST_VIN}"'"}' \
    -H "Content-Type: application/json" \
    -X PUT 'localhost:8080/api/vehicle'
}

_vehicleIntegrationDELETEtest () {
  curl -s -d '{"id":'"${VEHICLE_ID}"'}' \
    -H "Content-Type: application/json" \
    -X DELETE 'localhost:8080/api/vehicle'
}
