#!/bin/bash

FILE=.env
if test -f "$FILE"; then
  echo "Everything is OK"
  echo "Validate dependencies . . ."
  go mod tidy -compat=1.19
  echo "Trying to run the linter & tests . . ."
  staticcheck ./...
  go vet ./...
  go test ./... -v
  echo "Trying to run the server . . ."
  go run ./cmd/web/main.go
else
  echo "==========================================================="
  echo "|  $FILE (environment) file does not exist.               |"
  echo "|  Please Crete new .env file from .env.example.          |"
  echo "|  by running this script: //:~$ cp .env.example .env     |"
  echo "==========================================================="
  exit 0
fi