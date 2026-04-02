package tui

import (
	"testing"

	"github.com/ai-cain/websnap/internal/domain"
)

func TestNewModelDefaults(t *testing.T) {
	t.Parallel()

	model := NewModel(&fakeShotRunner{})

	if model.focus != fieldURL {
		t.Fatalf("focus = %d, want %d", model.focus, fieldURL)
	}

	if model.inputs[inputIndex(fieldWidth)].Value() != "1440" {
		t.Fatalf("width default = %q, want %q", model.inputs[inputIndex(fieldWidth)].Value(), "1440")
	}

	if model.inputs[inputIndex(fieldHeight)].Value() != "900" {
		t.Fatalf("height default = %q, want %q", model.inputs[inputIndex(fieldHeight)].Value(), "900")
	}

	if model.mode != modeViewport {
		t.Fatalf("mode = %v, want %v", model.mode, modeViewport)
	}
}

func TestModelBuildRequest(t *testing.T) {
	t.Parallel()

	model := NewModel(&fakeShotRunner{})
	model.inputs[inputIndex(fieldURL)].SetValue("https://example.com")
	model.inputs[inputIndex(fieldWidth)].SetValue("1600")
	model.inputs[inputIndex(fieldHeight)].SetValue("1000")
	model.inputs[inputIndex(fieldSelector)].SetValue("#app")
	model.inputs[inputIndex(fieldOut)].SetValue("captures/home.png")
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

func TestModelBuildRequestForFullPage(t *testing.T) {
	t.Parallel()

	model := NewModel(&fakeShotRunner{})
	model.inputs[inputIndex(fieldURL)].SetValue("https://example.com")
	model.inputs[inputIndex(fieldWidth)].SetValue("1600")
	model.inputs[inputIndex(fieldHeight)].SetValue("1000")
	model.mode = modeFullPage

	req, err := model.buildRequest()
	if err != nil {
		t.Fatalf("buildRequest() error = %v", err)
	}

	if !req.FullPage {
		t.Fatal("request should use FullPage=true")
	}
}
