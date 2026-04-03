package tui

import (
	"context"
	"strings"
	"testing"

	"github.com/ai-cain/websnap/internal/domain"
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

func TestGroupSelectionInstructionsMentionAllArrows(t *testing.T) {
	t.Parallel()

	model := NewModel(&fakeStudio{})
	model.screen = screenGroupSelect

	got := model.instructionsForCurrentScreen()
	if !strings.Contains(got, "←/→") {
		t.Fatalf("instructions = %q, want to mention ←/→", got)
	}
}

func TestGroupSelectionRendersNestedCards(t *testing.T) {
	t.Parallel()

	model := NewModel(&fakeStudio{})
	model.phase = phaseEditing
	model.width = 140
	model.groups = []groupMenuItem{
		{title: "Chrome", detail: "1 ventana • navegador"},
		{title: "Antigravity", detail: "3 ventanas • app"},
	}

	view := model.renderGroupSelection()
	if strings.Count(view, "╭") < 2 {
		t.Fatalf("group selection should render nested cards, got view:\n%s", view)
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
		{name: "wide", width: 180, want: 3},
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
