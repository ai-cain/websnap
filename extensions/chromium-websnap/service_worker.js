const SERVER_URL = "ws://127.0.0.1:38971/ws";

let socket = null;
let reconnectTimer = null;

function detectBrowserName() {
  const ua = self.navigator.userAgent.toLowerCase();
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

function isWebTab(tab) {
  if (!tab || typeof tab.url !== "string") {
    return false;
  }

  return tab.url.startsWith("http://") || tab.url.startsWith("https://");
}

async function buildSnapshot() {
  const browser = detectBrowserName();
  const windows = await chrome.windows.getAll({
    populate: true,
    windowTypes: ["normal"]
  });

  const payloadWindows = windows
    .map((window) => {
      const tabs = (window.tabs || [])
        .filter(isWebTab)
        .map((tab) => ({
          tabId: tab.id,
          index: tab.index,
          title: tab.title || tab.url || "Untitled tab",
          url: tab.url || "",
          active: Boolean(tab.active)
        }));

      if (tabs.length === 0) {
        return null;
      }

      const activeTab = tabs.find((tab) => tab.active) || tabs[0];
      return {
        windowId: window.id,
        appName: browser,
        title: activeTab?.title || `${browser} window`,
        tabs
      };
    })
    .filter(Boolean);

  return {
    type: "snapshot",
    browser,
    windows: payloadWindows
  };
}

function sendMessage(message) {
  if (!socket || socket.readyState !== WebSocket.OPEN) {
    return;
  }

  socket.send(JSON.stringify(message));
}

async function sendSnapshot() {
  try {
    sendMessage(await buildSnapshot());
  } catch (error) {
    console.warn("websnap: unable to build snapshot", error);
  }
}

function scheduleReconnect() {
  if (reconnectTimer !== null) {
    return;
  }

  reconnectTimer = setTimeout(() => {
    reconnectTimer = null;
    connect();
  }, 1500);
}

async function captureWindowTab(message) {
  const windowId = Number(message.windowId);
  const tabId = Number(message.tabId);

  if (!windowId || !tabId) {
    sendMessage({
      id: message.id,
      type: "capture-result",
      error: "windowId and tabId are required"
    });
    return;
  }

  try {
    await chrome.windows.update(windowId, { focused: true });
    await chrome.tabs.update(tabId, { active: true });
    await new Promise((resolve) => setTimeout(resolve, 350));

    const pngDataUrl = await chrome.tabs.captureVisibleTab(windowId, {
      format: "png"
    });

    sendMessage({
      id: message.id,
      type: "capture-result",
      pngDataUrl
    });

    await sendSnapshot();
  } catch (error) {
    sendMessage({
      id: message.id,
      type: "capture-result",
      error: error?.message || String(error)
    });
  }
}

function handleSocketMessage(event) {
  let message;
  try {
    message = JSON.parse(event.data);
  } catch (error) {
    console.warn("websnap: invalid server message", error);
    return;
  }

  if (message.type === "capture") {
    captureWindowTab(message);
  }
}

function connect() {
  if (socket && (socket.readyState === WebSocket.OPEN || socket.readyState === WebSocket.CONNECTING)) {
    return;
  }

  socket = new WebSocket(SERVER_URL);

  socket.addEventListener("open", async () => {
    sendMessage({
      type: "hello",
      browser: detectBrowserName()
    });
    await sendSnapshot();
  });

  socket.addEventListener("message", handleSocketMessage);

  socket.addEventListener("close", () => {
    socket = null;
    scheduleReconnect();
  });

  socket.addEventListener("error", () => {
    if (socket) {
      socket.close();
    }
  });
}

async function refreshSnapshotIfConnected() {
  if (socket && socket.readyState === WebSocket.OPEN) {
    await sendSnapshot();
    return;
  }

  connect();
}

chrome.runtime.onInstalled.addListener(() => {
  connect();
});

chrome.runtime.onStartup.addListener(() => {
  connect();
});

chrome.action.onClicked.addListener(() => {
  refreshSnapshotIfConnected();
});

chrome.tabs.onActivated.addListener(() => {
  refreshSnapshotIfConnected();
});

chrome.tabs.onUpdated.addListener(() => {
  refreshSnapshotIfConnected();
});

chrome.tabs.onRemoved.addListener(() => {
  refreshSnapshotIfConnected();
});

chrome.tabs.onCreated.addListener(() => {
  refreshSnapshotIfConnected();
});

chrome.windows.onFocusChanged.addListener(() => {
  refreshSnapshotIfConnected();
});

chrome.windows.onRemoved.addListener(() => {
  refreshSnapshotIfConnected();
});

chrome.windows.onCreated.addListener(() => {
  refreshSnapshotIfConnected();
});

connect();
