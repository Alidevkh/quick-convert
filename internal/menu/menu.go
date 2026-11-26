// Package menu writes install.cmd and uninstall.cmd, the two scripts that
// add and remove QuickConvert's right-click menu entries.
//
// The entries are built from the action list (package actions), so every
// menu entry matches something the program can do. Everything goes under
// HKEY_CURRENT_USER: no administrator rights are needed.
//
// What one file type gets in the registry (example: .mov):
//
//	...\SystemFileAssociations\.mov\shell\QuickConvert          "Convert to" (a submenu)
//	...\SystemFileAssociations\.mov\shell\QuickConvert\shell\01mp4\command
//	    = "C:\Users\me\AppData\Local\QuickConvert\quickconvert.exe" mp4 "%1"
//
// Windows replaces %1 with the file that was right-clicked.
package menu

import (
	"fmt"
	"strings"

	"quickconvert/internal/actions"
)

// classesRoot is where per-user file-type menu entries live.
const classesRoot = `HKCU\Software\Classes\SystemFileAssociations\`

// InstallScript returns the text of install.cmd.
func InstallScript() string {
	var s script
	s.add("@echo off")
	s.add("rem Installs QuickConvert for the current user (no admin needed).")
	s.add("setlocal")
	s.add(`set "DEST=%LOCALAPPDATA%\QuickConvert"`)
	s.add("")
	s.add("rem 0. Stop if run from inside the zip (Windows shows the zip like a folder).")
	s.add(`if not exist "%~dp0quickconvert.exe" goto :notunzipped`)
	s.add("")
	s.add("rem 1. Remove menu entries from any older version.")
	addDeletes(&s)
	s.add("")
	s.add("rem 2. Copy the program files (skipped when run from the install folder).")
	s.add(`if /I not "%~dp0"=="%DEST%\" robocopy "%~dp0." "%DEST%" /E /NFL /NDL /NJH /NJS /NP >nul`)
	// Files unzipped from a download carry Windows' "from the internet" mark;
	// SmartScreen then stops the Send to shortcuts ("unknown publisher").
	// Remove the mark from QuickConvert's own installed files only.
	s.add(`powershell -NoProfile -ExecutionPolicy Bypass -Command "Get-ChildItem -LiteralPath $env:DEST -Recurse | Unblock-File"`)
	s.add("")
	s.add("rem 3. Add the right-click menu entries.")
	for _, ext := range actions.AllExts() {
		addEntries(&s, ext)
	}
	s.add("")
	s.add(`rem 4. Add the "Send to" shortcuts (cmd cannot make shortcuts, PowerShell can).`)
	for _, item := range actions.SendToItems() {
		s.add(sendToShortcut(item))
	}
	s.add("")
	s.add("echo QuickConvert installed.")
	s.add("echo Right-click a file ^> Show more options ^> Convert to.")
	s.add("pause")
	s.add("exit /b 0")
	s.add("")
	s.add(":notunzipped")
	s.add("echo Please unzip QuickConvert first, then run install.cmd from the unzipped folder.")
	s.add("pause")
	s.add("exit /b 1")
	return s.String()
}

// SelfTestScript returns the text of selftest.cmd, which runs the
// PowerShell self-test (selftest.ps1) next to it.
func SelfTestScript() string {
	var s script
	s.add("@echo off")
	s.add("rem Installs QuickConvert and tests every conversion on this computer.")
	s.add(`powershell -NoProfile -ExecutionPolicy Bypass -File "%~dp0selftest.ps1" %*`)
	s.add("pause")
	return s.String()
}

