package tui

import (
	"context"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ai-cain/websnap/internal/domain"
	"github.com/charmbracelet/lipgloss"
)

func TestNewModelStartsInGroupSelection(t *testing.T) {
	t.Parallel()

	model := NewModel(&fakeStudio{})

	if model.screen != screenGroupSelect {
		t.Fatalf("screen = %d, want %d", model.screen, screenGroupSelect)
	}

	if model.phase != phaseBusy {
		t.Fatalf("phase = %d, want %d", model.phase, phaseBusy)
	}

	if model.busyLabel == "" {
		t.Fatal("busyLabel should describe target loading")
	}
}

func TestModelLoadsTargetsIntoGroups(t *testing.T) {
	t.Parallel()

	model := NewModel(&fakeStudio{})

	next, _ := model.Update(targetsLoadedMsg{
		targets: []domain.LiveTarget{
			{
				WindowHandle: 100,
				Title:        "README.md - Antigravity",
				AppName:      "antigravity",
				Type:         domain.LiveTargetApp,
			},
			{
				WindowHandle: 101,
				Title:        "home.ts - Antigravity",
				AppName:      "antigravity",
				Type:         domain.LiveTargetApp,
			},
			{
				WindowHandle: 200,
				Title:        "Your Repositories - Google Chrome",
				AppName:      "chrome",
				Type:         domain.LiveTargetBrowser,
				CanListTabs:  true,
			},
		},
	})

	got := next.(Model)
	if got.phase != phaseEditing {
		t.Fatalf("phase = %d, want %d", got.phase, phaseEditing)
	}

	if len(got.groups) != 2 {
		t.Fatalf("len(groups) = %d, want 2", len(got.groups))
	}

	if got.groups[0].title != "Chrome" {
		t.Fatalf("groups[0].title = %q, want Chrome first", got.groups[0].title)
	}

	if got.groups[1].title != "Antigravity" {
		t.Fatalf("groups[1].title = %q, want Antigravity second", got.groups[1].title)
	}
}

func TestModelLoadsEmptyTargetListWithoutGroups(t *testing.T) {
	t.Parallel()

	model := NewModel(&fakeStudio{})
	next, _ := model.Update(targetsLoadedMsg{})
	got := next.(Model)

	if len(got.groups) != 0 {
		t.Fatalf("len(groups) = %d, want 0 when nothing is open", len(got.groups))
	}
}

func TestModelBuildLiveRequest(t *testing.T) {
	t.Parallel()

	model := NewModel(&fakeStudio{})
	model.selectedTarget = domain.LiveTarget{
		WindowHandle: 1312510,
		Title:        "portfolio - Explorador de archivos",
		AppName:      "explorer",
		Type:         domain.LiveTargetFolder,
	}
	model.hasSelectedTarget = true
	model.liveOut.SetValue("captures/folder.png")

	req, err := model.buildLiveRequest()
	if err != nil {
		t.Fatalf("buildLiveRequest() error = %v", err)
	}

	if req.Target.WindowHandle != 1312510 || req.Out != "captures/folder.png" || req.TabIndex != -1 {
		t.Fatalf("request = %#v, want folder target with no tab and custom output", req)
	}
}

func TestSuggestLiveOutputPathPrefersMeaningfulBrowserSegment(t *testing.T) {
	t.Parallel()

	target := domain.LiveTarget{AppName: "chrome", Title: `12°02'15.0"S 76°57'45.7"W - Google Maps`}
	tab := domain.BrowserTab{Title: `12°02'15.0"S 76°57'45.7"W - Google Maps`}

	got := suggestLiveOutputPath(target, tab, true)
	want := filepath.Join("captures", "google-maps.png")
	if got != want {
		t.Fatalf("suggestLiveOutputPath() = %q, want %q", got, want)
	}
}

func TestSuggestLiveOutputPathStripsBrowserSuffix(t *testing.T) {
	t.Parallel()

	target := domain.LiveTarget{AppName: "chrome", Title: "Your Repositories - Google Chrome"}

	got := suggestLiveOutputPath(target, domain.BrowserTab{}, false)
	want := filepath.Join("captures", "your-repositories.png")
	if got != want {
		t.Fatalf("suggestLiveOutputPath() = %q, want %q", got, want)
	}
}

func TestSuggestLiveOutputPathPrefersFolderNameOverGenericExplorerSuffix(t *testing.T) {
	t.Parallel()

	target := domain.LiveTarget{AppName: "explorer", Title: "portfolio - Explorador de archivos"}

	got := suggestLiveOutputPath(target, domain.BrowserTab{}, false)
	want := filepath.Join("captures", "portfolio.png")
	if got != want {
		t.Fatalf("suggestLiveOutputPath() = %q, want %q", got, want)
	}
}

