package tui

import (
	"testing"

	"github.com/ai-cain/websnap/internal/domain"
	tea "github.com/charmbracelet/bubbletea"
)

func TestModelSelectsBrowserTargetAndShowsTabSelection(t *testing.T) {
	t.Parallel()

	studio := &fakeStudio{
		tabsByHandle: map[int64][]domain.BrowserTab{
			131584: {
				{Index: 0, Title: "WhatsApp", Selected: true},
				{Index: 1, Title: "YouTube", Selected: false},
			},
		},
	}

	model := NewModel(studio)
	model.phase = phaseEditing
	model.targets = []targetMenuItem{
		newURLMenuItem(),
		newLiveTargetMenuItem(domain.LiveTarget{
			WindowHandle: 131584,
			Title:        "WhatsApp - Google Chrome",
			AppName:      "chrome",
			Type:         domain.LiveTargetBrowser,
			CanListTabs:  true,
		}),
	}
	model.targetIndex = 1

	next, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	busy := next.(Model)
	if busy.phase != phaseBusy {
		t.Fatalf("phase = %d, want %d", busy.phase, phaseBusy)
	}

	msg := cmd()
	final, _ := busy.Update(msg)
	got := final.(Model)

	if got.screen != screenTabSelect {
		t.Fatalf("screen = %d, want %d", got.screen, screenTabSelect)
	}

	if len(got.tabs) != 2 {
		t.Fatalf("len(tabs) = %d, want 2", len(got.tabs))
	}
}

func TestModelSingleBrowserTabAutoAdvancesToLiveOptions(t *testing.T) {
	t.Parallel()

	studio := &fakeStudio{
		tabsByHandle: map[int64][]domain.BrowserTab{
			131584: {
				{Index: 0, Title: "WhatsApp", Selected: true},
			},
		},
	}

	model := NewModel(studio)
	model.phase = phaseEditing
	model.targets = []targetMenuItem{
		newURLMenuItem(),
		newLiveTargetMenuItem(domain.LiveTarget{
			WindowHandle: 131584,
			Title:        "WhatsApp - Google Chrome",
			AppName:      "chrome",
			Type:         domain.LiveTargetBrowser,
			CanListTabs:  true,
		}),
	}
	model.targetIndex = 1

	next, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	busy := next.(Model)
	msg := cmd()
	final, _ := busy.Update(msg)
	got := final.(Model)

	if got.screen != screenLiveOptions {
		t.Fatalf("screen = %d, want %d", got.screen, screenLiveOptions)
	}

	if !got.hasSelectedTab || got.selectedTab.Index != 0 {
		t.Fatalf("selectedTab = %#v, want first tab auto-selected", got.selectedTab)
	}
}

func TestModelEnterOnLiveOptionsTransitionsToSuccess(t *testing.T) {
	t.Parallel()

	studio := &fakeStudio{
		liveResult: domain.CaptureResult{
			Path:   "C:/captures/live.png",
			Width:  1550,
			Height: 830,
		},
	}

	model := NewModel(studio)
	model.phase = phaseEditing
	model.screen = screenLiveOptions
	model.selectedTarget = domain.LiveTarget{
		WindowHandle: 1312510,
		Title:        "portfolio - Explorador de archivos",
		AppName:      "explorer",
		Type:         domain.LiveTargetFolder,
	}
	model.hasSelectedTarget = true
	model.liveOut.SetValue("captures/live.png")

	next, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	busy := next.(Model)
	if busy.phase != phaseBusy {
		t.Fatalf("phase = %d, want %d", busy.phase, phaseBusy)
	}

	msg := cmd()
	final, _ := busy.Update(msg)
	success := final.(Model)

	if success.phase != phaseSuccess {
		t.Fatalf("phase = %d, want %d", success.phase, phaseSuccess)
	}

	if success.lastPath != "C:/captures/live.png" {
		t.Fatalf("lastPath = %q, want %q", success.lastPath, "C:/captures/live.png")
	}

	if studio.lastLiveReq.Target.WindowHandle != 1312510 {
		t.Fatalf("lastLiveReq = %#v, want selected live target handle", studio.lastLiveReq)
	}
}
