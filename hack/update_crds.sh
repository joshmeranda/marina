#!/usr/bin/env sh

set -e

crd_repo=https://github.com/joshmeranda/marina-operator.git
commit=e29fc3be06c1
tag_or_branch=main

temp_dir=$(mktemp --directory update-crds.XXXXXXXXXX)
src_crd_subdir=config/crd/bases
dst_crd_dir=$(pwd)/crds

cleanup() {	
	echo cleaning up temp dir "$temp_dir"
	rm --recursive --force "$temp_dir"
}

trap cleanup EXIT

if [ -n "$tag_or_branch" ]; then
	git clone --branch "$tag_or_branch" "$crd_repo" "$temp_dir"
else
	git clone "$crd_repo" "$temp_dir"
fi

if [ -n "$commit" ]; then
	cd "$temp_dir"
	git reset --hard "$commit"
	cd -
fi

cp --verbose --recursive "$temp_dir/$src_crd_subdir" "$dst_crd_dir"