func TestGroupSelectionInstructionsMentionAllArrows(t *testing.T) {
	t.Parallel()

	model := NewModel(&fakeStudio{})
	model.screen = screenGroupSelect

	got := model.instructionsForCurrentScreen()
	if !strings.Contains(got, "←/→") {
		t.Fatalf("instructions = %q, want to mention ←/→", got)
	}
}

func TestGroupSelectionRendersCards(t *testing.T) {
	t.Parallel()

	model := NewModel(&fakeStudio{})
	model.phase = phaseEditing
	model.width = 140
	model.groups = []groupMenuItem{
		{title: "Chrome", detail: "1 ventana • navegador"},
		{title: "Antigravity", detail: "3 ventanas • app"},
	}

	view := model.renderGroupSelection()
	borders := strings.Count(view, "┌") + strings.Count(view, "╭")
	if borders < 2 {
		t.Fatalf("group selection should render cards, got view:\n%s", view)
	}
}

func TestGroupGridColumnsRespondToWidth(t *testing.T) {
	t.Parallel()

	model := NewModel(&fakeStudio{})
	model.groups = []groupMenuItem{
		{title: "Chrome"},
		{title: "Antigravity"},
		{title: "Explorer"},
		{title: "Settings"},
		{title: "Terminal"},
	}

	tests := []struct {
		name  string
		width int
		want  int
	}{
		{name: "narrow", width: 80, want: 1},
		{name: "medium", width: 120, want: 2},
		{name: "wide", width: 180, want: 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model.width = tt.width
			if got := model.groupGridColumns(); got != tt.want {
				t.Fatalf("groupGridColumns() = %d, want %d for width %d", got, tt.want, tt.width)
			}
		})
	}
}

func TestRenderChoiceCardKeepsSingleContentRow(t *testing.T) {
	t.Parallel()

	model := NewModel(&fakeStudio{})
	card := model.renderChoiceCard("Chrome", "1 ventana • navegador", true, 42)

	if strings.Count(card, "\n") != 2 {
		t.Fatalf("renderChoiceCard() should stay compact with border-only height, got:\n%s", card)
	}
}

func TestRenderFooterShowsQuitOnTopLevel(t *testing.T) {
	t.Parallel()

	model := NewModel(&fakeStudio{})
	model.screen = screenGroupSelect

	footer := model.renderFooter()
	if !strings.Contains(footer, "Esc") || !strings.Contains(footer, "quit") {
		t.Fatalf("renderFooter() = %q, want Esc quit on top-level screen", footer)
	}
	if strings.Contains(footer, "Esc   back") {
		t.Fatalf("renderFooter() = %q, should not advertise Esc back on top-level screen", footer)
	}
}

func TestRenderFooterShowsBackOnNestedScreens(t *testing.T) {
	t.Parallel()

	model := NewModel(&fakeStudio{})
	model.screen = screenTargetSelect

	footer := model.renderFooter()
	if !strings.Contains(footer, "Esc") || !strings.Contains(footer, "back") {
		t.Fatalf("renderFooter() = %q, want Esc back on nested screen", footer)
	}
}

func TestViewFitsWithinViewportWidth(t *testing.T) {
	t.Parallel()

	model := NewModel(&fakeStudio{})
	model.phase = phaseEditing
	model.width = 120
	model.groups = []groupMenuItem{
		{title: "Chrome", detail: "1 ventana • navegador"},
		{title: "Explorador", detail: "1 ventana • carpeta"},
		{title: "Antigravity", detail: "2 ventanas • app"},
		{title: "Windows Host", detail: "1 ventana • app"},
	}

	view := model.View()
	maxWidth := 0
	for _, line := range strings.Split(view, "\n") {
		if width := lipgloss.Width(line); width > maxWidth {
			maxWidth = width
		}
	}

	if maxWidth > model.width {
		t.Fatalf("view width = %d, want <= viewport width %d\n%s", maxWidth, model.width, view)
	}
}

type fakeStudio struct {
	targets      []domain.LiveTarget
	tabsByHandle map[int64][]domain.BrowserTab
	liveResult   domain.CaptureResult
	liveErr      error
	lastLiveReq  domain.LiveCaptureRequest
}

func (f *fakeStudio) ListTargets(_ context.Context) ([]domain.LiveTarget, error) {
	return f.targets, nil
}

func (f *fakeStudio) ListTabs(_ context.Context, target domain.LiveTarget) ([]domain.BrowserTab, error) {
	if f.tabsByHandle == nil {
		return nil, nil
	}

	return f.tabsByHandle[target.WindowHandle], nil
}

func (f *fakeStudio) CaptureLive(_ context.Context, req domain.LiveCaptureRequest) (domain.CaptureResult, error) {
	f.lastLiveReq = req
	return f.liveResult, f.liveErr
}
