mkdir releases
wails build -platform linux/amd64
tar -czvf releases/rolens-${{ matrix.platform }}-amd64.tar.gz --directory build/bin Rolens

# rm -rf build/bin
# wails build -platform linux/arm64
# tar -czvf releases/rolens-${{ matrix.platform }}-arm64.tar.gz --directory build/bin Rolens
