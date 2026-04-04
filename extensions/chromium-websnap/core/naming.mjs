const BROWSER_SUFFIXES = [
  ' - Google Chrome',
  ' - Chrome',
  ' - Microsoft Edge',
  ' - Edge',
  ' - Brave',
  ' - Opera'
];

const GENERIC_SEGMENTS = new Set([
  'google-chrome',
  'chrome',
  'microsoft-edge',
  'edge',
  'brave',
  'opera',
  'file-explorer',
  'explorador-de-archivos'
]);

export function stripBrowserTitleSuffix(value = '') {
  const trimmed = String(value).trim();
  for (const suffix of BROWSER_SUFFIXES) {
    if (trimmed.endsWith(suffix)) {
      return trimmed.slice(0, -suffix.length).trim();
    }
  }
  return trimmed;
}

export function splitTitleSegments(value = '') {
  let normalized = String(value);
  for (const separator of [' — ', ' – ', ' | ', ' • ']) {
    normalized = normalized.split(separator).join(' - ');
  }

  return normalized
    .split(' - ')
    .map((part) => part.trim())
    .filter(Boolean);
}

export function sanitizeFileStem(value = '') {
  const input = String(value).trim().toLowerCase();
  if (!input) {
    return '';
  }

  let output = '';
  let lastDash = false;
  for (const char of input) {
    if (/^[\p{L}\p{N}]$/u.test(char)) {
      output += char;
      lastDash = false;
      continue;
    }

    if (char === '.' || char === '-' || char === '_') {
      output += char;
      lastDash = false;
      continue;
    }

    if (!lastDash) {
      output += '-';
      lastDash = true;
    }
  }

  return output.replace(/^[-._]+|[-._]+$/g, '');
}

export function preferredCaptureStem({ title = '', url = '', browserName = '' } = {}) {
  const cleanedTitle = stripBrowserTitleSuffix(title);
  const segments = splitTitleSegments(cleanedTitle);

  let best = '';
  let bestScore = -Infinity;
  for (const segment of segments) {
    if (!segment || isGenericSegment(segment, browserName)) {
      continue;
    }

    const score = scoreSegment(segment);
    if (score > bestScore) {
      best = segment;
      bestScore = score;
    }
  }

  const candidate = best || segments[0] || hostLabel(url) || browserName || 'capture';
  const sanitized = sanitizeFileStem(candidate);
  if (sanitized) {
    return sanitized;
  }

  return sanitizeFileStem(hostLabel(url)) || sanitizeFileStem(browserName) || 'capture';
}

export function buildCaptureFilename({ title = '', url = '', browserName = '', now = new Date() } = {}) {
  const stamp = formatTimestamp(now);
  const stem = preferredCaptureStem({ title, url, browserName });
  return `websnap/${stamp}-${stem}.png`;
}

function hostLabel(url) {
  try {
    return new URL(url).hostname.replace(/^www\./i, '');
  } catch {
    return '';
  }
}

function isGenericSegment(segment, browserName) {
  const normalized = sanitizeFileStem(segment);
  if (!normalized) {
    return true;
  }

  return GENERIC_SEGMENTS.has(normalized) || normalized === sanitizeFileStem(browserName);
}

function scoreSegment(segment) {
  let letters = 0;
  let digits = 0;
  let separators = 0;
  for (const char of segment) {
    if (/^\p{L}$/u.test(char)) {
      letters += 1;
      continue;
    }
    if (/^\p{N}$/u.test(char)) {
      digits += 1;
      continue;
    }
    if (/^[\s_-]$/u.test(char)) {
      separators += 1;
    }
  }

  return (letters * 4) + separators - (digits * 2);
}

function formatTimestamp(value) {
  const date = value instanceof Date ? value : new Date(value);
  const year = date.getFullYear();
  const month = String(date.getMonth() + 1).padStart(2, '0');
  const day = String(date.getDate()).padStart(2, '0');
  const hours = String(date.getHours()).padStart(2, '0');
  const minutes = String(date.getMinutes()).padStart(2, '0');
  const seconds = String(date.getSeconds()).padStart(2, '0');
  return `${year}${month}${day}-${hours}${minutes}${seconds}`;
}
