#!/bin/bash

set -e

WORKSPACE="$(pwd -P)"

go build -o "${WORKSPACE}"/bin/jiraexporter "${WORKSPACE}"/main.go
