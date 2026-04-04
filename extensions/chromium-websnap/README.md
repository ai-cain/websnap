# websnap Browser Bridge Extension

Load this folder as an unpacked extension in a Chromium browser:

1. Open `chrome://extensions` (or `edge://extensions`)
2. Enable **Developer mode**
3. Click **Load unpacked**
4. Select this `extensions/chromium-websnap` folder

When `websnap` is running, the extension connects to `ws://127.0.0.1:38971/ws`.

## What it does

- Streams open Chromium windows/tabs back to `websnap`
- Lets `websnap` capture **visible tab content** instead of the whole browser window
- Keeps apps/folders on the native Windows capture path

## Current scope

- Designed for Chromium browsers first
- Chrome and Edge are the primary targets
- Browser capture is extension-backed; desktop apps/folders stay native
