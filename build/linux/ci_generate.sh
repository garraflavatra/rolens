#!/bin/sh

mkdir releases
wails build -platform linux/amd64
tar -czvf releases/rolens-$1-amd64.tar.gz --directory build/bin Rolens

# rm -rf build/bin
# wails build -platform linux/arm64
# tar -czvf releases/rolens-$1-arm64.tar.gz --directory build/bin Rolens
