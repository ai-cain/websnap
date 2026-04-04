package router

import (
	"context"

	"github.com/ai-cain/websnap/internal/domain"
	"github.com/ai-cain/websnap/internal/port"
)

type Catalog struct {
	desktop port.LiveTargetCatalog
	web     port.LiveTargetCatalog
}

func NewCatalog(desktop, web port.LiveTargetCatalog) *Catalog {
	return &Catalog{
		desktop: desktop,
		web:     web,
	}
}

func (c *Catalog) ListTargets(ctx context.Context) ([]domain.LiveTarget, error) {
	desktopTargets, err := c.desktop.ListTargets(ctx)
	if err != nil {
		return nil, err
	}

	if c.web == nil {
		return desktopTargets, nil
	}

	webTargets, err := c.web.ListTargets(ctx)
	if err != nil || len(webTargets) == 0 {
		return desktopTargets, nil
	}

	filtered := make([]domain.LiveTarget, 0, len(desktopTargets)+len(webTargets))
	for _, target := range desktopTargets {
		if target.Type == domain.LiveTargetBrowser {
			continue
		}
		filtered = append(filtered, target)
	}

	filtered = append(filtered, webTargets...)
	return filtered, nil
}

func (c *Catalog) ListTabs(ctx context.Context, target domain.LiveTarget) ([]domain.BrowserTab, error) {
	if target.Provider == domain.LiveTargetProviderBrowserExtension && c.web != nil {
		return c.web.ListTabs(ctx, target)
	}

	return c.desktop.ListTabs(ctx, target)
}

type Capturer struct {
	desktop port.LiveCapturer
	web     port.LiveCapturer
}

func NewCapturer(desktop, web port.LiveCapturer) *Capturer {
	return &Capturer{
		desktop: desktop,
		web:     web,
	}
}

func (c *Capturer) Capture(ctx context.Context, req domain.LiveCaptureRequest) (domain.LiveCaptureImage, error) {
	if req.Target.Provider == domain.LiveTargetProviderBrowserExtension && c.web != nil {
		return c.web.Capture(ctx, req)
	}

	return c.desktop.Capture(ctx, req)
}
