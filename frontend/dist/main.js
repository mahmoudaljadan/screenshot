const canvas = document.getElementById('canvas');
const ctx = canvas.getContext('2d');
const opsEl = document.getElementById('ops');
const toolEl = document.getElementById('tool');
const colorEl = document.getElementById('color');

let ops = [];
let undone = [];
let drag = null;
let baseImagePath = '';
let baseImage = null;

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
    }
  };
}

function draw() {
  ctx.clearRect(0, 0, canvas.width, canvas.height);
  if (baseImage) ctx.drawImage(baseImage, 0, 0, canvas.width, canvas.height);
  for (const op of ops) drawOp(ctx, op);
  if (drag) drawPreview(ctx, drag);
  opsEl.textContent = JSON.stringify(ops, null, 2);
}

function drawOp(ctx, op) {
  const p = op.payload;
  ctx.strokeStyle = p.color || '#ff3b30';
  ctx.fillStyle = p.color || '#ff3b30';
  ctx.lineWidth = p.strokeWidth || 2;
  if (op.kind === 'rect') {
    ctx.strokeRect(p.x, p.y, p.w, p.h);
  } else if (op.kind === 'line' || op.kind === 'arrow') {
    ctx.beginPath();
    ctx.moveTo(p.x1, p.y1);
    ctx.lineTo(p.x2, p.y2);
    ctx.stroke();
  } else if (op.kind === 'text') {
    ctx.font = `${p.size || 18}px ui-monospace`;
    ctx.fillText(p.text || 'Text', p.x, p.y);
  } else if (op.kind === 'blur' || op.kind === 'pixelate') {
    ctx.strokeStyle = '#f59e0b';
    ctx.strokeRect(p.x, p.y, p.w, p.h);
  }
}

function drawPreview(ctx, d) {
  drawOp(ctx, { kind: d.kind, payload: d.payload });
}

function canvasPoint(e) {
  const rect = canvas.getBoundingClientRect();
  const scaleX = canvas.width / rect.width;
  const scaleY = canvas.height / rect.height;
  return {
    x: Math.round((e.clientX - rect.left) * scaleX),
    y: Math.round((e.clientY - rect.top) * scaleY)
  };
}

canvas.addEventListener('mousedown', (e) => {
  const { x, y } = canvasPoint(e);
  const kind = toolEl.value;
  if (kind === 'text') {
    const text = prompt('Text');
    if (!text) return;
    pushOp({ kind, payload: { x, y, text, color: colorEl.value, size: 18 } });
    return;
  }
  drag = { kind, startX: x, startY: y, payload: { x, y, w: 0, h: 0, color: colorEl.value, strokeWidth: 2 } };
});

canvas.addEventListener('mousemove', (e) => {
  if (!drag) return;
  const { x, y } = canvasPoint(e);
  if (drag.kind === 'rect' || drag.kind === 'blur' || drag.kind === 'pixelate') {
    drag.payload.w = x - drag.startX;
    drag.payload.h = y - drag.startY;
  } else {
    drag.payload.x1 = drag.startX;
    drag.payload.y1 = drag.startY;
    drag.payload.x2 = x;
    drag.payload.y2 = y;
    drag.payload.color = colorEl.value;
    drag.payload.strokeWidth = 2;
  }
  draw();
});

canvas.addEventListener('mouseup', () => {
  if (!drag) return;
  const payload = normalizePayload(drag.kind, drag.payload);
  pushOp({ kind: drag.kind, payload });
  drag = null;
});

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

document.getElementById('undo').addEventListener('click', () => {
  if (!ops.length) return;
  undone.push(ops.pop());
  draw();
});

document.getElementById('redo').addEventListener('click', () => {
  if (!undone.length) return;
  ops.push(undone.pop());
  draw();
});

document.getElementById('capture').addEventListener('click', async () => {
  const result = await backend().StartCapture('screen');
  if (!result?.imagePath) return;
  baseImagePath = result.imagePath;
  const img = new Image();
  img.src = `file://${result.imagePath}`;
  await img.decode();
  baseImage = img;
  draw();
});

document.getElementById('save').addEventListener('click', async () => {
  if (!baseImagePath) {
    alert('Capture first');
    return;
  }
  const req = {
    baseImagePath,
    ops,
    format: 'png',
    quality: 90,
    outputPath: ''
  };
  const result = await backend().SaveAnnotated(req);
  if (result?.outputPath) alert(`Saved: ${result.outputPath}`);
});

draw();
