# Feature roadmap by version

This document separates two things that teams often mix badly:

1. **what is committed**
2. **what is desirable**

A serious proposal does not promise everything at once. It versions growth.

---

## Current status

- Repository state: **proposal / pre-alpha**
- Published version: **none**
- Immediate target: **`v0.1.0`**

---

## Roadmap conventions

- **In progress**: implemented partially and actively being completed
- **Committed**: expected to land in that version unless scope changes materially
- **Candidate**: valid direction, but not committed yet
- **Deferred**: intentionally postponed

---

## `v0.1.0` — Useful CLI bootstrap

**Goal:** deliver a first executable version that is small, clear, and defensible.

**Status:** In progress

### Implemented now

- interactive terminal mode
- `websnap shot <url>`
- URL validation
- configurable viewport:
  - `--width`
  - `--height`
- explicit output path with `--out`
- automatic `media/img` creation
- clear terminal-facing errors
- minimum viable CLI output contract

### Remaining to close `v0.1.0`

- browser prerequisite diagnostics should be more explicit
- release packaging and publishable install story still need to be closed
- command help and docs should keep tightening around the real bootstrap behavior

### Still not included

- `--selector`
- `--full-page`
- `--clip`
- GIF
- video
- watch mode
- file-based configuration

---

## `v0.2.0` — Capture targets

**Goal:** move from “basic capture” to “useful UI capture”.

**Status:** Committed

### Scope

- `--selector`
- `--full-page`
- consistent file naming strategy
- basic protection against invalid paths
- more descriptive selector-related failures

### Main risk

Element capture introduces a stronger dependency on DOM timing and render stability. It should be solved without making the CLI opaque.

---

## `v0.3.0` — Reproducibility and developer experience

**Goal:** make the tool more reliable for demos, documentation, and automation.

**Status:** Committed

### Scope

- `--delay`
- `--timeout`
- more explicit console output
- predictable exit codes
- better ergonomics for `localhost`

### Value

This version makes the tool defensible in real work, not just in happy-path demos.

---

## `v0.4.0` — Advanced capture control

**Goal:** add more precise control over the captured area.

**Status:** Candidate

### Candidate scope

- `--clip x,y,width,height`
- geometric validation for the requested area
- simple viewport presets

### Why it is not committed yet

Precise clipping looks simple, but it brings visual validation rules and consistency concerns that should come only after the main path is stable.

---

## `v0.5.0` — Experimental GIF pipeline

**Goal:** open the door to motion capture without pretending it is stable too early.

**Status:** Candidate

### Candidate scope

- `websnap gif <url>`
- `--duration`
- `--fps`
- sequential frame capture
- FFmpeg integration

### Why it lives here

GIF is **not** a minor extension of screenshot capture. It is a separate pipeline:

- multiple frames
- temporal coordination
- encoding
- performance cost
- extra runtime dependency

That is why it should arrive as an explicit experimental track.

---

## `v1.0.0` — Stable release

**Goal:** establish `websnap` as a reliable screenshot tool from the terminal.

**Status:** Committed as direction, exact scope still to be refined

### Expected outcome

- stable CLI contract
- solid screenshot path
- mature documentation
- predictable error model
- clear binary distribution story

### Important note

`v1.0.0` does not need GIF to be considered successful. Product stability can focus on screenshots first.

---

## `v1.1.0` — CLI internationalization

**Goal:** keep English as the base language while enabling Spanish as a supported user-facing locale.

**Status:** Committed after core stabilization

### Scope

- localized CLI help output
- localized error and feedback messages
- message catalog structure for `en` and `es`
- fallback to English when a translation is missing

### Architectural rule

Localization belongs in the presentation layer.  
Use cases and domain errors should expose codes or typed failures, not hard-coded prose.

---

## Post-`v1.1.0` backlog

**Status:** Deferred

Valid ideas, but outside the current commitment:

- video (`mp4` / `webm`)
- `.websnap.yaml`
- watch mode
- uploads to S3 / Cloudinary
- project-level presets
- advanced authentication

---

## Evolution rule

If a feature adds runtime complexity, new dependencies, or a distinct execution pipeline, it does **not** enter because it sounds exciting.  
The main path must stay protected first: `shot`.
