param([string]$platform)

mkdir releases
wails build -platform windows/amd64 -nsis
Remove-Item build\bin\Rolens.exe
Compress-Archive -Path build\bin\* -DestinationPath releases\rolens-$platform-amd64-installer.zip

Remove-Item -Recurse -Confirm:$false .\build\bin
wails build -platform windows/arm64 -nsis
Remove-Item build\bin\Rolens.exe
Compress-Archive -Path build\bin\* -DestinationPath releases\rolens-$platform-arm64-installer.zip
