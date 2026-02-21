const launcherEl = document.getElementById('launcher');
const editorEl = document.getElementById('editor');
const canvas = document.getElementById('canvas');
const ctx = canvas.getContext('2d');
const hintEl = document.getElementById('hint');
const toolbarEl = document.getElementById('annotationToolbar');

const toolEl = document.getElementById('tool');
const colorEl = document.getElementById('color');
const undoBtn = document.getElementById('undo');
const redoBtn = document.getElementById('redo');
const saveBtn = document.getElementById('save');
const cancelBtn = document.getElementById('cancel');

const captureRegionBtn = document.getElementById('captureRegion');
const captureScreenBtn = document.getElementById('captureScreen');

let phase = 'idle'; // idle | selecting | annotating
let captureMode = ''; // region | screen
let selection = null;
let selecting = null;

let ops = [];
let undone = [];
let drag = null;
let baseImagePath = '';
let baseImage = null;
let imageView = null;

function backend() {
  if (window.go?.main?.App) return window.go.main.App;
  return {
    async StartCapture() {
      alert('Wails backend not bound.');
      return null;
    },
    async SaveAnnotated() {
      alert('Wails backend not bound.');
      return null;
    },
    async EnterEditorMode() {},
    async ExitEditorMode() {},
    async CheckCapturePermission() { return { granted: true }; },
    async OpenCapturePermissionSettings() {},
    async LoadCaptureImage() { return ''; }
  };
}

function debugLog(...args) {
  console.log('[go-wails-shot]', ...args);
}

function showError(message) {
  if (!message) return;
  alert(message);
}

function describeError(err) {
  if (!err) return 'Unknown capture error';
  if (typeof err === 'string') return err;
  if (err.message) return err.message;
  try {
    return JSON.stringify(err);
  } catch {
    return 'Unknown capture error';
  }
}

async function enterLauncherWindowMode() {
  try {
    debugLog('ExitEditorMode -> launcher');
    await backend().ExitEditorMode();
  } catch (err) {
    debugLog('ExitEditorMode failed', err);
  }
}

function enterIdle() {
  debugLog('enterIdle');
  phase = 'idle';
  captureMode = '';
  selection = null;
  selecting = null;
  drag = null;
  ops = [];
  undone = [];
  baseImage = null;
  baseImagePath = '';
  imageView = null;

  launcherEl.classList.remove('hidden');
  editorEl.classList.add('hidden');
  hintEl.classList.add('hidden');
  toolbarEl.classList.add('hidden');
  canvas.style.cursor = 'default';
  draw();
  void enterLauncherWindowMode();
}

function enterEditor() {
  launcherEl.classList.add('hidden');
  editorEl.classList.remove('hidden');
  resizeCanvas();
}

function draw() {
  ctx.clearRect(0, 0, canvas.width, canvas.height);
  if (!baseImage) return;
  const view = getImageView();
  if (!view) return;

  ctx.fillStyle = '#020617';
  ctx.fillRect(0, 0, canvas.width, canvas.height);
  ctx.drawImage(baseImage, view.x, view.y, view.w, view.h);

  ctx.save();
  ctx.translate(view.x, view.y);
  ctx.scale(view.scale, view.scale);
  if (selection) {
    ctx.beginPath();
    ctx.rect(selection.x, selection.y, selection.w, selection.h);
    ctx.clip();
  }
  for (const op of ops) drawOp(ctx, op);
  if (drag) drawOp(ctx, { kind: drag.kind, payload: drag.payload });
  ctx.restore();
}

function drawOp(ctx, op) {
  const p = op.payload;
  ctx.strokeStyle = p.color || '#ff3b30';
  ctx.fillStyle = p.color || '#ff3b30';
  ctx.lineWidth = p.strokeWidth || 2;

  if (op.kind === 'rect') {
    ctx.strokeRect(p.x, p.y, p.w, p.h);
    return;
  }
  if (op.kind === 'line') {
    ctx.beginPath();
    ctx.moveTo(p.x1, p.y1);
    ctx.lineTo(p.x2, p.y2);
    ctx.stroke();
    return;
  }
  if (op.kind === 'arrow') {
    drawArrow(ctx, p);
    return;
  }
  if (op.kind === 'text') {
    ctx.font = `${p.size || 18}px ui-monospace`;
    ctx.fillText(p.text || 'Text', p.x, p.y);
    return;
  }
  if (op.kind === 'blur' || op.kind === 'pixelate') {
    ctx.strokeStyle = '#f59e0b';
    ctx.strokeRect(p.x, p.y, p.w, p.h);
  }
}

