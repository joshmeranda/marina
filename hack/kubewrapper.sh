#!/usr/bin/env bash
# A script to allow you to place your kubebuilder directories wherever you damn well plase without having to build yourself a scaffolding plugin.
#
# Due to the order we move the directories around, the controller apis can be a subdirectory to the controllers dir, but not the other way around.

config_file=kubewrapper.yaml

apis_dir="$(yq eval .apis "$config_file")"
controllers_dir="$(yq eval .controllers "$config_file")"
kubebuilder_bin="$(yq eval .kubebuilder_bin "$config_file")"
manager_main="$(yq eval .manager_main "$config_file")"

usage="Usage: kubewrapper.sh <command> [-f <config_file>]
  Commands:
    setup                        move the specified directories to the kubebuilder locations
    cleanup                      restore diretories to their original locations

  Args:
    --file, -f <config_fiel>     the file to read the config from [kubewrapper.,yaml]
    --help, -h                   show this help text"

setup() {
	if [ -n "$apis_dir" ]; then
		mv "$apis_dir" api
	fi

	if [ -n "$controllers_dir" ]; then
		mv "$controllers_dir" controllers
	fi

	if [ -n "$manager_main" ]; then
		mv "$manager_main" main.go
	fi
}

restore() {
	if [ -n "$controllers_dir" ]; then
		mv controllers "$controllers_dir"
	fi
	
	if [ -n "$apis_dir" ]; then
		mv api "$apis_dir"
	fi

	if [ -n "$manager_main" ]; then
		mv main.go "$manager_main"
	fi
}

if [ $# -eq 0 ]; then
	echo "Expected at least one arg but found $#"
	echo "$usage"
	exit 1
fi

while [ $# -gt 0 ]; do
	case "$1" in
		--file|-f)
			config_file="$2"
			shift
			;;
		--help|-h)
			echo "$usage"
			exit 0
			;;
		setup)
			command=setup
			;;
		restore)
			command=restore
			;;
		*)
			echo "$usage"
			;;
	esac
	shift
done

if [ -z "$command" ]; then
	echo "Expected at least one command but found none"
	echo "$usage"
	exit 1
fi

printf 'running %s\n' "$command"
$command