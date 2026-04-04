# websnap

`websnap` is now being reshaped into the browser-extension product.

## Current repository focus

This repository keeps the Chromium extension source for the web capture path.

Current extension folder:

```text
extensions/chromium-websnap
```

## What the extension does now

The extension is now a **standalone MV3 browser capture MVP**:

- lists open `http(s)` tabs across normal Chromium windows
- captures the **visible tab content** instead of the whole browser chrome
- downloads the PNG directly from the extension
- offers a popup menu with Capture / Tabs / Settings

## Local loading

In Chrome or Edge:

1. Open `chrome://extensions` or `edge://extensions`
2. Enable developer mode
3. Click **Load unpacked**
4. Select `extensions/chromium-websnap`
5. Open the extension popup and capture a tab

## Notes

- The Go CLI and desktop/live capture code were moved to the separate `desksnap` repository.
- This repo is now the web-facing side of the split.
