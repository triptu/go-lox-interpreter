#!/bin/sh
# builds and runs the interpreter
set -e # Exit early if any commands fail
(
  cd "$(dirname "$0")" # Ensure compile steps are run within the repository directory
  go build -o ./build/golox ./cmd/main.go
)
exec ./build/golox "$@"