// UninstallScript returns the text of uninstall.cmd.
func UninstallScript() string {
	var s script
	s.add("@echo off")
	s.add("rem Removes QuickConvert's menu entries and its folder.")
	addDeletes(&s)
	for _, item := range actions.SendToItems() {
		s.add(`del "%APPDATA%\Microsoft\Windows\SendTo\` + item.Label + `.lnk" 2>nul`)
	}
	s.add("echo QuickConvert menu entries removed.")
	s.add(`cd /d "%TEMP%"`) // leave the folder, or Windows won't let us delete it
	// Deleting the folder this script may be running from: "(goto)" ends
	// the script first, so cmd doesn't complain that the file is gone.
	s.add(`(goto) 2>nul & rmdir /S /Q "%LOCALAPPDATA%\QuickConvert"`)
	return s.String()
}

// sendToShortcut returns the command that creates "<Label>.lnk" in the
// user's Send to folder, running "quickconvert.exe <id>". Windows then adds
// every selected file after <id>. DEST is read from the environment
// ($env:DEST), which install.cmd set just before.
func sendToShortcut(item actions.SendTo) string {
	return `powershell -NoProfile -ExecutionPolicy Bypass -Command "` +
		`$s=(New-Object -ComObject WScript.Shell).CreateShortcut([Environment]::GetFolderPath('SendTo')+'\` + item.Label + `.lnk'); ` +
		`$s.TargetPath=$env:DEST+'\quickconvert.exe'; $s.Arguments='` + item.ActionID + `'; $s.Save()"`
}

// addEntries writes the registry lines for one file extension.
func addEntries(s *script, ext string) {
	if items := actions.Submenu(ext); len(items) > 0 {
		key := classesRoot + ext + `\shell\QuickConvert`
		s.regValue(key, "MUIVerb", "Convert to")
		s.regValue(key, "SubCommands", "")            // empty = the submenu is the "shell" key below
		s.regValue(key, "MultiSelectModel", "Player") // still shown when more than 15 files are selected
		for n, item := range items {
			verb := fmt.Sprintf(`%s\shell\%02d%s`, key, n+1, item.ActionID) // 01, 02... keeps menu order
			s.regVerb(verb, item)
		}
	}
	for _, item := range actions.Extras(ext) {
		s.regVerb(classesRoot+ext+`\shell\QuickConvert.`+item.ActionID, item)
	}
}

// addDeletes writes "reg delete" for every top-level key we may create.
func addDeletes(s *script) {
	for _, ext := range actions.AllExts() {
		s.regDelete(classesRoot + ext + `\shell\QuickConvert`)
		// Every top-level entry on every type, not only today's: an older
		// version may have added one that no longer exists (e.g.
		// "Remove metadata" on .mkv).
		for _, id := range extraIDs() {
			s.regDelete(classesRoot + ext + `\shell\QuickConvert.` + id)
		}
	}
}

// extraIDs lists every top-level action ID used on any file type.
func extraIDs() []string {
	seen := map[string]bool{}
	var ids []string
	for _, ext := range actions.AllExts() {
		for _, item := range actions.Extras(ext) {
			if !seen[item.ActionID] {
				seen[item.ActionID] = true
				ids = append(ids, item.ActionID)
			}
		}
	}
	return ids
}

// script builds a Windows batch file. Batch files need CRLF line endings.
type script struct{ b strings.Builder }

func (s *script) add(line string) { s.b.WriteString(line + "\r\n") }

func (s *script) String() string { return s.b.String() }

// regValue sets a named text value on a key.
func (s *script) regValue(key, name, value string) {
	s.add(fmt.Sprintf(`reg add "%s" /v %s /d "%s" /f >nul`, key, name, value))
}

// regVerb writes one menu entry: its label, multi-file support, and the
// command line. In a batch file "%" must be doubled, so %1 is written %%1,
// and \" puts a quote inside the value.
func (s *script) regVerb(key string, item actions.MenuItem) {
	s.regValue(key, "MUIVerb", item.Label)
	s.regValue(key, "MultiSelectModel", "Player") // allow selecting many files at once
	s.add(fmt.Sprintf(`reg add "%s\command" /ve /d "\"%%DEST%%\quickconvert.exe\" %s \"%%%%1\"" /f >nul`, key, item.ActionID))
}

func (s *script) regDelete(key string) {
	s.add(fmt.Sprintf(`reg delete "%s" /f >nul 2>&1`, key))
}
