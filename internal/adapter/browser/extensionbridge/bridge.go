package extensionbridge

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"image/png"
	"net"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ai-cain/websnap/internal/domain"
	apperrors "github.com/ai-cain/websnap/internal/support/errors"
	"github.com/gorilla/websocket"
)

const defaultAddress = "127.0.0.1:38971"

type Bridge struct {
	addr string

	mu       sync.RWMutex
	clients  map[string]*client
	server   *http.Server
	listener net.Listener
	started  bool
	nextID   atomic.Uint64
}

type bridgeMessage struct {
	ID         string                 `json:"id,omitempty"`
	Type       string                 `json:"type"`
	Browser    string                 `json:"browser,omitempty"`
	WindowID   int                    `json:"windowId,omitempty"`
	TabID      int                    `json:"tabId,omitempty"`
	Windows    []browserWindowPayload `json:"windows,omitempty"`
	PNGDataURL string                 `json:"pngDataUrl,omitempty"`
	Error      string                 `json:"error,omitempty"`
}

type browserWindowPayload struct {
	WindowID int                 `json:"windowId"`
	AppName  string              `json:"appName"`
	Title    string              `json:"title"`
	Tabs     []browserTabPayload `json:"tabs"`
}

type browserTabPayload struct {
	TabID  int    `json:"tabId"`
	Index  int    `json:"index"`
	Title  string `json:"title"`
	URL    string `json:"url"`
	Active bool   `json:"active"`
}

type browserSnapshot struct {
	Browser string
	Windows []browserWindowPayload
}

type captureResponse struct {
	pngDataURL string
	err        error
}

type client struct {
	bridge    *Bridge
	conn      *websocket.Conn
	writeMu   sync.Mutex
	pendingMu sync.Mutex
	pending   map[string]chan captureResponse

	mu       sync.RWMutex
	browser  string
	snapshot browserSnapshot
}

func New(addr string) *Bridge {
	if strings.TrimSpace(addr) == "" {
		addr = defaultAddress
	}

	return &Bridge{
		addr:    addr,
		clients: map[string]*client{},
	}
}

