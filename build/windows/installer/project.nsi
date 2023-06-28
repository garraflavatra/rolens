Unicode true

## Please note: Template replacements don't work in this file. They are provided with default defines like
## mentioned underneath.
## If the keyword is not defined, "wails_tools.nsh" will populate them with the values from ProjectInfo.
## If they are defined here, "wails_tools.nsh" will not touch them. This allows to use this project.nsi manually
## from outside of Wails for debugging and development of the installer.

## For development first make a wails nsis build to populate the "wails_tools.nsh":
## > wails build --target windows/amd64 --nsis
## Then you can call makensis on this file with specifying the path to your binary:
## For a AMD64 only installer:
## > makensis -DARG_WAILS_AMD64_BINARY=..\..\bin\app.exe
## For a ARM64 only installer:
## > makensis -DARG_WAILS_ARM64_BINARY=..\..\bin\app.exe
## For a installer with both architectures:
## > makensis -DARG_WAILS_AMD64_BINARY=..\..\bin\app-amd64.exe -DARG_WAILS_ARM64_BINARY=..\..\bin\app-arm64.exe

## The following information is taken from the ProjectInfo file, but can be overwritten here.
## !define INFO_PROJECTNAME    "MyProject"           # Default "{{.Name}}"
## !define INFO_COMPANYNAME    "MyCompany"           # Default "{{.Info.CompanyName}}"
## !define INFO_PRODUCTNAME    "MyProduct"           # Default "{{.Info.ProductName}}"
## !define INFO_PRODUCTVERSION "1.0.0"               # Default "{{.Info.ProductVersion}}"
## !define INFO_COPYRIGHT      "Copyright"           # Default "{{.Info.Copyright}}"
## !define PRODUCT_EXECUTABLE  "Application.exe"     # Default "${INFO_PROJECTNAME}.exe"
## !define UNINST_KEY_NAME     "UninstKeyInRegistry" # Default "${INFO_COMPANYNAME}${INFO_PRODUCTNAME}"
## !define REQUEST_EXECUTION_LEVEL "admin"           # Default "admin"  see also https://nsis.sourceforge.io/Docs/Chapter4.html

!include "wails_tools.nsh"

# The version information for this two must consist of 4 parts
VIProductVersion "${INFO_PRODUCTVERSION}.0"
VIFileVersion    "${INFO_PRODUCTVERSION}.0"

# Product information
VIAddVersionKey "CompanyName"     "${INFO_COMPANYNAME}"
VIAddVersionKey "FileDescription" "${INFO_PRODUCTNAME} Installer"
VIAddVersionKey "ProductVersion"  "${INFO_PRODUCTVERSION}"
VIAddVersionKey "FileVersion"     "${INFO_PRODUCTVERSION}"
VIAddVersionKey "LegalCopyright"  "${INFO_COPYRIGHT}"
VIAddVersionKey "ProductName"     "${INFO_PRODUCTNAME}"

!include "MUI2.nsh"
!define MUI_ICON "..\icon.ico"
!define MUI_UNICON "..\icon.ico"

# Bitmap on the left side of the welcome page. Must be 164x314 pixels in size.
!define MUI_HEADERIMAGE
!define MUI_HEADERIMAGE_BITMAP ".\banner_h.bmp"
!define MUI_WELCOMEFINISHPAGE_BITMAP ".\banner_v.bmp"
!define MUI_WELCOMEPAGE_TITLE "Welcome to the Rolens installer!"

# Finish page information
!define MUI_FINISHPAGE_RUN "$INSTDIR\${INFO_PROJECTNAME}.exe"
!define MUI_FINISHPAGE_RUN_TEXT "Start Rolens when finished"
!define MUI_FINISHPAGE_TITLE "Thanks for installing!"
!define MUI_FINISHPAGE_LINK "Visit Rolens on the Web!"
!define MUI_FINISHPAGE_LINK_LOCATION "https://garraflavatra.github.io/rolens/"
!define MUI_FINISHPAGE_LINK_COLOR 880000

!define MUI_FINISHPAGE_NOAUTOCLOSE # Wait on the INSTFILES page so the user can take a look into the details of the installation steps
!define MUI_ABORTWARNING # This will warn the user if they exit from the installer.

!insertmacro MUI_PAGE_WELCOME
# !insertmacro MUI_PAGE_LICENSE "resources\eula.txt"
!insertmacro MUI_PAGE_DIRECTORY
!insertmacro MUI_PAGE_INSTFILES
!insertmacro MUI_PAGE_FINISH

!insertmacro MUI_UNPAGE_INSTFILES

!insertmacro MUI_LANGUAGE "English"

## The following two statements can be used to sign the installer and the uninstaller. The path to the binaries are provided in %1
#!uninstfinalize 'signtool --file "%1"'
#!finalize 'signtool --file "%1"'

Name "${INFO_PRODUCTNAME}"
OutFile "..\..\bin\${INFO_PROJECTNAME}-${ARCH}-installer.exe" # Name of the installer's file.
InstallDir "$PROGRAMFILES64\${INFO_PRODUCTNAME}" # Default installing folder ($PROGRAMFILES is Program Files folder).
ShowInstDetails show # This will always show the installation details.

Function .onInit
   !insertmacro wails.checkArchitecture
FunctionEnd

Section
    !insertmacro wails.webview2runtime

    SetOutPath $INSTDIR

    !insertmacro wails.files

    CreateShortcut "$SMPROGRAMS\${INFO_PRODUCTNAME}.lnk" "$INSTDIR\${PRODUCT_EXECUTABLE}"
    CreateShortCut "$DESKTOP\${INFO_PRODUCTNAME}.lnk" "$INSTDIR\${PRODUCT_EXECUTABLE}"

    !insertmacro wails.writeUninstaller
SectionEnd

Section "uninstall"
    RMDir /r "$AppData\${PRODUCT_EXECUTABLE}" # Remove the WebView2 DataPath

    RMDir /r $INSTDIR

    Delete "$SMPROGRAMS\${INFO_PRODUCTNAME}.lnk"
    Delete "$DESKTOP\${INFO_PRODUCTNAME}.lnk"

    !insertmacro wails.deleteUninstaller
SectionEnd
