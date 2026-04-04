import { formatHost, getActiveTabFromSnapshot } from './core/browser.mjs';

const state = {
  snapshot: null,
  selectedWindowId: 0,
  selectedTabId: 0,
  options: {
    saveAs: false
  },
  busy: false
};

const elements = {
  refreshButton: document.querySelector('#refreshButton'),
  captureActiveButton: document.querySelector('#captureActiveButton'),
  saveAsCheckbox: document.querySelector('#saveAsCheckbox'),
  statusText: document.querySelector('#statusText'),
  windowsContainer: document.querySelector('#windowsContainer'),
  browserBadge: document.querySelector('#browserBadge')
};

bootstrap().catch((error) => setStatus(error.message || String(error), 'error'));

async function bootstrap() {
  bindEvents();
  await refreshState();
}

function bindEvents() {
  elements.refreshButton.addEventListener('click', () => void refreshState());
  elements.captureActiveButton.addEventListener('click', () => void captureActiveTab());
  elements.saveAsCheckbox.addEventListener('change', () => void updateOptions());
}

async function refreshState() {
  setBusy(true);
  setStatus('Loading browser tabs…');

  try {
    const [{ snapshot }, { options }] = await Promise.all([
      sendMessage({ type: 'getSnapshot' }),
      sendMessage({ type: 'getOptions' })
    ]);

    state.snapshot = snapshot;
    state.options = options;
    elements.saveAsCheckbox.checked = Boolean(options?.saveAs);
    elements.browserBadge.textContent = titleCase(snapshot?.browser || 'chromium');

    primeSelection(snapshot);
    renderWindows();

    if ((snapshot?.windows || []).length === 0) {
      setStatus('No http(s) tabs are open right now. Open a normal web page and refresh.');
    } else {
      setStatus('Choose any tab below or capture the active one.', 'success');
    }
  } catch (error) {
    setStatus(error.message || String(error), 'error');
    elements.windowsContainer.innerHTML = renderEmptyState('Unable to load tabs from the extension.');
  } finally {
    setBusy(false);
  }
}

async function updateOptions() {
  try {
    const { options } = await sendMessage({
      type: 'setOptions',
      options: {
        saveAs: elements.saveAsCheckbox.checked
      }
    });
    state.options = options;
  } catch (error) {
    setStatus(error.message || String(error), 'error');
  }
}

async function captureActiveTab() {
  setBusy(true);
  setStatus('Capturing active tab…');

  try {
    const response = await sendMessage({ type: 'captureActiveTab' });
    setStatus(`Saved ${response.filename}`, 'success');
    await refreshState();
  } catch (error) {
    setStatus(error.message || String(error), 'error');
  } finally {
    setBusy(false);
  }
}

async function captureSelectedTab(windowId, tabId) {
  setBusy(true);
  setStatus('Capturing selected tab…');

  try {
    const response = await sendMessage({
      type: 'captureTab',
      windowId,
      tabId,
      saveAs: elements.saveAsCheckbox.checked
    });
    setStatus(`Saved ${response.filename}`, 'success');
    state.selectedWindowId = windowId;
    state.selectedTabId = tabId;
    await refreshState();
  } catch (error) {
    setStatus(error.message || String(error), 'error');
  } finally {
    setBusy(false);
  }
}

function renderWindows() {
  const windows = state.snapshot?.windows || [];
  if (windows.length === 0) {
    elements.windowsContainer.innerHTML = renderEmptyState('Open a normal web page and press refresh.');
    return;
  }

  elements.windowsContainer.innerHTML = windows.map((currentWindow) => renderWindowCard(currentWindow)).join('');
  bindDynamicEvents();
}

function bindDynamicEvents() {
  document.querySelectorAll('[data-select-tab]').forEach((button) => {
    button.addEventListener('click', () => {
      state.selectedWindowId = Number(button.dataset.windowId);
      state.selectedTabId = Number(button.dataset.tabId);
      renderWindows();
    });
  });

  document.querySelectorAll('[data-capture-tab]').forEach((button) => {
    button.addEventListener('click', () => {
      void captureSelectedTab(Number(button.dataset.windowId), Number(button.dataset.tabId));
    });
  });
}

function renderWindowCard(currentWindow) {
  const subtitle = `${currentWindow.tabs.length} ${currentWindow.tabs.length === 1 ? 'tab' : 'tabs'}`;
  const tabRows = currentWindow.tabs.map((tab) => renderTabRow(currentWindow.windowId, tab)).join('');

  return `
    <article class="window-card">
      <div class="window-header">
        <div>
          <p class="window-title">${escapeHtml(currentWindow.title)}</p>
          <p class="window-subtitle">${escapeHtml(subtitle)}</p>
        </div>
        <span class="badge">Window ${currentWindow.windowId}</span>
      </div>
      <div class="tab-list">${tabRows}</div>
    </article>
  `;
}

function renderTabRow(windowId, tab) {
  const selected = state.selectedWindowId === windowId && state.selectedTabId === tab.tabId;
  const host = formatHost(tab.url);
  const activePill = tab.active ? '<span class="active-pill">Active</span>' : '';

  return `
    <div class="tab-row ${selected ? 'selected' : ''}">
      <button class="tab-main ghost-button" type="button" data-select-tab="true" data-window-id="${windowId}" data-tab-id="${tab.tabId}">
        <div class="tab-title">${escapeHtml(tab.title)}</div>
        <div class="tab-url">${escapeHtml(host || tab.url)}</div>
      </button>
      <div class="tab-meta">
        ${activePill}
        <button type="button" data-capture-tab="true" data-window-id="${windowId}" data-tab-id="${tab.tabId}">Capture</button>
      </div>
    </div>
  `;
}

function renderEmptyState(message) {
  return `<div class="empty-state">${escapeHtml(message)}</div>`;
}

function primeSelection(snapshot) {
  const active = getActiveTabFromSnapshot(snapshot);
  if (!active) {
    state.selectedWindowId = 0;
    state.selectedTabId = 0;
    return;
  }

  const selectionStillExists = snapshot?.windows?.some((currentWindow) =>
    currentWindow.windowId === state.selectedWindowId &&
    currentWindow.tabs.some((tab) => tab.tabId === state.selectedTabId)
  );

  if (!selectionStillExists) {
    state.selectedWindowId = active.windowId;
    state.selectedTabId = active.tabId;
  }
}

function setBusy(nextBusy) {
  state.busy = nextBusy;
  elements.refreshButton.disabled = nextBusy;
  elements.captureActiveButton.disabled = nextBusy;
}

function setStatus(message, tone = '') {
  elements.statusText.textContent = message;
  elements.statusText.className = `status ${tone}`.trim();
}

function titleCase(value) {
  const text = String(value || '').trim();
  if (!text) {
    return 'Chromium';
  }
  return text.charAt(0).toUpperCase() + text.slice(1);
}

function escapeHtml(value) {
  return String(value || '')
    .replaceAll('&', '&amp;')
    .replaceAll('<', '&lt;')
    .replaceAll('>', '&gt;')
    .replaceAll('"', '&quot;')
    .replaceAll("'", '&#39;');
}

function sendMessage(message) {
  return new Promise((resolve, reject) => {
    chrome.runtime.sendMessage(message, (response) => {
      const error = chrome.runtime.lastError;
      if (error) {
        reject(new Error(error.message));
        return;
      }

      if (!response) {
        reject(new Error('No response received from the extension runtime'));
        return;
      }

      if (response.ok === false) {
        reject(new Error(response.error || 'Unknown extension error'));
        return;
      }

      resolve(response);
    });
  });
}
