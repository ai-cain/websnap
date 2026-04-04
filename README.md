# websnap

Web-first screenshot capture is now being prepared as the browser-extension product.

## Current repository focus

This repository keeps the Chromium extension source for the web capture path.

Current extension folder:

```text
extensions/chromium-websnap
```

## Local loading

In Chrome or Edge:

1. Open `chrome://extensions` or `edge://extensions`
2. Enable developer mode
3. Click **Load unpacked**
4. Select `extensions/chromium-websnap`

## Notes

- The Go CLI and desktop/live capture code were moved to the separate `desksnap` repository.
- This repo is now the web-facing side of the split.