function drawArrow(ctx, p) {
  const x1 = p.x1;
  const y1 = p.y1;
  const x2 = p.x2;
  const y2 = p.y2;
  const headSize = p.headSize || 14;
  const angle = Math.atan2(y2 - y1, x2 - x1);

  ctx.beginPath();
  ctx.moveTo(x1, y1);
  ctx.lineTo(x2, y2);
  ctx.stroke();

  const a1 = angle + Math.PI * 0.82;
  const a2 = angle - Math.PI * 0.82;
  const hx1 = x2 + headSize * Math.cos(a1);
  const hy1 = y2 + headSize * Math.sin(a1);
  const hx2 = x2 + headSize * Math.cos(a2);
  const hy2 = y2 + headSize * Math.sin(a2);

  ctx.beginPath();
  ctx.moveTo(x2, y2);
  ctx.lineTo(hx1, hy1);
  ctx.moveTo(x2, y2);
  ctx.lineTo(hx2, hy2);
  ctx.stroke();
}

function canvasPoint(e) {
  const rect = canvas.getBoundingClientRect();
  const canvasX = (e.clientX - rect.left) * (canvas.width / rect.width);
  const canvasY = (e.clientY - rect.top) * (canvas.height / rect.height);
  const view = getImageView();
  if (!view) return { x: 0, y: 0 };
  const x = (canvasX - view.x) / view.scale;
  const y = (canvasY - view.y) / view.scale;
  return {
    x: Math.round(x),
    y: Math.round(y)
  };
}

function normalizeRect(startX, startY, endX, endY) {
  const x = Math.min(startX, endX);
  const y = Math.min(startY, endY);
  const w = Math.abs(endX - startX);
  const h = Math.abs(endY - startY);
  return { x, y, w, h };
}

function inSelection(pt) {
  if (!selection) return false;
  return pt.x >= selection.x && pt.x <= selection.x + selection.w && pt.y >= selection.y && pt.y <= selection.y + selection.h;
}

function resizeCanvas() {
  const dpr = window.devicePixelRatio || 1;
  const rect = canvas.getBoundingClientRect();
  const nextW = Math.max(1, Math.floor(rect.width * dpr));
  const nextH = Math.max(1, Math.floor(rect.height * dpr));
  if (canvas.width !== nextW || canvas.height !== nextH) {
    canvas.width = nextW;
    canvas.height = nextH;
  }
  imageView = null;
}

function getImageView() {
  if (!baseImage) return null;
  if (imageView) return imageView;
  const iw = baseImage.naturalWidth || baseImage.width;
  const ih = baseImage.naturalHeight || baseImage.height;
  const cw = canvas.width;
  const ch = canvas.height;
  const scale = Math.min(cw / iw, ch / ih, 1);
  const w = iw * scale;
  const h = ih * scale;
  const x = Math.floor((cw - w) / 2);
  const y = Math.floor((ch - h) / 2);
  imageView = { x, y, w, h, scale };
  return imageView;
}

function normalizePayload(kind, p) {
  if (kind === 'rect' || kind === 'blur' || kind === 'pixelate') {
    const x = p.w < 0 ? p.x + p.w : p.x;
    const y = p.h < 0 ? p.y + p.h : p.y;
    const w = Math.abs(p.w);
    const h = Math.abs(p.h);
    if (kind === 'blur') return { x, y, w, h, radius: 3 };
    if (kind === 'pixelate') return { x, y, w, h, size: 12 };
    return { x, y, w, h, color: p.color, strokeWidth: p.strokeWidth, fill: false };
  }
  return p;
}

