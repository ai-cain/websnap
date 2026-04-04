import { formatHost, getActiveTabFromSnapshot } from './core/browser.mjs';
import { countTabs, resolveSelectedTab } from './core/popup-state.mjs';

const state = {
  snapshot: null,
  selectedWindowId: 0,
  selectedTabId: 0,
  options: {
    saveAs: false
  },
  busy: false,
  view: 'capture'
};

const elements = {
  refreshButton: document.querySelector('#refreshButton'),
  captureActiveButton: document.querySelector('#captureActiveButton'),
  captureSelectedButton: document.querySelector('#captureSelectedButton'),
  saveAsCheckbox: document.querySelector('#saveAsCheckbox'),
  statusText: document.querySelector('#statusText'),
  windowsContainer: document.querySelector('#windowsContainer'),
  browserBadge: document.querySelector('#browserBadge'),
  tabCountBadge: document.querySelector('#tabCountBadge'),
  activeTabSummary: document.querySelector('#activeTabSummary'),
  selectedTabSummary: document.querySelector('#selectedTabSummary'),
  viewButtons: [...document.querySelectorAll('[data-view-button]')],
  panels: {
    capture: document.querySelector('#view-capture'),
    tabs: document.querySelector('#view-tabs'),
    settings: document.querySelector('#view-settings')
  }
};

bootstrap().catch((error) => setStatus(error.message || String(error), 'error'));

async function bootstrap() {
  bindEvents();
  setView(state.view);
  await refreshState();
}

function bindEvents() {
  elements.refreshButton.addEventListener('click', () => void refreshState());
  elements.captureActiveButton.addEventListener('click', () => void captureActiveTab());
  elements.captureSelectedButton.addEventListener('click', () => void captureCurrentSelection());
  elements.saveAsCheckbox.addEventListener('change', () => void updateOptions());
  elements.viewButtons.forEach((button) => {
    button.addEventListener('click', () => setView(button.dataset.view));
  });
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
    elements.tabCountBadge.textContent = `${countTabs(snapshot)} tabs`;

    primeSelection(snapshot);
    renderAll();

    if ((snapshot?.windows || []).length === 0) {
      setStatus('No http(s) tabs are open right now. Open a normal web page and refresh.');
    } else {
      setStatus('Choose a tab from the menu and capture it when ready.', 'success');
    }
  } catch (error) {
    setStatus(error.message || String(error), 'error');
    elements.windowsContainer.innerHTML = renderEmptyState('Unable to load tabs from the extension.');
    elements.activeTabSummary.innerHTML = renderSummaryEmpty('Nothing available yet.');
    elements.selectedTabSummary.innerHTML = renderSummaryEmpty('Nothing selected yet.');
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
    setStatus(elements.saveAsCheckbox.checked ? 'The extension will ask where to save each capture.' : 'Captures will download directly into the websnap folder.', 'success');
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

async function captureCurrentSelection() {
  const selected = getSelectedTab();
  if (!selected) {
    setStatus('Select a web tab first from the Tabs menu.', 'error');
    setView('tabs');
    return;
  }

  await captureSelectedTab(selected.windowId, selected.tabId);
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
    setView('capture');
    await refreshState();
  } catch (error) {
    setStatus(error.message || String(error), 'error');
  } finally {
    setBusy(false);
  }
}

function renderAll() {
  renderCapturePanel();
  renderWindows();
  updateCaptureButtonState();
}

function renderCapturePanel() {
  const active = getActiveTabFromSnapshot(state.snapshot);
  const selected = getSelectedTab();

  elements.activeTabSummary.innerHTML = active
    ? renderTabSummary(active, 'Current focused web tab')
    : renderSummaryEmpty('No active web tab is available right now.');

  elements.selectedTabSummary.innerHTML = selected
    ? renderTabSummary(selected, selected.active ? 'Currently selected and active' : 'Chosen from the Tabs menu')
    : renderSummaryEmpty('Pick a tab from the Tabs menu to capture it later.');
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
      renderAll();
      setView('capture');
      setStatus('Selected tab updated. You can capture it now.', 'success');
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
      <div class="tab-main">
        <div class="tab-title">${escapeHtml(tab.title)}</div>
        <div class="tab-url">${escapeHtml(host || tab.url)}</div>
      </div>
      <div class="tab-meta">
        ${activePill}
        <div class="tab-actions">
          <button type="button" data-select-tab="true" data-window-id="${windowId}" data-tab-id="${tab.tabId}">Select</button>
          <button type="button" data-capture-tab="true" data-window-id="${windowId}" data-tab-id="${tab.tabId}">Capture</button>
        </div>
      </div>
    </div>
  `;
}

function renderTabSummary(tab, subtitle) {
  return `
    <div class="summary-title">${escapeHtml(tab.title)}</div>
    <div class="summary-meta">${escapeHtml(formatHost(tab.url) || tab.url)}</div>
    <div class="summary-meta">${escapeHtml(subtitle)}</div>
  `;
}

function renderSummaryEmpty(message) {
  return `<div class="summary-empty">${escapeHtml(message)}</div>`;
}

function renderEmptyState(message) {
  return `<div class="empty-state">${escapeHtml(message)}</div>`;
}

function primeSelection(snapshot) {
  const active = getActiveTabFromSnapshot(snapshot);
  const selected = resolveSelectedTab(snapshot, state.selectedWindowId, state.selectedTabId, active);
  if (!selected) {
    state.selectedWindowId = 0;
    state.selectedTabId = 0;
    return;
  }

  state.selectedWindowId = selected.windowId;
  state.selectedTabId = selected.tabId;
}

function getSelectedTab() {
  return resolveSelectedTab(
    state.snapshot,
    state.selectedWindowId,
    state.selectedTabId,
    getActiveTabFromSnapshot(state.snapshot)
  );
}

function updateCaptureButtonState() {
  elements.captureSelectedButton.disabled = state.busy || !getSelectedTab();
}

function setBusy(nextBusy) {
  state.busy = nextBusy;
  elements.refreshButton.disabled = nextBusy;
  elements.captureActiveButton.disabled = nextBusy;
  elements.saveAsCheckbox.disabled = nextBusy;
  updateCaptureButtonState();
}

function setView(view) {
  state.view = ['capture', 'tabs', 'settings'].includes(view) ? view : 'capture';
  elements.viewButtons.forEach((button) => {
    button.classList.toggle('is-active', button.dataset.view === state.view);
  });
  for (const [name, panel] of Object.entries(elements.panels)) {
    panel.hidden = name !== state.view;
  }
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
