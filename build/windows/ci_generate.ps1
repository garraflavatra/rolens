param([string]$platform)

mkdir releases
wails build -platform windows/amd64
Compress-Archive -Path build\bin\* -DestinationPath releases\rolens-$platform-amd64.zip

Remove-Item -Recurse -Confirm:$false .\build\bin
wails build -platform windows/arm64
Compress-Archive -Path build\bin\* -DestinationPath releases\rolens-$platform-arm64.zip
