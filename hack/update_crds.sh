#!/usr/bin/env sh

set -e

crd_repo=https://github.com/joshmeranda/marina-operator.git
commit=eb3f05accc451ad0835da2b26c68f4976c80bbc6
tag_or_branch=main

temp_dir=$(mktemp --directory)
srd_crd_subdir=config/crd/bases
dst_crd_dir=crds

trap "rm --recursive --force $temp_dir" EXIT

if [ -n "$tag_or_branch" ]; then
	git clone --branch "$tag_or_branch" "$crd_repo" "$temp_dir"
else
	git clone "$crd_repo" "$temp_dir"
fi

if [ -n "$commit" ]; then
	cd "$temp_dir"
	git reset --hard "$commit"
	exit 1
	cd -
fi

git clone $crd_repo "$temp_dir"
cp --recursive --force "$temp_dir/$srd_crd_subdir" "$dst_crd_dir"