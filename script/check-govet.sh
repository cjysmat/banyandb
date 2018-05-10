#!/usr/bin/env bash

source "$(dirname "$BASH_SOURCE")/.precommit"

IFS=$'\n'
files=( $(validate_diff --diff-filter=ACMR --name-only -- '*.go' | grep -v '^vendor/' || true) )
unset IFS

errors=()
for f in "${files[@]}"; do
	# we use "git show" here to validate that what's committed passes go vet
	failedVet=$(go tool vet -printf=false "$f")
	if [ "$failedVet" ]; then
		errors+=( "$failedVet" )
	fi
done

if [ ${#errors[@]} -eq 0 ]; then
	echo 'Congratulations!  All Go source files have been vetted.'
else
	{
		echo "Errors from govet:"
		for err in "${errors[@]}"; do
			echo "$err"
		done
		echo
		echo 'Please fix the above errors. You can test via "go vet" and commit the result.'
		echo
	} >&2
	false
fi
