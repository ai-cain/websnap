# websnap Chromium extension

Load this folder as an unpacked extension in a Chromium browser:

1. Open `chrome://extensions` (or `edge://extensions`)
2. Enable **Developer mode**
3. Click **Load unpacked**
4. Select this `extensions/chromium-websnap` folder

## Current state

This folder is now a **standalone extension MVP** for browser capture.

It no longer depends on the old local WebSocket bridge.

## What it does

- lists open `http(s)` tabs across Chromium windows
- captures visible tab content with no browser chrome around it
- downloads the PNG directly from the extension
- shows a popup menu with:
  - **Capture** — quick actions for active/selected tab
  - **Tabs** — browse and select any open tab
  - **Settings** — save behavior and extension notes

## Included files

- `manifest.json` — MV3 manifest
- `service_worker.js` — runtime logic for tab discovery, capture, and download
- `popup.html` / `popup.css` / `popup.js` — popup menu UI
- `core/` — shared browser/naming/popup helpers
- `test/` — lightweight Node tests for pure helper logic

## Quick local verification

```powershell
node --test extensions/chromium-websnap/test/*.test.mjs
```

## Next direction

The next obvious steps are polishing UX, adding settings/history, and preparing packaging/distribution.