function pushOp(op) {
  const id = crypto.randomUUID?.() || `${Date.now()}-${Math.random()}`;
  ops.push({ id, kind: op.kind, z: ops.length, payload: op.payload });
  undone = [];
  draw();
}

function positionToolbar() {
  if (!selection) return;
  const view = getImageView();
  if (!view) return;
  const rect = canvas.getBoundingClientRect();
  const canvasToCssX = rect.width / canvas.width;
  const canvasToCssY = rect.height / canvas.height;

  const selLeftCanvas = view.x + selection.x * view.scale;
  const selTopCanvas = view.y + selection.y * view.scale;
  const selWidthCanvas = selection.w * view.scale;

  const selLeftPx = selLeftCanvas * canvasToCssX;
  const selTopPx = selTopCanvas * canvasToCssY;
  const selWidthPx = selWidthCanvas * canvasToCssX;

  const top = Math.max(10, selTopPx - 52);
  let left = selLeftPx;
  if (captureMode === 'screen') {
    left = selLeftPx + selWidthPx / 2;
    toolbarEl.style.transform = 'translateX(-50%)';
  } else {
    toolbarEl.style.transform = 'none';
  }

  toolbarEl.style.top = `${top}px`;
  toolbarEl.style.left = `${left}px`;
}

async function beginCapture(mode) {
  debugLog('beginCapture', { mode });
  captureRegionBtn.disabled = true;
  captureScreenBtn.disabled = true;

  try {
    const perm = await backend().CheckCapturePermission();
    debugLog('CheckCapturePermission result', perm);
    if (perm && perm.granted === false) {
      const message = perm.message || 'Screen capture permission is required.';
      const hint = perm.settingsHint ? `\n\n${perm.settingsHint}` : '';
      const allow = confirm(`${message}${hint}\n\nOpen system privacy settings now?`);
      if (allow) {
        await backend().OpenCapturePermissionSettings();
      }
      enterIdle();
      return;
    }

    debugLog('StartCapture requested', { requestMode: mode, uiMode: mode });
    const result = await backend().StartCapture(mode);
    if (!result?.imagePath) {
      debugLog('StartCapture returned empty result', result);
      showError('Capture failed: no image was returned.');
      enterIdle();
      return;
    }
    debugLog('StartCapture success', result);

    baseImagePath = result.imagePath;
    captureMode = mode;
    ops = [];
    undone = [];

    const img = new Image();
    const dataURL = await backend().LoadCaptureImage(result.imagePath);
    if (!dataURL) {
      throw new Error('LoadCaptureImage returned empty data');
    }
    img.src = dataURL;
    await img.decode();
    baseImage = img;
    imageView = null;

    debugLog('EnterEditorMode -> fullscreen (post-capture)');
    await backend().EnterEditorMode();
    enterEditor();

    if (mode === 'screen' || mode === 'region') {
      phase = 'annotating';
      selection = { x: 0, y: 0, w: canvas.width, h: canvas.height };
      canvas.style.cursor = 'crosshair';
      hintEl.classList.add('hidden');
      toolbarEl.classList.remove('hidden');
      positionToolbar();
    }
    draw();
  } catch (err) {
    const msg = describeError(err);
    debugLog('StartCapture failed', err, msg);
    if (isPermissionError(msg)) {
      const allow = confirm('Screen Recording permission is required. Open macOS Privacy settings now?');
      if (allow) {
        try {
          await backend().OpenCapturePermissionSettings();
          debugLog('Opened Screen Recording settings');
        } catch (openErr) {
          debugLog('Failed to open Screen Recording settings', openErr);
        }
      }
    }
    showError(`Capture failed: ${describeError(err)}`);
    enterIdle();
  } finally {
    captureRegionBtn.disabled = false;
    captureScreenBtn.disabled = false;
  }
}

function isPermissionError(msg) {
  const lower = String(msg || '').toLowerCase();
  return lower.includes('err_capture_permission') ||
    (lower.includes('screen recording') && lower.includes('permission')) ||
    lower.includes('not authorized');
}

