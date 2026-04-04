import { buildSnapshot, detectBrowserName, getActiveTabFromSnapshot } from './core/browser.mjs';
import { buildCaptureFilename } from './core/naming.mjs';

const DEFAULT_OPTIONS = {
  saveAs: false
};

chrome.runtime.onInstalled.addListener(async () => {
  const current = await storageGet(DEFAULT_OPTIONS);
  await storageSet({ ...DEFAULT_OPTIONS, ...current });
});

chrome.runtime.onMessage.addListener((message, _sender, sendResponse) => {
  handleMessage(message)
    .then((result) => sendResponse({ ok: true, ...result }))
    .catch((error) => sendResponse({ ok: false, error: error?.message || String(error) }));

  return true;
});

async function handleMessage(message) {
  switch (message?.type) {
    case 'getSnapshot': {
      return { snapshot: await createSnapshot() };
    }
    case 'getOptions': {
      return { options: await storageGet(DEFAULT_OPTIONS) };
    }
    case 'setOptions': {
      const nextOptions = { ...DEFAULT_OPTIONS, ...(message.options || {}) };
      await storageSet(nextOptions);
      return { options: nextOptions };
    }
    case 'captureActiveTab': {
      const snapshot = await createSnapshot();
      const active = getActiveTabFromSnapshot(snapshot);
      if (!active) {
        throw new Error('No active web tab is available to capture');
      }

      const options = await storageGet(DEFAULT_OPTIONS);
      return captureTab({
        windowId: active.windowId,
        tabId: active.tabId,
        saveAs: options.saveAs
      });
    }
    case 'captureTab': {
      const options = await storageGet(DEFAULT_OPTIONS);
      return captureTab({
        windowId: Number(message.windowId),
        tabId: Number(message.tabId),
        saveAs: typeof message.saveAs === 'boolean' ? message.saveAs : options.saveAs
      });
    }
    default:
      throw new Error(`Unsupported message type: ${message?.type || 'unknown'}`);
  }
}

async function createSnapshot() {
  const browser = detectBrowserName(self.navigator.userAgent);
  const windows = await chrome.windows.getAll({
    populate: true,
    windowTypes: ['normal']
  });

  return buildSnapshot(windows, browser);
}

async function captureTab({ windowId, tabId, saveAs }) {
  if (!windowId || !tabId) {
    throw new Error('windowId and tabId are required');
  }

  await chrome.windows.update(windowId, { focused: true });
  await chrome.tabs.update(tabId, { active: true });
  await delay(250);

  const tab = await chrome.tabs.get(tabId);
  const pngDataUrl = await chrome.tabs.captureVisibleTab(windowId, { format: 'png' });
  const filename = buildCaptureFilename({
    title: tab?.title || '',
    url: tab?.url || '',
    browserName: detectBrowserName(self.navigator.userAgent),
    now: new Date()
  });

  const downloadId = await chrome.downloads.download({
    url: pngDataUrl,
    filename,
    saveAs: Boolean(saveAs)
  });

  return {
    downloadId,
    filename,
    tab: {
      tabId,
      windowId,
      title: tab?.title || '',
      url: tab?.url || ''
    }
  };
}

function delay(ms) {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

function storageGet(defaults) {
  return new Promise((resolve, reject) => {
    chrome.storage.local.get(defaults, (items) => {
      const error = chrome.runtime.lastError;
      if (error) {
        reject(new Error(error.message));
        return;
      }
      resolve(items);
    });
  });
}

function storageSet(values) {
  return new Promise((resolve, reject) => {
    chrome.storage.local.set(values, () => {
      const error = chrome.runtime.lastError;
      if (error) {
        reject(new Error(error.message));
        return;
      }
      resolve();
    });
  });
}
