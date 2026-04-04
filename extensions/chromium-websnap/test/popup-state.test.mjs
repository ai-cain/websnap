import test from 'node:test';
import assert from 'node:assert/strict';

import { countTabs, findTabInSnapshot, resolveSelectedTab } from '../core/popup-state.mjs';

const snapshot = {
  windows: [
    {
      windowId: 7,
      tabs: [
        { tabId: 11, title: 'Docs', url: 'https://example.com/docs', active: false },
        { tabId: 12, title: 'Timeline', url: 'https://x.com/home', active: true }
      ]
    }
  ]
};

test('countTabs sums all tabs from every window', () => {
  assert.equal(countTabs(snapshot), 2);
});

test('findTabInSnapshot returns the exact selected tab', () => {
  assert.deepEqual(findTabInSnapshot(snapshot, 7, 11), {
    windowId: 7,
    tabId: 11,
    title: 'Docs',
    url: 'https://example.com/docs',
    active: false
  });
});

test('resolveSelectedTab falls back to the active tab when selection is stale', () => {
  assert.deepEqual(resolveSelectedTab(snapshot, 999, 888, {
    windowId: 7,
    tabId: 12,
    title: 'Timeline',
    url: 'https://x.com/home',
    active: true
  }), {
    windowId: 7,
    tabId: 12,
    title: 'Timeline',
    url: 'https://x.com/home',
    active: true
  });
});