canvas.addEventListener('mousedown', (e) => {
  if (phase === 'idle' || !baseImage) return;
  const pt = canvasPoint(e);

  if (phase === 'selecting') {
    selecting = { x: pt.x, y: pt.y, w: 0, h: 0, startX: pt.x, startY: pt.y };
    return;
  }

  if (phase !== 'annotating') return;
  if (!inSelection(pt)) return;

  const kind = toolEl.value;
  if (kind === 'text') {
    const text = prompt('Text');
    if (!text) return;
    pushOp({ kind, payload: { x: pt.x, y: pt.y, text, color: colorEl.value, size: 18 } });
    return;
  }

  drag = {
    kind,
    startX: pt.x,
    startY: pt.y,
    payload: { x: pt.x, y: pt.y, w: 0, h: 0, color: colorEl.value, strokeWidth: 2 }
  };
});

canvas.addEventListener('mousemove', (e) => {
  if (!baseImage) return;
  const pt = canvasPoint(e);

  if (phase === 'selecting' && selecting) {
    const r = normalizeRect(selecting.startX, selecting.startY, pt.x, pt.y);
    selecting.x = r.x;
    selecting.y = r.y;
    selecting.w = r.w;
    selecting.h = r.h;
    draw();
    return;
  }

  if (phase !== 'annotating' || !drag) return;

  if (drag.kind === 'rect' || drag.kind === 'blur' || drag.kind === 'pixelate') {
    drag.payload.w = pt.x - drag.startX;
    drag.payload.h = pt.y - drag.startY;
  } else {
    drag.payload.x1 = drag.startX;
    drag.payload.y1 = drag.startY;
    drag.payload.x2 = pt.x;
    drag.payload.y2 = pt.y;
    drag.payload.color = colorEl.value;
    drag.payload.strokeWidth = 2;
    if (drag.kind === 'arrow') drag.payload.headSize = 14;
  }
  draw();
});

canvas.addEventListener('mouseup', () => {
  if (phase === 'selecting' && selecting) {
    if (selecting.w < 8 || selecting.h < 8) {
      selecting = null;
      draw();
      return;
    }
    selection = { x: selecting.x, y: selecting.y, w: selecting.w, h: selecting.h };
    selecting = null;
    phase = 'annotating';
    hintEl.classList.add('hidden');
    toolbarEl.classList.remove('hidden');
    positionToolbar();
    draw();
    return;
  }

  if (phase !== 'annotating' || !drag) return;
  const payload = normalizePayload(drag.kind, drag.payload);
  pushOp({ kind: drag.kind, payload });
  drag = null;
});

undoBtn.addEventListener('click', () => {
  if (!ops.length) return;
  undone.push(ops.pop());
  draw();
});

redoBtn.addEventListener('click', () => {
  if (!undone.length) return;
  ops.push(undone.pop());
  draw();
});

captureRegionBtn.addEventListener('click', async () => {
  debugLog('Capture Region clicked');
  await beginCapture('region');
});

captureScreenBtn.addEventListener('click', async () => {
  debugLog('Capture Full Screen clicked');
  await beginCapture('screen');
});

cancelBtn.addEventListener('click', () => {
  debugLog('Cancel clicked');
  enterIdle();
});

saveBtn.addEventListener('click', async () => {
  debugLog('Save clicked');
  if (!baseImagePath) {
    alert('Capture first');
    return;
  }
  try {
    const req = {
      baseImagePath,
      ops,
      format: 'png',
      quality: 90,
      outputPath: ''
    };
    debugLog('Save request', req);
    const result = await backend().SaveAnnotated(req);
    debugLog('Save response', result);
    if (result?.outputPath) {
      alert(`Saved: ${result.outputPath}`);
    } else {
      showError('Save failed: empty response');
    }
  } catch (err) {
    const msg = describeError(err);
    debugLog('Save failed', err, msg);
    showError(`Save failed: ${msg}`);
  }
});

window.addEventListener('resize', () => {
  resizeCanvas();
  if (phase === 'annotating') positionToolbar();
  draw();
});

enterIdle();
