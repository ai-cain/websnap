import test from 'node:test';
import assert from 'node:assert/strict';

import { buildSnapshot, detectBrowserName, formatHost, getActiveTabFromSnapshot, isWebUrl } from '../core/browser.mjs';

test('detectBrowserName picks edge when edg/ is present', () => {
  assert.equal(detectBrowserName('Mozilla/5.0 Edg/123.0'), 'edge');
});

test('isWebUrl only accepts http and https', () => {
  assert.equal(isWebUrl('https://example.com'), true);
  assert.equal(isWebUrl('http://example.com'), true);
  assert.equal(isWebUrl('chrome://extensions'), false);
  assert.equal(isWebUrl('file:///C:/docs/test.html'), false);
});

test('buildSnapshot filters non-web tabs and keeps active tab title', () => {
  const snapshot = buildSnapshot([
    {
      id: 7,
      tabs: [
        { id: 1, index: 0, title: 'Extensions', url: 'chrome://extensions', active: false },
        { id: 2, index: 1, title: 'X timeline', url: 'https://x.com/home', active: true }
      ]
    }
  ], 'chrome');

  assert.equal(snapshot.windows.length, 1);
  assert.equal(snapshot.windows[0].title, 'X timeline');
  assert.equal(snapshot.windows[0].tabs.length, 1);
  assert.equal(snapshot.windows[0].tabs[0].tabId, 2);
});

test('getActiveTabFromSnapshot falls back to the first tab when needed', () => {
  const active = getActiveTabFromSnapshot({
    windows: [
      {
        windowId: 10,
        tabs: [
          { tabId: 55, title: 'Docs', url: 'https://example.com/docs', active: false }
        ]
      }
    ]
  });

  assert.deepEqual(active, {
    windowId: 10,
    tabId: 55,
    title: 'Docs',
    url: 'https://example.com/docs',
    active: false
  });
});

test('formatHost strips www from hostnames', () => {
  assert.equal(formatHost('https://www.example.com/path'), 'example.com');
});
