# websnap Chromium extension scaffold

Load this folder as an unpacked extension in a Chromium browser:

1. Open `chrome://extensions` (or `edge://extensions`)
2. Enable **Developer mode**
3. Click **Load unpacked**
4. Select this `extensions/chromium-websnap` folder

## Current state

This folder is the preserved browser-extension scaffold from the original split.

It is **not a standalone finished product yet**.
The current service worker still expects a local WebSocket bridge at:

```text
ws://127.0.0.1:38971/ws
```

That bridge no longer lives in this repository after the CLI/desktop code moved to `desksnap`.

## What remains here

- Chromium window/tab discovery prototype
- visible-tab capture prototype
- base manifest and service worker structure

## Next direction

To turn this into the real `websnap` product, the extension still needs its own standalone UX and distribution flow, instead of depending on the old local bridge.
