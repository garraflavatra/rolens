mkdir releases
mkdir emptydir

# AMD/Intel
wails build -platform darwin/amd64
create-dmg \
  --volname Rolens \
  --window-size 155 250 \
  --volicon build/appicon.png \
  --eula LICENSE \
  --app-drop-link 750 500 \
  --icon-size 100 \
  --background build/darwin/dmg_background.png \
  --add-file Rolens.app build/bin/Rolens.app 595 250 \
  build/bin/Rolens.dmg emptydir
tar -czvf releases/rolens-$1-amd64.tar.gz --directory build/bin Rolens.dmg

# ARM/AppleM1
rm -rf build/bin
wails build -platform darwin/arm64
create-dmg \
  --volname Rolens \
  --window-size 155 250 \
  --volicon build/appicon.png \
  --eula LICENSE \
  --app-drop-link 750 500 \
  --icon-size 100 \
  --background build/darwin/dmg_background.png \
  --add-file Rolens.app build/bin/Rolens.app 595 250 \
  build/bin/Rolens.dmg emptydir
tar -czvf releases/rolens-$1-arm64.tar.gz --directory build/bin Rolens.dmg
