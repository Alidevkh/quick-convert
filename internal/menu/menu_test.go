package menu

import (
	"regexp"
	"strings"
	"testing"

	"quickconvert/internal/actions"
)

const root = `HKCU\Software\Classes\SystemFileAssociations\`

func TestCRLF(t *testing.T) {
	for name, script := range map[string]string{"install": InstallScript(), "uninstall": UninstallScript()} {
		if strings.Count(script, "\n") != strings.Count(script, "\r\n") {
			t.Errorf("%s: every line must end with CRLF (Windows batch files)", name)
		}
		if !strings.HasPrefix(script, "@echo off\r\n") {
			t.Errorf("%s: must start with @echo off", name)
		}
	}
}

func TestInstallHasMovMp4(t *testing.T) {
	want := `reg add "` + root + `.mov\shell\QuickConvert\shell\01mp4\command" /ve /d "\"%DEST%\quickconvert.exe\" mp4 \"%%1\"" /f >nul`
	if !strings.Contains(InstallScript(), want) {
		t.Errorf("missing line:\n%s", want)
	}
}

func TestInstallCascadeKey(t *testing.T) {
	s := InstallScript()
	for _, want := range []string{
		`reg add "` + root + `.mov\shell\QuickConvert" /v MUIVerb /d "Convert to" /f >nul`,
		`reg add "` + root + `.mov\shell\QuickConvert" /v SubCommands /d "" /f >nul`,
		`reg add "` + root + `.mov\shell\QuickConvert\shell\02mp3" /v MUIVerb /d "MP3 (audio only)" /f >nul`,
		`reg add "` + root + `.mov\shell\QuickConvert\shell\02mp3" /v MultiSelectModel /d "Player" /f >nul`,
	} {
		if !strings.Contains(s, want) {
			t.Errorf("missing line:\n%s", want)
		}
	}
}

func TestInstallNoBarePercent1(t *testing.T) {
	if regexp.MustCompile(`[^%]%1`).MatchString(InstallScript()) {
		t.Error(`a bare %1 in a batch file is eaten by cmd; it must be written %%1`)
	}
}

func TestInstallExtras(t *testing.T) {
	s := InstallScript()
	for _, want := range []string{
		root + `.pdf\shell\QuickConvert.shrink" /v MUIVerb /d "Reduce size"`,
		root + `.heic\shell\QuickConvert.clean" /v MUIVerb /d "Remove metadata"`,
		root + `.pdf\shell\QuickConvert.shrink\command" /ve /d "\"%DEST%\quickconvert.exe\" shrink \"%%1\""`,
	} {
		if !strings.Contains(s, want) {
			t.Errorf("missing:\n%s", want)
		}
	}
}

