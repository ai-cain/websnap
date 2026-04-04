export function detectBrowserName(userAgent = "") {
  const ua = String(userAgent).toLowerCase();
  if (ua.includes("edg/")) {
    return "edge";
  }
  if (ua.includes("opr/")) {
    return "opera";
  }
  if (ua.includes("brave")) {
    return "brave";
  }
  return "chrome";
}

export function isWebUrl(url) {
  return /^https?:\/\//i.test(String(url || "").trim());
}

export function isWebTab(tab) {
  return Boolean(tab) && isWebUrl(tab.url || tab.pendingUrl || "");
}

export function formatHost(url) {
  try {
    return new URL(url).host.replace(/^www\./i, "");
  } catch {
    return "";
  }
}

export function normalizeTab(tab) {
  const url = tab?.url || tab?.pendingUrl || "";
  return {
    tabId: Number(tab?.id ?? tab?.tabId ?? 0),
    index: Number(tab?.index ?? 0),
    title: String(tab?.title || url || "Untitled tab"),
    url,
    active: Boolean(tab?.active)
  };
}

export function buildSnapshot(windows, browser) {
  return {
    browser,
    windows: (windows || [])
      .map((currentWindow) => normalizeWindow(currentWindow, browser))
      .filter(Boolean)
  };
}

export function getActiveTabFromSnapshot(snapshot) {
  for (const currentWindow of snapshot?.windows || []) {
    for (const tab of currentWindow.tabs || []) {
      if (tab.active) {
        return { ...tab, windowId: currentWindow.windowId };
      }
    }
  }

  const firstWindow = snapshot?.windows?.[0];
  const firstTab = firstWindow?.tabs?.[0];
  if (!firstWindow || !firstTab) {
    return null;
  }

  return { ...firstTab, windowId: firstWindow.windowId };
}

function normalizeWindow(currentWindow, browser) {
  const tabs = (currentWindow?.tabs || [])
    .filter(isWebTab)
    .map(normalizeTab);

  if (tabs.length === 0) {
    return null;
  }

  const activeTab = tabs.find((tab) => tab.active) || tabs[0];
  return {
    windowId: Number(currentWindow?.id ?? currentWindow?.windowId ?? 0),
    appName: browser,
    title: String(activeTab?.title || `${browser} window`),
    tabs
  };
}
