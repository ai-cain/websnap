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

	if model.inputs[fieldWidth].Value() != "1440" {
		t.Fatalf("width default = %q, want %q", model.inputs[fieldWidth].Value(), "1440")
	}

	if model.inputs[fieldHeight].Value() != "900" {
		t.Fatalf("height default = %q, want %q", model.inputs[fieldHeight].Value(), "900")
	}
}

func TestModelBuildRequest(t *testing.T) {
	t.Parallel()

	model := NewModel(&fakeShotRunner{})
	model.inputs[fieldURL].SetValue("https://example.com")
	model.inputs[fieldWidth].SetValue("1600")
	model.inputs[fieldHeight].SetValue("1000")
	model.inputs[fieldOut].SetValue("captures/home.png")

	req, err := model.buildRequest()
	if err != nil {
		t.Fatalf("buildRequest() error = %v", err)
	}

	expected := domain.CaptureRequest{
		URL:    "https://example.com",
		Width:  1600,
		Height: 1000,
		Out:    "captures/home.png",
	}

	if req != expected {
		t.Fatalf("request = %#v, want %#v", req, expected)
	}
}
