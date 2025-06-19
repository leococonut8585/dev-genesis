!define PRODUCT_NAME "Dev Genesis"
!define PRODUCT_VERSION "1.0.0"
!define PRODUCT_PUBLISHER "Leo Sakaguchi"
!define PRODUCT_WEB_SITE "https://github.com/leococonut8585/dev-genesis"

!include "MUI2.nsh"

Name "${PRODUCT_NAME} ${PRODUCT_VERSION}"
OutFile "../build/DevGenesisInstaller.exe"
InstallDir "$PROGRAMFILES64\DevGenesis"
RequestExecutionLevel admin

!define MUI_ABORTWARNING
!define MUI_ICON "../assets/icon.ico"
!define MUI_UNICON "../assets/icon.ico"

!insertmacro MUI_PAGE_WELCOME
!insertmacro MUI_PAGE_LICENSE "../LICENSE"
!insertmacro MUI_PAGE_DIRECTORY
!insertmacro MUI_PAGE_INSTFILES
!insertmacro MUI_PAGE_FINISH

!insertmacro MUI_UNPAGE_WELCOME
!insertmacro MUI_UNPAGE_CONFIRM
!insertmacro MUI_UNPAGE_INSTFILES
!insertmacro MUI_UNPAGE_FINISH

!insertmacro MUI_LANGUAGE "English"
!insertmacro MUI_LANGUAGE "Japanese"

Section "MainSection" SEC01
    SetOutPath "$INSTDIR"
    File "../build/dev-genesis-windows-amd64.exe"
    File "../backend/scripts/install/*.ps1"
    
    ; Create shortcuts
    CreateDirectory "$SMPROGRAMS\Dev Genesis"
    CreateShortcut "$SMPROGRAMS\Dev Genesis\Dev Genesis.lnk" "$INSTDIR\dev-genesis-windows-amd64.exe"
    CreateShortcut "$DESKTOP\Dev Genesis.lnk" "$INSTDIR\dev-genesis-windows-amd64.exe"
    
    ; Add to PATH
    EnVar::SetHKLM
    EnVar::AddValue "PATH" "$INSTDIR"
SectionEnd

Section -Post
    WriteUninstaller "$INSTDIR\uninst.exe"
    WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${PRODUCT_NAME}" "DisplayName" "${PRODUCT_NAME}"
    WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${PRODUCT_NAME}" "UninstallString" "$INSTDIR\uninst.exe"
    WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${PRODUCT_NAME}" "DisplayVersion" "${PRODUCT_VERSION}"
    WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${PRODUCT_NAME}" "Publisher" "${PRODUCT_PUBLISHER}"
SectionEnd

Section Uninstall
    Delete "$INSTDIR\*.*"
    RMDir "$INSTDIR"
    Delete "$SMPROGRAMS\Dev Genesis\*.*"
    RMDir "$SMPROGRAMS\Dev Genesis"
    Delete "$DESKTOP\Dev Genesis.lnk"
    DeleteRegKey HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${PRODUCT_NAME}"
SectionEnd