// Package platform holds the small pieces that talk to Windows itself:
// showing notifications and limiting how many heavy jobs run at once.
// Each has a Windows version (*_windows.go) and a stand-in used when the
// code is built and tested on other systems (*_other.go).
package platform

import (
	"encoding/base64"
	"unicode/utf16"
)

// ToastScript is the PowerShell script that shows a Windows notification.
// The title and text are read from the environment variables QC_TITLE and
// QC_BODY, so file names with quotes or other symbols cannot break it.
//
// Windows only shows notifications for registered apps; we borrow the
// registration of Windows PowerShell, so the sender shows as
// "Windows PowerShell".
const ToastScript = `
[Windows.UI.Notifications.ToastNotificationManager, Windows.UI.Notifications, ContentType = WindowsRuntime] > $null
$template = [Windows.UI.Notifications.ToastNotificationManager]::GetTemplateContent([Windows.UI.Notifications.ToastTemplateType]::ToastText02)
$texts = $template.GetElementsByTagName('text')
$texts.Item(0).AppendChild($template.CreateTextNode($env:QC_TITLE)) > $null
$texts.Item(1).AppendChild($template.CreateTextNode($env:QC_BODY)) > $null
$appId = '{1AC14E77-02E7-4E5D-B744-2EB1AE5198B7}\WindowsPowerShell\v1.0\powershell.exe'
[Windows.UI.Notifications.ToastNotificationManager]::CreateToastNotifier($appId).Show([Windows.UI.Notifications.ToastNotification]::new($template))
`

// encodeCommand prepares a script for "powershell -EncodedCommand", which
// expects the text as UTF-16 little-endian bytes, base64-encoded.
func encodeCommand(script string) string {
	units := utf16.Encode([]rune(script))
	bytes := make([]byte, 0, len(units)*2)
	for _, u := range units {
		bytes = append(bytes, byte(u), byte(u>>8))
	}
	return base64.StdEncoding.EncodeToString(bytes)
}
