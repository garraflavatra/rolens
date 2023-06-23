#!/bin/sh

# Cleanup
rm -rf releases
rm -rf build/bin
mkdir releases
mkdir -p build/bin

# Settings
cat > build/darwin/dmg_settings.json << EOF
{
  "title": "Rolens",
  "background": "$(pwd)/build/darwin/dmg_background.png",
  "icon-size": 100,
  "window": {
    "size": { "width": 155, "height": 250 },
    "position": { "x": 360, "y": 360 }
  },
  "contents": [
    { "x": 750, "y": 500, "type": "link", "path": "/Applications" },
    { "x": 595, "y": 250, "type": "file", "path": "$(pwd)/build/bin/Rolens.app" }
  ]
}
EOF

# AMD/Intel
wails build -platform darwin/amd64
# create-dmg \
#   --volname Rolens \
#   --window-size 155 250 \
#   --volicon build/appicon.png \
#   --eula LICENSE \
#   --app-drop-link 750 500 \
#   --icon-size 100 \
#   --background build/darwin/dmg_background.png \
#   --add-file Rolens.app build/bin/Rolens.app 595 250 \
#   build/bin/Rolens.dmg emptydir
# appdmg build/darwin/dmg_settings.json build/bin/Rolens.dmg
tar -czvf releases/rolens-$1-amd64.tar.gz --directory build/bin Rolens.app

# Cleanup
rm -rf build/bin/Rolens.app

# ARM/AppleM1
wails build -platform darwin/arm64
# create-dmg \
#   --volname Rolens \
#   --window-size 155 250 \
#   --volicon build/appicon.png \
#   --eula LICENSE \
#   --app-drop-link 750 500 \
#   --icon-size 100 \
#   --background build/darwin/dmg_background.png \
#   --add-file Rolens.app build/bin/Rolens.app 595 250 \
#   build/bin/Rolens.dmg emptydir
# appdmg build/darwin/dmg_settings.json build/bin/Rolens.dmg
tar -czvf releases/rolens-$1-arm64.tar.gz --directory build/bin Rolens.app

# Cleanup
rm -rf build/bin/Rolens.app
