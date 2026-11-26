// Package platform holds the small pieces that talk to Windows itself:
// showing notifications and limiting how many heavy jobs run at once.
// Each has a Windows version (*_windows.go) and a stand-in used when the
// code is built and tested on other systems (*_other.go).
package platform

import (
	"encoding/base64"
	"fmt"
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

// AskScript shows a small "Which pages?" box in front of other windows and
// prints what the user typed (nothing if they pressed Cancel). Texts come
// from QC_PROMPT and QC_TITLE; the answer is printed as UTF-8 without a
// byte-order mark, so Arabic digits survive.
const AskScript = `
Add-Type -AssemblyName System.Windows.Forms
[Console]::OutputEncoding = New-Object System.Text.UTF8Encoding $false
$form = New-Object System.Windows.Forms.Form
$form.Text = $env:QC_TITLE
$form.TopMost = $true
$form.StartPosition = 'CenterScreen'
$form.FormBorderStyle = 'FixedDialog'
$form.MinimizeBox = $false
$form.MaximizeBox = $false
$form.ClientSize = New-Object System.Drawing.Size(360, 150)
$label = New-Object System.Windows.Forms.Label
$label.Text = $env:QC_PROMPT
$label.SetBounds(12, 10, 336, 60)
$box = New-Object System.Windows.Forms.TextBox
$box.SetBounds(12, 75, 336, 24)
$ok = New-Object System.Windows.Forms.Button
$ok.Text = 'OK'
$ok.DialogResult = 'OK'
$ok.SetBounds(192, 110, 75, 28)
$cancel = New-Object System.Windows.Forms.Button
$cancel.Text = 'Cancel'
$cancel.DialogResult = 'Cancel'
$cancel.SetBounds(273, 110, 75, 28)
$form.Controls.AddRange(@($label, $box, $ok, $cancel))
$form.AcceptButton = $ok
$form.CancelButton = $cancel
$form.Add_Shown({ $form.Activate(); $box.Focus() })
if ($form.ShowDialog() -eq 'OK') { [Console]::Out.Write($box.Text) }
`

// askProcessFlags says how the PowerShell that shows the box is started:
// CREATE_NO_WINDOW hides its console. HideWindow must NOT be used: Windows
// would apply it to the first window the process shows — the box itself.
func askProcessFlags() (hideWindow bool, creationFlags uint32) {
	return false, 0x08000000
}

// AskPrompt is the question shown in the "Which pages?" box.
func AskPrompt(pages int) string {
	return fmt.Sprintf("Which pages? (1-%d)\n\nExamples:  1-3, 5   or   5, 1-3", pages)
}
