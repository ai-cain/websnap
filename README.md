# websnap

> Go CLI proposal for reproducible web UI screenshots from the terminal.

---

## Project status

| Field | Status |
| --- | --- |
| Current phase | Bootstrap / `v0.1.0` in progress |
| Published release | None |
| Immediate target | `v0.1.0` |
| Base stack | **Go + chromedp** |
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

---

## Current bootstrap commands

Already implemented in the current bootstrap:

```bash
websnap shot https://example.com
websnap shot https://example.com --width 1440 --height 900
websnap shot https://example.com --out ./captures/home.png
```

Planned next:

```bash
websnap shot https://example.com --selector ".hero"
websnap shot https://example.com --full-page
```

---

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
- the `shot` orchestration path
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

## Roadmap summary

- `v0.1.0` — CLI bootstrap + basic screenshot
- `v0.2.0` — selector and full-page capture
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
