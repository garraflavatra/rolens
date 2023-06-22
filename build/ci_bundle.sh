#!/bin/sh

# Available platforms:
# - windows-2019
# - windows-2022
# - macos-11
# - macos-12
# - macos-13
# - ubuntu-20.04
# - ubuntu-22.04

version=$(<./build/version.txt)

mkdir bundle

# mv artifacts/*/rolens-macos-11.zip bundle/
# mv artifacts/*/rolens-macos-11.zip bundle/

ls
echo "$version"
