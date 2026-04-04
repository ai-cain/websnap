export function countTabs(snapshot) {
  return (snapshot?.windows || []).reduce((total, currentWindow) => total + (currentWindow.tabs?.length || 0), 0);
}

export function findTabInSnapshot(snapshot, windowId, tabId) {
  for (const currentWindow of snapshot?.windows || []) {
    if (currentWindow.windowId !== Number(windowId)) {
      continue;
    }

    for (const tab of currentWindow.tabs || []) {
      if (tab.tabId === Number(tabId)) {
        return { ...tab, windowId: currentWindow.windowId };
      }
    }
  }

  return null;
}

export function resolveSelectedTab(snapshot, selectedWindowId, selectedTabId, activeTab) {
  const explicit = findTabInSnapshot(snapshot, selectedWindowId, selectedTabId);
  if (explicit) {
    return explicit;
  }

  if (activeTab) {
    return activeTab;
  }

  return null;
}
