package tui

import (
	"context"
	"testing"

	"github.com/ai-cain/websnap/internal/domain"
)

func TestNewModelStartsInTargetSelection(t *testing.T) {
	t.Parallel()

	model := NewModel(&fakeStudio{})

	if model.screen != screenTargetSelect {
		t.Fatalf("screen = %d, want %d", model.screen, screenTargetSelect)
	}

	if model.phase != phaseBusy {
		t.Fatalf("phase = %d, want %d", model.phase, phaseBusy)
	}

	if model.busyLabel == "" {
		t.Fatal("busyLabel should describe target loading")
	}
}

func TestModelLoadsTargetsIntoMenu(t *testing.T) {
	t.Parallel()

	model := NewModel(&fakeStudio{})

	next, _ := model.Update(targetsLoadedMsg{
		targets: []domain.LiveTarget{
			{
				WindowHandle: 131584,
				Title:        "WhatsApp - Google Chrome",
				AppName:      "chrome",
				Type:         domain.LiveTargetBrowser,
				CanListTabs:  true,
			},
			{
				WindowHandle: 1312510,
				Title:        "portfolio - Explorador de archivos",
				AppName:      "explorer",
				Type:         domain.LiveTargetFolder,
			},
		},
	})

	got := next.(Model)
	if got.phase != phaseEditing {
		t.Fatalf("phase = %d, want %d", got.phase, phaseEditing)
	}

	if len(got.targets) != 2 {
		t.Fatalf("len(targets) = %d, want 2 live targets only", len(got.targets))
	}

	if got.targets[0].kind != menuItemLiveTarget {
		t.Fatalf("targets[0].kind = %d, want %d", got.targets[0].kind, menuItemLiveTarget)
	}
}

func TestModelLoadsEmptyTargetListWithoutURLFallback(t *testing.T) {
	t.Parallel()

	model := NewModel(&fakeStudio{})
	next, _ := model.Update(targetsLoadedMsg{})
	got := next.(Model)

	if len(got.targets) != 0 {
		t.Fatalf("len(targets) = %d, want 0 when nothing is open", len(got.targets))
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
