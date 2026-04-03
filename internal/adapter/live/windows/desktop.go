package windows

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"unicode/utf16"

	"github.com/ai-cain/websnap/internal/domain"
	apperrors "github.com/ai-cain/websnap/internal/support/errors"
)

type executor interface {
	Run(ctx context.Context, script string) ([]byte, error)
}

type Desktop struct {
	exec executor
}

func New() *Desktop {
	return newDesktop(powerShellExecutor{binary: "powershell"})
}

func newDesktop(exec executor) *Desktop {
	return &Desktop{exec: exec}
}

func (d *Desktop) ListTargets(ctx context.Context) ([]domain.LiveTarget, error) {
	if d == nil || d.exec == nil {
		return nil, apperrors.New(apperrors.CodeInvalidArgument, "desktop target catalog is not configured")
	}

	var payload []struct {
		WindowHandle int64                 `json:"windowHandle"`
		Title        string                `json:"title"`
		AppName      string                `json:"appName"`
		Type         domain.LiveTargetType `json:"type"`
		CanListTabs  bool                  `json:"canListTabs"`
	}

	if err := d.runJSON(ctx, listTargetsScript(), &payload); err != nil {
		return nil, apperrors.Wrap(apperrors.CodeCaptureFailed, "failed to enumerate open targets", err)
	}

	targets := make([]domain.LiveTarget, 0, len(payload))
	for _, item := range payload {
		targets = append(targets, domain.LiveTarget{
			WindowHandle: item.WindowHandle,
			Title:        item.Title,
			AppName:      item.AppName,
			Type:         item.Type,
			CanListTabs:  item.CanListTabs,
		})
	}

	return targets, nil
}

func (d *Desktop) ListTabs(ctx context.Context, target domain.LiveTarget) ([]domain.BrowserTab, error) {
	if d == nil || d.exec == nil {
		return nil, apperrors.New(apperrors.CodeInvalidArgument, "desktop target catalog is not configured")
	}

	if target.WindowHandle <= 0 {
		return nil, apperrors.New(apperrors.CodeInvalidArgument, "window handle is required")
	}

	var payload []struct {
		Index    int    `json:"index"`
		Title    string `json:"title"`
		Selected bool   `json:"selected"`
	}

	if err := d.runJSON(ctx, listTabsScript(target.WindowHandle), &payload); err != nil {
		return nil, apperrors.Wrap(apperrors.CodeCaptureFailed, "failed to enumerate browser tabs", err)
	}

	tabs := make([]domain.BrowserTab, 0, len(payload))
	for _, item := range payload {
		tabs = append(tabs, domain.BrowserTab{
			Index:    item.Index,
			Title:    item.Title,
			Selected: item.Selected,
		})
	}

	return tabs, nil
}

func (d *Desktop) Capture(ctx context.Context, req domain.LiveCaptureRequest) (domain.LiveCaptureImage, error) {
	if d == nil || d.exec == nil {
		return domain.LiveCaptureImage{}, apperrors.New(apperrors.CodeInvalidArgument, "desktop live capturer is not configured")
	}

	if err := req.Validate(); err != nil {
		return domain.LiveCaptureImage{}, err
	}

	var payload struct {
		Width     int64  `json:"width"`
		Height    int64  `json:"height"`
		PNGBase64 string `json:"pngBase64"`
	}

	if err := d.runJSON(ctx, captureWindowScript(req), &payload); err != nil {
		return domain.LiveCaptureImage{}, apperrors.Wrap(apperrors.CodeCaptureFailed, "failed to capture selected live target", err)
	}

	png, err := base64.StdEncoding.DecodeString(strings.TrimSpace(payload.PNGBase64))
	if err != nil {
		return domain.LiveCaptureImage{}, apperrors.Wrap(apperrors.CodeCaptureFailed, "failed to decode captured image", err)
	}

	return domain.LiveCaptureImage{
		PNG:    png,
		Width:  payload.Width,
		Height: payload.Height,
	}, nil
}

func (d *Desktop) runJSON(ctx context.Context, script string, target any) error {
	output, err := d.exec.Run(ctx, script)
	if err != nil {
		return err
	}

	trimmed := strings.TrimSpace(string(output))
	if trimmed == "" {
		trimmed = "[]"
	}

	if err := json.Unmarshal([]byte(trimmed), target); err != nil {
		return fmt.Errorf("invalid powershell json output: %w", err)
	}

	return nil
}

