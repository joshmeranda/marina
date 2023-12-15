#!/usr/bin/env bash

set -e

check_git() {
	local status="$(git status --porcelain)"
	if [ -n "$status" ]; then
		echo git is no longer clean
		echo "$status"	
	fi

	exit 1
}

check() {
	echo running "'$@'"
	$@
	check_git
}

check_git
check go mod tidy
check go fmt ./...