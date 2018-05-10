#!/usr/bin/env bash

source "$(dirname "$BASH_SOURCE")/.precommit"

IFS=$'\n'
files=( $(validate_diff --diff-filter=ACMR --name-only -- '*.go' | grep -v '^vendor/\|autogen' || true) )
unset IFS

errors=()
for f in "${files[@]}"; do
	# we use "git show" here to validate that what's committed passes go vet
	failedLint=$(golint "$f")
	if [ "$failedLint" ]; then
		errors+=( "$failedLint" )
	fi
done

if [ ${#errors[@]} -eq 0 ]; then
	echo 'Congratulations!  All Go source files have been linted.'
else
	{
		echo "Errors from golint:"
		for err in "${errors[@]}"; do
			echo "$err"
		done
		echo
		echo 'Please fix the above errors. You can test via "golint" and commit the result.'
		echo
	} >&2
	false
fi