type powerShellExecutor struct {
	binary string
}

func (e powerShellExecutor) Run(ctx context.Context, script string) ([]byte, error) {
	encoded := encodeUTF16LEBase64(script)
	cmd := exec.CommandContext(ctx, e.binary, "-NoProfile", "-NonInteractive", "-EncodedCommand", encoded)
	output, err := cmd.Output()
	if err == nil {
		return output, nil
	}

	if exitErr, ok := err.(*exec.ExitError); ok {
		stderr := strings.TrimSpace(string(exitErr.Stderr))
		if stderr != "" {
			return nil, fmt.Errorf("powershell failed: %s", stderr)
		}
	}

	return nil, err
}

func encodeUTF16LEBase64(value string) string {
	encoded := utf16.Encode([]rune(value))
	bytes := make([]byte, 0, len(encoded)*2)
	for _, r := range encoded {
		bytes = append(bytes, byte(r), byte(r>>8))
	}

	return base64.StdEncoding.EncodeToString(bytes)
}

func listTargetsScript() string {
	return `
Add-Type @"
using System;
using System.Text;
using System.Runtime.InteropServices;
public static class Win32 {
  public delegate bool EnumWindowsProc(IntPtr hWnd, IntPtr lParam);
  [DllImport("user32.dll")] public static extern bool EnumWindows(EnumWindowsProc lpEnumFunc, IntPtr lParam);
  [DllImport("user32.dll", SetLastError=true)] public static extern int GetWindowText(IntPtr hWnd, StringBuilder text, int maxCount);
  [DllImport("user32.dll")] public static extern int GetWindowTextLength(IntPtr hWnd);
  [DllImport("user32.dll")] [return: MarshalAs(UnmanagedType.Bool)] public static extern bool IsWindowVisible(IntPtr hWnd);
  [DllImport("user32.dll")] public static extern uint GetWindowThreadProcessId(IntPtr hWnd, out uint processId);
}
"@

$browserProcesses = @("chrome", "msedge", "brave", "opera")
$results = New-Object System.Collections.Generic.List[object]

$callback = [Win32+EnumWindowsProc]{
  param($hWnd, $lParam)

  if (-not [Win32]::IsWindowVisible($hWnd)) { return $true }

  $len = [Win32]::GetWindowTextLength($hWnd)
  if ($len -le 0) { return $true }

  $sb = New-Object System.Text.StringBuilder ($len + 1)
  [void][Win32]::GetWindowText($hWnd, $sb, $sb.Capacity)
  $title = $sb.ToString().Trim()
  if (-not $title -or $title -eq "Program Manager") { return $true }

  [uint32]$processId = 0
  [void][Win32]::GetWindowThreadProcessId($hWnd, [ref]$processId)
  if ($processId -eq 0) { return $true }

  try {
    $process = Get-Process -Id $processId -ErrorAction Stop
  } catch {
    return $true
  }

  $appName = $process.ProcessName.ToLowerInvariant()
  $type = "app"
  $canListTabs = $false

  if ($browserProcesses -contains $appName) {
    $type = "browser"
    $canListTabs = $true
  } elseif ($appName -eq "explorer") {
    $type = "folder"
  }

  $results.Add([pscustomobject]@{
    windowHandle = $hWnd.ToInt64()
    title = $title
    appName = $appName
    type = $type
    canListTabs = $canListTabs
  })

  return $true
}

[void][Win32]::EnumWindows($callback, [IntPtr]::Zero)
$results |
  Sort-Object @{Expression="type"; Ascending=$true}, @{Expression="title"; Ascending=$true} |
  ConvertTo-Json -Depth 4
`
}