func (b *Bridge) Start() error {
	if b == nil {
		return apperrors.New(apperrors.CodeInvalidArgument, "browser extension bridge is not configured")
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	if b.started {
		return nil
	}

	listener, err := net.Listen("tcp", b.addr)
	if err != nil {
		return err
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/ws", b.handleWebsocket)

	server := &http.Server{
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	b.listener = listener
	b.server = server
	b.started = true

	go func() {
		_ = server.Serve(listener)
	}()

	return nil
}

func (b *Bridge) Close(ctx context.Context) error {
	if b == nil {
		return nil
	}

	b.mu.Lock()
	server := b.server
	listener := b.listener
	clients := make([]*client, 0, len(b.clients))
	for _, current := range b.clients {
		clients = append(clients, current)
	}
	b.clients = map[string]*client{}
	b.server = nil
	b.listener = nil
	b.started = false
	b.mu.Unlock()

	for _, current := range clients {
		current.failPending(fmt.Errorf("browser extension disconnected"))
		_ = current.conn.Close()
	}

	if server != nil {
		if ctx == nil {
			ctx = context.Background()
		}
		return server.Shutdown(ctx)
	}

	if listener != nil {
		return listener.Close()
	}

	return nil
}

func (b *Bridge) Address() string {
	if b == nil {
		return ""
	}

	b.mu.RLock()
	defer b.mu.RUnlock()

	if b.listener != nil {
		return b.listener.Addr().String()
	}

	return b.addr
}

func (b *Bridge) ListTargets(_ context.Context) ([]domain.LiveTarget, error) {
	if b == nil {
		return nil, apperrors.New(apperrors.CodeInvalidArgument, "browser extension bridge is not configured")
	}

	b.mu.RLock()
	clients := make([]*client, 0, len(b.clients))
	for _, current := range b.clients {
		clients = append(clients, current)
	}
	b.mu.RUnlock()

	targets := make([]domain.LiveTarget, 0)
	for _, current := range clients {
		current.mu.RLock()
		snapshot := current.snapshot
		current.mu.RUnlock()
		targets = append(targets, buildTargetsFromSnapshot(snapshot)...)
	}

	return targets, nil
}

func (b *Bridge) ListTabs(_ context.Context, target domain.LiveTarget) ([]domain.BrowserTab, error) {
	if b == nil {
		return nil, apperrors.New(apperrors.CodeInvalidArgument, "browser extension bridge is not configured")
	}

	client := b.clientForBrowser(target.AppName)
	if client == nil {
		return nil, apperrors.New(apperrors.CodeCaptureFailed, "browser extension is not connected for the selected browser")
	}

	windowID := target.BrowserWindowID
	if windowID <= 0 {
		windowID = int(target.WindowHandle)
	}

	client.mu.RLock()
	defer client.mu.RUnlock()

	for _, window := range client.snapshot.Windows {
		if window.WindowID != windowID {
			continue
		}

		tabs := make([]domain.BrowserTab, 0, len(window.Tabs))
		for _, tab := range window.Tabs {
			tabs = append(tabs, domain.BrowserTab{
				Index:    tab.Index,
				ID:       tab.TabID,
				WindowID: window.WindowID,
				URL:      tab.URL,
				Title:    tab.Title,
				Selected: tab.Active,
			})
		}
		return tabs, nil
	}

	return nil, apperrors.New(apperrors.CodeCaptureFailed, "browser window is not available in the extension snapshot")
}

func (b *Bridge) Capture(ctx context.Context, req domain.LiveCaptureRequest) (domain.LiveCaptureImage, error) {
	if b == nil {
		return domain.LiveCaptureImage{}, apperrors.New(apperrors.CodeInvalidArgument, "browser extension bridge is not configured")
	}

	client := b.clientForBrowser(req.Target.AppName)
	if client == nil {
		return domain.LiveCaptureImage{}, apperrors.New(apperrors.CodeCaptureFailed, "browser extension is not connected for the selected browser")
	}

	windowID := req.Target.BrowserWindowID
	if windowID <= 0 {
		windowID = int(req.Target.WindowHandle)
	}

	tabID := req.TabID
	if tabID <= 0 {
		if activeTabID := client.activeTabID(windowID); activeTabID > 0 {
			tabID = activeTabID
		}
	}

	if tabID <= 0 {
		return domain.LiveCaptureImage{}, apperrors.New(apperrors.CodeInvalidArgument, "browser tab id is required for extension capture")
	}

	requestID := fmt.Sprintf("capture-%d", b.nextID.Add(1))
	replyCh := make(chan captureResponse, 1)

	client.registerPending(requestID, replyCh)
	defer client.unregisterPending(requestID)

	if err := client.writeJSON(bridgeMessage{
		ID:       requestID,
		Type:     "capture",
		Browser:  normalizedBrowserName(req.Target.AppName),
		WindowID: windowID,
		TabID:    tabID,
	}); err != nil {
		return domain.LiveCaptureImage{}, apperrors.Wrap(apperrors.CodeCaptureFailed, "failed to request browser tab capture", err)
	}

	select {
	case <-ctx.Done():
		return domain.LiveCaptureImage{}, ctx.Err()
	case response := <-replyCh:
		if response.err != nil {
			return domain.LiveCaptureImage{}, apperrors.Wrap(apperrors.CodeCaptureFailed, "browser extension capture failed", response.err)
		}

		pngBytes, width, height, err := decodePNGDataURL(response.pngDataURL)
		if err != nil {
			return domain.LiveCaptureImage{}, apperrors.Wrap(apperrors.CodeCaptureFailed, "browser extension returned an invalid PNG payload", err)
		}

		return domain.LiveCaptureImage{
			PNG:    pngBytes,
			Width:  width,
			Height: height,
		}, nil
	}
}

func (b *Bridge) handleWebsocket(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(_ *http.Request) bool {
			return true
		},
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	current := &client{
		bridge:  b,
		conn:    conn,
		pending: map[string]chan captureResponse{},
	}

	current.readLoop()
}

func (b *Bridge) registerClient(browser string, current *client) {
	if current == nil {
		return
	}

	browser = normalizedBrowserName(browser)
	if browser == "" {
		return
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	if previous, ok := b.clients[browser]; ok && previous != current {
		previous.failPending(fmt.Errorf("browser extension connection was replaced"))
		_ = previous.conn.Close()
	}

	b.clients[browser] = current
}

func (b *Bridge) unregisterClient(current *client) {
	if current == nil {
		return
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	if current.browser == "" {
		return
	}

	browser := normalizedBrowserName(current.browser)
	if registered, ok := b.clients[browser]; ok && registered == current {
		delete(b.clients, browser)
	}
}

func (b *Bridge) clientForBrowser(appName string) *client {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.clients[normalizedBrowserName(appName)]
}

func (c *client) readLoop() {
	defer func() {
		c.bridge.unregisterClient(c)
		c.failPending(fmt.Errorf("browser extension disconnected"))
		_ = c.conn.Close()
	}()

	for {
		var message bridgeMessage
		if err := c.conn.ReadJSON(&message); err != nil {
			return
		}

		switch message.Type {
		case "hello":
			browser := normalizedBrowserName(message.Browser)
			c.mu.Lock()
			c.browser = browser
			c.snapshot.Browser = browser
			c.mu.Unlock()
			c.bridge.registerClient(browser, c)
		case "snapshot":
			browser := normalizedBrowserName(message.Browser)
			c.mu.Lock()
			if browser == "" {
				browser = c.browser
			}
			c.snapshot = browserSnapshot{
				Browser: browser,
				Windows: message.Windows,
			}
			c.mu.Unlock()
		case "capture-result":
			c.resolvePending(message.ID, captureResponse{
				pngDataURL: message.PNGDataURL,
				err:        responseError(message.Error),
			})
		}
	}
}

func (c *client) writeJSON(message bridgeMessage) error {
	c.writeMu.Lock()
	defer c.writeMu.Unlock()
	return c.conn.WriteJSON(message)
}

func (c *client) registerPending(id string, ch chan captureResponse) {
	c.pendingMu.Lock()
	defer c.pendingMu.Unlock()
	c.pending[id] = ch
}

func (c *client) unregisterPending(id string) {
	c.pendingMu.Lock()
	defer c.pendingMu.Unlock()
	delete(c.pending, id)
}

func (c *client) resolvePending(id string, response captureResponse) {
	c.pendingMu.Lock()
	ch := c.pending[id]
	delete(c.pending, id)
	c.pendingMu.Unlock()

	if ch == nil {
		return
	}

	select {
	case ch <- response:
	default:
	}
}

func (c *client) failPending(err error) {
	c.pendingMu.Lock()
	pending := c.pending
	c.pending = map[string]chan captureResponse{}
	c.pendingMu.Unlock()

	for _, ch := range pending {
		select {
		case ch <- captureResponse{err: err}:
		default:
		}
	}
}

func (c *client) activeTabID(windowID int) int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	for _, window := range c.snapshot.Windows {
		if window.WindowID != windowID {
			continue
		}
		for _, tab := range window.Tabs {
			if tab.Active {
				return tab.TabID
			}
		}
	}

	return 0
}

func buildTargetsFromSnapshot(snapshot browserSnapshot) []domain.LiveTarget {
	browser := normalizedBrowserName(snapshot.Browser)
	targets := make([]domain.LiveTarget, 0, len(snapshot.Windows))
	for _, window := range snapshot.Windows {
		title := strings.TrimSpace(window.Title)
		if title == "" {
			title = activeTabTitle(window.Tabs)
		}
		if title == "" {
			title = "Browser Window"
			if browser != "" {
				title = strings.ToUpper(browser[:1]) + browser[1:] + " Window"
			}
		}

		targets = append(targets, domain.LiveTarget{
			WindowHandle:    int64(window.WindowID),
			Title:           title,
			AppName:         firstNonEmpty(window.AppName, browser),
			Type:            domain.LiveTargetBrowser,
			CanListTabs:     true,
			Provider:        domain.LiveTargetProviderBrowserExtension,
			BrowserWindowID: window.WindowID,
		})
	}

	return targets
}

func activeTabTitle(tabs []browserTabPayload) string {
	for _, tab := range tabs {
		if tab.Active && strings.TrimSpace(tab.Title) != "" {
			return tab.Title
		}
	}
	for _, tab := range tabs {
		if strings.TrimSpace(tab.Title) != "" {
			return tab.Title
		}
	}
	return ""
}

func decodePNGDataURL(value string) ([]byte, int64, int64, error) {
	if strings.TrimSpace(value) == "" {
		return nil, 0, 0, fmt.Errorf("png payload is empty")
	}

	payload := value
	if strings.HasPrefix(payload, "data:") {
		comma := strings.Index(payload, ",")
		if comma < 0 {
			return nil, 0, 0, fmt.Errorf("png data url is invalid")
		}
		payload = payload[comma+1:]
	}

	pngBytes, err := base64.StdEncoding.DecodeString(payload)
	if err != nil {
		return nil, 0, 0, err
	}

	config, err := png.DecodeConfig(bytes.NewReader(pngBytes))
	if err != nil {
		return nil, 0, 0, err
	}

	return pngBytes, int64(config.Width), int64(config.Height), nil
}

func normalizedBrowserName(value string) string {
	value = strings.TrimSpace(strings.ToLower(value))
	switch value {
	case "microsoft-edge", "edge":
		return "edge"
	case "google-chrome", "chrome":
		return "chrome"
	case "brave":
		return "brave"
	case "opera":
		return "opera"
	default:
		return value
	}
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value != "" {
			return value
		}
	}
	return ""
}

func responseError(value string) error {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	return fmt.Errorf("%s", value)
}
