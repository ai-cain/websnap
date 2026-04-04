# websnap

> Go CLI proposal for reproducible web UI screenshots from the terminal.

---

## Project status

| Field | Status |
| --- | --- |
| Current phase | Capture targets complete / ready for `v0.3.0` |
| Published release | None |
| Immediate target | `v0.3.0` |
| Base stack | **Go + chromedp + Bubble Tea** |
| Documentation language | English-first |
| CLI i18n | Deferred to an advanced version |

> **Important:** this repository now contains an initial Go bootstrap for `shot`, but it is **not a published release yet**.  
> The README documents both the implemented base and the planned direction for the remaining roadmap.

---

## What problem it solves

`websnap` is meant to solve one focused problem well: **capture web UIs in a reproducible, scriptable, terminal-friendly way**.

Typical use cases:

- documenting pages and components
- generating visual evidence for reviews
- capturing local (`localhost`) and remote URLs
- preparing assets for demos, PRs, and portfolio material
- automating screenshots from scripts or CI workflows

---

## The key question: how can a terminal tool take a screenshot?

The terminal does **not** render the web.  
The CLI simply **orchestrates a headless browser**.

Proposed flow:

1. The user runs `websnap shot <url>`.
2. The CLI parses arguments and validates the request.
3. The CLI starts a headless Chromium instance.
4. Chromium loads the page and renders it off-screen.
5. The CLI tells the browser what to capture:
   - current viewport
   - full page
   - or a specific element
6. The CLI writes the PNG to disk and returns the resulting path.

In other words: **the terminal does not take the picture**; the CLI **directs** a headless browser to do it.

---

## First release goal

The first useful release should stay small but serious: **stable screenshot capture**.

### Committed scope for `v0.1.0`

- styled interactive terminal mode
- `shot` command
- URL-based capture
- configurable viewport via `--width` and `--height`
- explicit output path via `--out`
- automatic `media/img` creation
- clear error messages
- simple, defensible CLI contract

### Out of scope for `v0.1.0`

- GIF
- video
- watch mode
- third-party uploads
- config file
- advanced authentication
- test-runner-like automation

## Why Go

Go is the proposed choice for product reasons, not trend chasing:

- single distributable binary
- excellent fit for CLI tooling
- low operational friction in CI
- fast startup and simple maintenance
- strong base for growth without dragging unnecessary tooling into the project

Initial technical direction:

- **Language:** Go
- **Browser engine:** `chromedp`
- **Interactive terminal UI:** Bubble Tea + Bubbles + Lip Gloss
- **Future GIF processing:** FFmpeg, but outside `v0.1.0`

---

## Proposed output structure

Artifacts should be written to the **current execution directory**, not to the repository itself:

```text
media/
  └── img/
```

Later versions may extend that to:

```text
media/
  ├── img/
  └── gif/
```

---

## Current repository reality

At the moment, this repository contains:

- an initial Go CLI bootstrap
- a styled interactive TUI flow
- the `shot` orchestration path
- live capture for open apps and folders on Windows
- an extension-backed browser bridge for Chromium browsers
- `chromedp` browser integration
- filesystem output persistence
- domain and CLI tests
- versioned documentation and architecture notes

Still to be built:

- selector capture
- full-page capture
- advanced diagnostics around browser prerequisites
- GIF pipeline
- release packaging and publishing

---

## Documentation map

- [`docs/README.md`](docs/README.md) — documentation index
- [`docs/FEATURES.md`](docs/FEATURES.md) — versioned feature roadmap
- [`docs/ARCHITECTURE.md`](docs/ARCHITECTURE.md) — proposed Go architecture

---

## Installation and update

### Requirements

- Go `1.26+`
- Chrome or Chromium available on the machine

### Option A — run from the repository without installing

Useful while the project is still in bootstrap mode:

```bash
go run ./cmd/websnap
go run ./cmd/websnap shot https://example.com
```

### Option B — install from local source

From the repository root:

```bash
go install ./cmd/websnap
```

After that, ensure Go's binary directory is on your `PATH`.

Windows check:

```powershell
go env GOPATH
```

Then make sure this folder is reachable from the terminal:

```text
%GOPATH%\bin
```

If it is not in `PATH`, PowerShell will show:

```text
websnap : The term 'websnap' is not recognized ...
```

### Option C — install directly from GitHub

For public usage from the repository path:

```bash
go install github.com/ai-cain/websnap/cmd/websnap@latest
```

### Update to a newer version

If you installed through `go install`, update by reinstalling the target version:

```bash
go install github.com/ai-cain/websnap/cmd/websnap@latest
```

When tagged releases exist, install a fixed version explicitly:

```bash
go install github.com/ai-cain/websnap/cmd/websnap@v0.1.0
```

### Verify the Go toolchain correctly

Use:

```bash
go version
```

Not:

```bash
go --version
```

Go uses subcommands, not GNU-style long flags for version output.

---

## Current bootstrap commands

### Interactive mode

Launch the interactive terminal flow:

```bash
websnap
websnap interactive
```

The current TUI provides:

- open app / folder selection
- browser-window selection
- browser-tab selection
- output path (optional)
- keyboard-first navigation
- styled panels, focus states, and capture feedback

### Browser-only page capture

For Chromium browsers, `websnap` can use the unpacked extension in [`extensions/chromium-websnap`](extensions/chromium-websnap) so browser targets capture the **visible tab content** instead of the whole browser window chrome.

Load it with:

1. `chrome://extensions` or `edge://extensions`
2. enable **Developer mode**
3. **Load unpacked**
4. choose `extensions/chromium-websnap`

When the extension is loaded and `websnap` is running, press `R` in the interactive UI to refresh browser targets from the extension bridge.

### Scripted mode

```bash
websnap shot https://example.com
websnap shot https://example.com --width 1440 --height 900
websnap shot https://example.com --selector "#app"
websnap shot https://example.com --full-page
websnap shot https://example.com --out ./captures/home.png
```

Current rule:

- `--selector` and `--full-page` are **mutually exclusive**

---

## Roadmap summary

- `v0.1.0` — interactive CLI + basic screenshot ✅ completed as local milestone
- `v0.2.0` — selector and full-page capture ✅ completed as local milestone
- `v0.3.0` — reproducibility and developer experience
- `v0.4.0` — clip support and capture refinements
- `v0.5.0` — experimental GIF pipeline
- `v1.0.0` — stable release
- `v1.1.0` — CLI internationalization (`en` / `es`)

Detailed version planning lives in [`docs/FEATURES.md`](docs/FEATURES.md).

---

## Design principles

- **reproducibility over magic**
- **small V1, solid foundations**
- **simple CLI contract**
- **keep screenshots and GIF on separate evolution tracks**
- **extensible architecture without premature complexity**

---

## License

Still pending formal definition in a root `LICENSE` file.