func listTabsScript(handle int64) string {
	return fmt.Sprintf(`
Add-Type -AssemblyName UIAutomationClient
$hwnd = [IntPtr]%d
$root = [System.Windows.Automation.AutomationElement]::FromHandle($hwnd)
if ($null -eq $root) {
  @() | ConvertTo-Json -Depth 4
  exit 0
}

$stripCondition = New-Object System.Windows.Automation.PropertyCondition(
  [System.Windows.Automation.AutomationElement]::ClassNameProperty,
  'HorizontalTabStripRegionView'
)
$tabStrip = $root.FindFirst([System.Windows.Automation.TreeScope]::Descendants, $stripCondition)
if ($null -eq $tabStrip) {
  @() | ConvertTo-Json -Depth 4
  exit 0
}

$itemCondition = New-Object System.Windows.Automation.PropertyCondition(
  [System.Windows.Automation.AutomationElement]::ControlTypeProperty,
  [System.Windows.Automation.ControlType]::TabItem
)
$items = $tabStrip.FindAll([System.Windows.Automation.TreeScope]::Descendants, $itemCondition)
$result = for ($i = 0; $i -lt $items.Count; $i++) {
  $item = $items.Item($i)
  $pattern = $null
  $selected = $false
  if ($item.TryGetCurrentPattern([System.Windows.Automation.SelectionItemPattern]::Pattern, [ref]$pattern)) {
    $selected = $pattern.Current.IsSelected
  }

  [pscustomobject]@{
    index = $i
    title = $item.Current.Name
    selected = $selected
  }
}

$result | ConvertTo-Json -Depth 4
`, handle)
}

func captureWindowScript(req domain.LiveCaptureRequest) string {
	return fmt.Sprintf(`
Add-Type -AssemblyName UIAutomationClient
Add-Type -AssemblyName System.Drawing
Add-Type @"
using System;
using System.Runtime.InteropServices;
public static class Win32 {
  [StructLayout(LayoutKind.Sequential)]
  public struct RECT {
    public int Left;
    public int Top;
    public int Right;
    public int Bottom;
  }

  [DllImport("user32.dll")] public static extern bool GetWindowRect(IntPtr hWnd, out RECT rect);
  [DllImport("user32.dll")] public static extern bool SetForegroundWindow(IntPtr hWnd);
}
"@

$hwnd = [IntPtr]%d
$tabIndex = %d
[void][Win32]::SetForegroundWindow($hwnd)
Start-Sleep -Milliseconds 200

if ($tabIndex -ge 0) {
  $root = [System.Windows.Automation.AutomationElement]::FromHandle($hwnd)
  if ($null -ne $root) {
    $stripCondition = New-Object System.Windows.Automation.PropertyCondition(
      [System.Windows.Automation.AutomationElement]::ClassNameProperty,
      'HorizontalTabStripRegionView'
    )
    $tabStrip = $root.FindFirst([System.Windows.Automation.TreeScope]::Descendants, $stripCondition)
    if ($null -ne $tabStrip) {
      $itemCondition = New-Object System.Windows.Automation.PropertyCondition(
        [System.Windows.Automation.AutomationElement]::ControlTypeProperty,
        [System.Windows.Automation.ControlType]::TabItem
      )
      $items = $tabStrip.FindAll([System.Windows.Automation.TreeScope]::Descendants, $itemCondition)
      if ($tabIndex -lt $items.Count) {
        $item = $items.Item($tabIndex)
        $selectionPattern = $null
        if ($item.TryGetCurrentPattern([System.Windows.Automation.SelectionItemPattern]::Pattern, [ref]$selectionPattern)) {
          $selectionPattern.Select()
        } else {
          $invokePattern = $null
          if ($item.TryGetCurrentPattern([System.Windows.Automation.InvokePattern]::Pattern, [ref]$invokePattern)) {
            $invokePattern.Invoke()
          }
        }
      }
    }
  }

  Start-Sleep -Milliseconds 350
}

$rect = New-Object Win32+RECT
[void][Win32]::GetWindowRect($hwnd, [ref]$rect)
$width = $rect.Right - $rect.Left
$height = $rect.Bottom - $rect.Top

$bitmap = New-Object System.Drawing.Bitmap $width, $height
$graphics = [System.Drawing.Graphics]::FromImage($bitmap)
$graphics.CopyFromScreen($rect.Left, $rect.Top, 0, 0, $bitmap.Size)

$stream = New-Object System.IO.MemoryStream
$bitmap.Save($stream, [System.Drawing.Imaging.ImageFormat]::Png)
$base64 = [System.Convert]::ToBase64String($stream.ToArray())

$graphics.Dispose()
$bitmap.Dispose()
$stream.Dispose()

[pscustomobject]@{
  width = $width
  height = $height
  pngBase64 = $base64
} | ConvertTo-Json -Compress
`, req.Target.WindowHandle, req.TabIndex)
}
