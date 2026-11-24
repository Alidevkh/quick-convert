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
