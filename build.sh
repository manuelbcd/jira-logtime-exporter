#!/bin/bash

set -e

WORKSPACE="$(pwd -P)"

cd "${WORKSPACE}"/cmd

go build -o "${WORKSPACE}"/bin/jiraexporter "${WORKSPACE}"/cmd/main.go