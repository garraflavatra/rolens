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
    "size": { "width": 750, "height": 400 }
  },
  "contents": [
    { "x": 600, "y": 175, "type": "link", "path": "/Applications" },
    { "x": 150, "y": 175, "type": "file", "path": "$(pwd)/build/bin/Rolens.app" }
  ]
}
EOF

# AMD/Intel
wails build -platform darwin/amd64
appdmg build/darwin/dmg_settings.json build/bin/Rolens.dmg
zip -j releases/rolens-$1-amd64.zip build/bin/Rolens.dmg

# Cleanup
rm -rf build/bin/Rolens.dmg

# ARM/AppleM1
wails build -platform darwin/arm64
appdmg build/darwin/dmg_settings.json build/bin/Rolens.dmg
zip -j releases/rolens-$1-arm64.zip build/bin/Rolens.dmg

# Cleanup
rm -rf build/bin/Rolens.app