func TestEveryMenuItemRegistered(t *testing.T) {
	s := InstallScript()
	for _, ext := range actions.AllExts() {
		for _, item := range actions.Submenu(ext) {
			pattern := regexp.QuoteMeta(root+ext+`\shell\QuickConvert\shell\`) + `\d\d` + regexp.QuoteMeta(item.ActionID+`\command" /ve /d "\"%DEST%\quickconvert.exe\" `+item.ActionID+` \"%%1\""`)
			if !regexp.MustCompile(pattern).MatchString(s) {
				t.Errorf("%s %s: no command line", ext, item.ActionID)
			}
		}
		for _, item := range actions.Extras(ext) {
			want := root + ext + `\shell\QuickConvert.` + item.ActionID + `\command" /ve /d "\"%DEST%\quickconvert.exe\" ` + item.ActionID + ` \"%%1\""`
			if !strings.Contains(s, want) {
				t.Errorf("%s %s: no command line", ext, item.ActionID)
			}
		}
	}
}

func TestUninstallCoversInstall(t *testing.T) {
	keyRoot := regexp.MustCompile(regexp.QuoteMeta(root) + `\.[a-z0-9]+\\shell\\QuickConvert(\.[a-z]+)?`)
	un := UninstallScript()
	for _, key := range keyRoot.FindAllString(InstallScript(), -1) {
		if !strings.Contains(un, `reg delete "`+key+`" /f >nul 2>&1`) {
			t.Errorf("uninstall does not remove %s", key)
		}
	}
}

func TestInstallRemovesOldEntriesFirst(t *testing.T) {
	s := InstallScript()
	firstDelete := strings.Index(s, "reg delete")
	firstAdd := strings.Index(s, "reg add")
	if firstDelete < 0 || firstDelete > firstAdd {
		t.Error("install should remove old entries before adding new ones")
	}
}

func TestUninstallRemovesFolderLast(t *testing.T) {
	lines := strings.Split(strings.TrimRight(UninstallScript(), "\r\n"), "\r\n")
	last := lines[len(lines)-1]
	if last != `(goto) 2>nul & rmdir /S /Q "%LOCALAPPDATA%\QuickConvert"` {
		t.Errorf("last line = %q", last)
	}
}

// Explorer hides an entry when more than 15 files are selected unless the
// entry itself (not only its submenu items) allows many files.
func TestParentAllowsManyFiles(t *testing.T) {
	want := `reg add "` + root + `.mov\shell\QuickConvert" /v MultiSelectModel /d "Player" /f >nul`
	if !strings.Contains(InstallScript(), want) {
		t.Errorf("missing:\n%s", want)
	}
}

// Running install.cmd from inside the zip viewer would register a program
// that does not exist; the script must stop and say so.
func TestInstallChecksProgramIsThere(t *testing.T) {
	s := InstallScript()
	check := strings.Index(s, `if not exist "%~dp0quickconvert.exe"`)
	if check < 0 || check > strings.Index(s, "reg add") {
		t.Error("install.cmd must check for quickconvert.exe before adding entries")
	}
}

// rmdir cannot delete the folder the script is running in, so the
// uninstaller must leave it first.
func TestUninstallLeavesFolderBeforeDeleting(t *testing.T) {
	s := UninstallScript()
	cd := strings.Index(s, `cd /d "%TEMP%"`)
	if cd < 0 || cd > strings.Index(s, "rmdir") {
		t.Error("uninstall.cmd must leave its own folder (cd to TEMP) before rmdir")
	}
}

func TestInstallCreatesSendToShortcuts(t *testing.T) {
	s := InstallScript()
	for _, want := range []string{
		`'\Images to PDF.lnk'); $s.TargetPath=$env:DEST+'\quickconvert.exe'; $s.Arguments='imgpdf'; $s.Save()`,
		`'\Merge PDFs.lnk'); $s.TargetPath=$env:DEST+'\quickconvert.exe'; $s.Arguments='merge'; $s.Save()`,
	} {
		if !strings.Contains(s, want) {
			t.Errorf("missing:\n%s", want)
		}
	}
	if strings.Index(s, "Images to PDF.lnk") > strings.Index(s, "echo QuickConvert installed.") {
		t.Error("shortcuts must be created before the 'installed' message")
	}
}

func TestUninstallRemovesSendToShortcuts(t *testing.T) {
	s := UninstallScript()
	for _, want := range []string{
		`del "%APPDATA%\Microsoft\Windows\SendTo\Images to PDF.lnk" 2>nul`,
		`del "%APPDATA%\Microsoft\Windows\SendTo\Merge PDFs.lnk" 2>nul`,
	} {
		i := strings.Index(s, want)
		if i < 0 || i > strings.Index(s, "rmdir") {
			t.Errorf("missing (or after rmdir):\n%s", want)
		}
	}
}

func TestInstallJpgHasPDFEntry(t *testing.T) {
	want := root + `.jpg\shell\QuickConvert\shell\03pdf\command" /ve /d "\"%DEST%\quickconvert.exe\" pdf \"%%1\""`
	if !strings.Contains(InstallScript(), want) {
		t.Errorf("missing:\n%s", want)
	}
}
