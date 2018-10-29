#!/bin/sh

# This script ranges over non-vendored Go packages and runs golint.
# It keeps on until the end and finally exits with 1 if there are any problems.

f=0
for i in $(go list ./... | grep -v /vendor/); do
	golint -set_exit_status "$i"
	f=$((f + $?))
done
if [ "$f" -ne 0 ]; then exit 1; fi
