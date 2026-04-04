import test from 'node:test';
import assert from 'node:assert/strict';

import { buildCaptureFilename, preferredCaptureStem, sanitizeFileStem, stripBrowserTitleSuffix } from '../core/naming.mjs';

test('stripBrowserTitleSuffix removes browser suffixes', () => {
  assert.equal(stripBrowserTitleSuffix('Your Repositories - Google Chrome'), 'Your Repositories');
});

test('sanitizeFileStem normalizes spaces and symbols', () => {
  assert.equal(sanitizeFileStem('Hello, World!'), 'hello-world');
});

test('preferredCaptureStem prefers meaningful title segment over map coordinates', () => {
  const stem = preferredCaptureStem({
    title: `12°02'15.0"S 76°57'45.7"W - Google Maps`,
    url: 'https://www.google.com/maps/place/Lima',
    browserName: 'chrome'
  });

  assert.equal(stem, 'google-maps');
});

test('buildCaptureFilename returns timestamped png path', () => {
  const filename = buildCaptureFilename({
    title: 'Your Repositories - Google Chrome',
    url: 'https://github.com/ai-cain',
    browserName: 'chrome',
    now: new Date(2026, 3, 4, 18, 20, 30)
  });

  assert.equal(filename, 'websnap/20260404-182030-your-repositories.png');
});
