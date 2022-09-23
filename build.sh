# shellcheck disable=SC1113
#/bin/bash
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o room-service  main/main.go
