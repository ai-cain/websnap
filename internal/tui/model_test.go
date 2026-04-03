package tui

import (
	"context"
	"testing"

	"github.com/ai-cain/websnap/internal/domain"
	tea "github.com/charmbracelet/bubbletea"
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

	if len(got.targets) != 3 {
		t.Fatalf("len(targets) = %d, want 3 including URL entry", len(got.targets))
	}

	if got.targets[0].kind != menuItemURL {
		t.Fatalf("targets[0].kind = %d, want %d", got.targets[0].kind, menuItemURL)
	}
}

func TestModelCanOpenURLFormFromTargetMenu(t *testing.T) {
	t.Parallel()

	model := NewModel(&fakeStudio{})
	model.phase = phaseEditing
	model.targets = []targetMenuItem{
		newURLMenuItem(),
	}

	next, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	got := next.(Model)

	if got.screen != screenURLForm {
		t.Fatalf("screen = %d, want %d", got.screen, screenURLForm)
	}

	if cmd == nil {
		t.Fatal("entering URL form should return a blink command")
	}
}

func TestModelBuildRequest(t *testing.T) {
	t.Parallel()

	model := NewModel(&fakeStudio{})
	model.phase = phaseEditing
	model.screen = screenURLForm
	model.urlInputs[inputIndex(fieldURL)].SetValue("https://example.com")
	model.urlInputs[inputIndex(fieldWidth)].SetValue("1600")
	model.urlInputs[inputIndex(fieldHeight)].SetValue("1000")
	model.urlInputs[inputIndex(fieldSelector)].SetValue("#app")
	model.urlInputs[inputIndex(fieldOut)].SetValue("captures/home.png")
	model.mode = modeSelector

	req, err := model.buildRequest()
	if err != nil {
		t.Fatalf("buildRequest() error = %v", err)
	}

	expected := domain.CaptureRequest{
		URL:      "https://example.com",
		Width:    1600,
		Height:   1000,
		Selector: "#app",
		Out:      "captures/home.png",
	}

	if req != expected {
		t.Fatalf("request = %#v, want %#v", req, expected)
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
	targets        []domain.LiveTarget
	tabsByHandle   map[int64][]domain.BrowserTab
	urlResult      domain.CaptureResult
	liveResult     domain.CaptureResult
	urlErr         error
	liveErr        error
	lastURLRequest domain.CaptureRequest
	lastLiveReq    domain.LiveCaptureRequest
}

func (f *fakeStudio) CaptureURL(_ context.Context, req domain.CaptureRequest) (domain.CaptureResult, error) {
	f.lastURLRequest = req
	return f.urlResult, f.urlErr
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
