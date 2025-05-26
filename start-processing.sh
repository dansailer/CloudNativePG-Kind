#!/bin/bash
# File: start-processing.sh
go mod tidy
source .env
sleep 5s
go run writer.go
