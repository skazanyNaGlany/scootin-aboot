#!/bin/bash

rm coverage.out 2>/dev/null

go test -coverprofile=coverage.out -coverpkg=../api/... && go tool cover -html=coverage.out -o=coverage.html && echo "Coverage saved to ./coverage.html"

rm coverage.out 2>/dev/null
