package api

// dashboardHTML is the entire admin SPA: sidebar nav, 6 sections, live data.
// Single file by design so deployment is just the Go binary.
const dashboardHTML = `<!doctype html>
<html lang="tr">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>Oroya Admin</title>
  <style>
    :root {
      --bg: #F5F0EB;
      --surface: #FFFFFF;
      --primary: #E86A33;
      --primary-dark: #D35400;
      --text: #2C2C2C;
      --muted: #6B6B6B;
      --border: #E0D8D0;
      --success: #27AE60;
      --error: #C0392B;
      --code-bg: #FAF6F2;
    }
    * { box-sizing: border-box; }
    body {
      margin: 0;
      font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", system-ui, sans-serif;
      background: var(--bg);
      color: var(--text);
      font-size: 14px;
      line-height: 1.5;
    }
    .layout { display: flex; min-height: 100vh; }

    /* --- Sidebar --- */
    .sidebar {
      width: 240px;
      background: var(--surface);
      border-right: 1px solid var(--border);
      padding: 24px 0;
      flex-shrink: 0;
      position: sticky;
      top: 0;
      height: 100vh;
      overflow-y: auto;
    }
    .brand {
      padding: 0 24px 24px;
      font-size: 20px;
      font-weight: 700;
      color: var(--primary);
      border-bottom: 1px solid var(--border);
      margin-bottom: 16px;
    }
    .brand small { display: block; font-size: 11px; color: var(--muted); font-weight: 400; margin-top: 4px; }
    .nav { list-style: none; padding: 0; margin: 0; }
    .nav li a {
      display: block;
      padding: 10px 24px;
      color: var(--text);
      text-decoration: none;
      border-left: 3px solid transparent;
      cursor: pointer;
    }
    .nav li a:hover { background: var(--bg); }
    .nav li a.active {
      background: var(--bg);
      border-left-color: var(--primary);
      color: var(--primary);
      font-weight: 600;
    }
    .nav-section {
      padding: 12px 24px 4px;
      font-size: 11px;
      text-transform: uppercase;
      letter-spacing: 0.5px;
      color: var(--muted);
    }

    /* --- Main --- */
    .main {
      flex: 1;
      padding: 32px 40px;
      max-width: 1200px;
    }
    .page { display: none; }
    .page.active { display: block; }
    h1 { margin: 0 0 24px; font-size: 28px; }
    h2 { margin: 0 0 16px; font-size: 18px; }
    h3 { margin: 0 0 8px; font-size: 14px; color: var(--muted); text-transform: uppercase; letter-spacing: 0.5px; }
    .muted { color: var(--muted); }
    .lead { color: var(--muted); font-size: 15px; margin-bottom: 24px; }

    /* --- Cards --- */
    .card {
      background: var(--surface);
      border: 1px solid var(--border);
      border-radius: 8px;
      padding: 20px;
      margin-bottom: 16px;
    }
    .card.featured { border-color: var(--primary); border-width: 2px; }
    .card-grid {
      display: grid;
      grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
      gap: 16px;
      margin-bottom: 24px;
    }
    .stat {
      background: var(--surface);
      border: 1px solid var(--border);
      border-radius: 8px;
      padding: 16px;
    }
    .stat .label { color: var(--muted); font-size: 12px; text-transform: uppercase; letter-spacing: 0.5px; }
    .stat .value { font-size: 28px; font-weight: 700; margin-top: 4px; color: var(--primary); }

    /* --- Code blocks --- */
    pre, code { font-family: "SF Mono", Menlo, Monaco, Consolas, monospace; font-size: 13px; }
    code { background: var(--code-bg); padding: 2px 6px; border-radius: 3px; }
    pre {
      background: var(--code-bg);
      border: 1px solid var(--border);
      border-radius: 6px;
      padding: 12px 16px;
      overflow-x: auto;
      margin: 0;
    }
    .copy-row {
      display: flex;
      align-items: center;
      gap: 8px;
      margin-top: 8px;
    }
    button {
      background: var(--primary);
      color: white;
      border: none;
      padding: 8px 16px;
      border-radius: 4px;
      cursor: pointer;
      font-size: 13px;
      font-weight: 500;
    }
    button:hover { background: var(--primary-dark); }
    button.ghost { background: transparent; color: var(--text); border: 1px solid var(--border); }
    button.ghost:hover { background: var(--code-bg); }

    /* --- Tables --- */
    table { width: 100%; border-collapse: collapse; }
    th, td { text-align: left; padding: 8px 12px; border-bottom: 1px solid var(--border); }
    th { font-weight: 600; color: var(--muted); font-size: 12px; text-transform: uppercase; letter-spacing: 0.5px; }
    td.num { font-family: "SF Mono", Menlo, monospace; text-align: right; }
    .method {
      display: inline-block;
      padding: 2px 8px;
      border-radius: 3px;
      font-size: 11px;
      font-weight: 700;
      font-family: "SF Mono", Menlo, monospace;
      min-width: 50px;
      text-align: center;
    }
    .method.GET    { background: #E3F2FD; color: #1976D2; }
    .method.POST   { background: #E8F5E9; color: #2E7D32; }
    .method.PUT    { background: #FFF3E0; color: #E65100; }
    .method.DELETE { background: #FFEBEE; color: #C62828; }
    .lock { color: var(--primary); margin-left: 4px; font-size: 12px; }
    .path { font-family: "SF Mono", Menlo, monospace; }

    /* --- Status --- */
    .pill {
      display: inline-block;
      padding: 3px 10px;
      border-radius: 12px;
      font-size: 11px;
      font-weight: 600;
      text-transform: uppercase;
      letter-spacing: 0.5px;
    }
    .pill.ok { background: rgba(39, 174, 96, 0.15); color: var(--success); }
    .pill.bad { background: rgba(192, 57, 43, 0.15); color: var(--error); }
    .pill.warn { background: rgba(232, 106, 51, 0.15); color: var(--primary); }

    .row { display: flex; gap: 12px; align-items: center; flex-wrap: wrap; }
    .row > * { margin: 0; }

    /* --- Token gate banner --- */
    .token-banner {
      background: var(--surface);
      border: 1px solid var(--border);
      border-radius: 8px;
      padding: 12px 16px;
      margin-bottom: 16px;
      display: flex;
      gap: 8px;
      align-items: center;
    }
    .token-banner input {
      flex: 1;
      padding: 8px 12px;
      border: 1px solid var(--border);
      border-radius: 4px;
      font-family: "SF Mono", Menlo, monospace;
      font-size: 13px;
    }
    .token-banner.set { border-color: var(--success); }
  </style>
</head>
<body>
<div class="layout">

  <aside class="sidebar">
    <div class="brand">
      Oroya Admin
      <small>backend control panel</small>
    </div>
    <ul class="nav">
      <li class="nav-section">Frontend</li>
      <li><a data-page="connection" class="active">Connection Info</a></li>
      <li><a data-page="endpoints">API Endpoints</a></li>
      <li><a data-page="auth">Auth & Token</a></li>
      <li class="nav-section">Monitoring</li>
      <li><a data-page="health">Health</a></li>
      <li><a data-page="stats">Statistics</a></li>
      <li><a data-page="queue">Queue & Worker</a></li>
      <li><a data-page="storage">Storage</a></li>
    </ul>
  </aside>

  <main class="main">

    <!-- Token banner (used by protected admin endpoints) -->
    <div class="token-banner" id="token-banner">
      <strong>Admin Token:</strong>
      <input type="password" id="admin-token" placeholder="X-Admin-Token (sadece admin endpointleri için gerekli)">
      <button onclick="saveToken()">Save</button>
    </div>

    <!-- =================== CONNECTION =================== -->
    <section class="page active" data-page="connection">
      <h1>Frontend Connection Info</h1>
      <p class="lead">Frontend developer'a vermen gereken her şey burada. Tek env değişkeni yeterli, gizli anahtar yok.</p>

      <div class="card featured">
        <h3>Environment Variable</h3>
        <pre id="env-line">loading…</pre>
        <div class="copy-row">
          <button onclick="copy('env-line')">📋 Copy .env line</button>
          <span class="muted">Frontend'in <code>.env</code> dosyasına bunu ekle.</span>
        </div>
      </div>

      <div class="card">
        <h3>API Base URL</h3>
        <pre id="api-base">loading…</pre>
        <div class="copy-row">
          <button onclick="copy('api-base')" class="ghost">📋 Copy URL</button>
        </div>
      </div>

      <div class="card-grid">
        <div class="stat">
          <div class="label">Backend Status</div>
          <div class="value"><span class="pill ok">RUNNING</span></div>
        </div>
        <div class="stat">
          <div class="label">Total Endpoints</div>
          <div class="value" id="ep-count">…</div>
        </div>
        <div class="stat">
          <div class="label">CORS Origins</div>
          <div class="value" style="font-size: 14px;" id="cors-count">…</div>
        </div>
      </div>

      <div class="card">
        <h3>Allowed CORS Origins</h3>
        <pre id="cors-list">loading…</pre>
      </div>

      <div class="card">
        <h3>Quick Test (frontend bağlantısı çalışıyor mu)</h3>
        <pre>fetch('<span id="test-url"></span>/api/v1/info').then(r => r.json()).then(console.log)</pre>
        <p class="muted">Bu satırı browser console'a yapıştır — JSON dönerse bağlantı sağlam.</p>
      </div>
    </section>

    <!-- =================== ENDPOINTS =================== -->
    <section class="page" data-page="endpoints">
      <h1>API Endpoints</h1>
      <p class="lead">Tüm endpoint'lerin tam listesi. 🔒 işareti olanlar <code>Authorization: Bearer &lt;token&gt;</code> header'ı ister.</p>

      <div id="endpoints-list">loading…</div>
    </section>

    <!-- =================== AUTH =================== -->
    <section class="page" data-page="auth">
      <h1>Auth & Token Usage</h1>
      <p class="lead">Login akışı, token saklama ve refresh.</p>

      <div class="card">
        <h3>1. Register / Login → Token al</h3>
        <pre>POST /api/v1/auth/login
Content-Type: application/json

{
  "email":    "user@example.com",
  "password": "secret123"
}</pre>
        <p style="margin-top:12px;">Response:</p>
        <pre>{
  "access_token":  "eyJhbGciOi...",
  "refresh_token": "v1.MzM...",
  "expires_in":    3600,
  "token_type":    "bearer",
  "user": { "id": "...", "username": "...", "email": "..." }
}</pre>
      </div>

      <div class="card">
        <h3>2. Token'ı sakla, her isteğe ekle</h3>
        <pre>fetch('/api/v1/me', {
  headers: {
    'Authorization': 'Bearer ' + accessToken
  }
})</pre>
      </div>

      <div class="card">
        <h3>3. Token expire olunca refresh et</h3>
        <pre>POST /api/v1/auth/refresh
Content-Type: application/json

{ "refresh_token": "v1.MzM..." }</pre>
        <p class="muted" style="margin-top:8px;">Response aynı şekilde — yeni access_token + refresh_token.</p>
      </div>

      <div class="card">
        <h3>4. Logout</h3>
        <pre>POST /api/v1/auth/logout
Authorization: Bearer &lt;access_token&gt;</pre>
      </div>
    </section>

    <!-- =================== HEALTH =================== -->
    <section class="page" data-page="health">
      <h1>System Health</h1>
      <p class="lead">Backend'in hangi servislere erişebildiğini canlı kontrol.</p>
      <div class="card-grid" id="health-cards">loading…</div>
      <div class="card">
        <h3>Raw Response</h3>
        <pre id="health-raw">loading…</pre>
        <button onclick="loadHealth()" class="ghost" style="margin-top:8px;">🔄 Refresh</button>
      </div>
    </section>

    <!-- =================== STATS =================== -->
    <section class="page" data-page="stats">
      <h1>Platform Statistics</h1>
      <p class="lead">Veritabanındaki gerçek sayılar. (Admin token gerekir.)</p>
      <div class="card-grid" id="stats-cards">loading…</div>
      <div class="card">
        <h3>Runtime</h3>
        <pre id="runtime-raw">loading…</pre>
        <button onclick="loadStats()" class="ghost" style="margin-top:8px;">🔄 Refresh</button>
      </div>
    </section>

    <!-- =================== QUEUE =================== -->
    <section class="page" data-page="queue">
      <h1>Video Processing Queue</h1>
      <p class="lead">FFmpeg worker'ın anlık durumu.</p>
      <div class="card-grid" id="queue-cards">loading…</div>
      <div class="card">
        <h3>Worker Detail</h3>
        <pre id="worker-raw">loading…</pre>
        <button onclick="loadQueue()" class="ghost" style="margin-top:8px;">🔄 Refresh</button>
      </div>
    </section>

    <!-- =================== STORAGE =================== -->
    <section class="page" data-page="storage">
      <h1>Storage Buckets</h1>
      <p class="lead">Supabase Storage bucket kullanımı.</p>
      <div class="card">
        <table>
          <thead><tr><th>Bucket</th><th>Files</th><th>Size</th></tr></thead>
          <tbody id="storage-table"><tr><td colspan="3">loading…</td></tr></tbody>
        </table>
        <button onclick="loadStorage()" class="ghost" style="margin-top:16px;">🔄 Refresh</button>
      </div>
    </section>

  </main>
</div>

<script>
// ---- Navigation ----
document.querySelectorAll('.nav a').forEach(a => {
  a.addEventListener('click', e => {
    e.preventDefault();
    const page = a.dataset.page;
    document.querySelectorAll('.nav a').forEach(x => x.classList.remove('active'));
    a.classList.add('active');
    document.querySelectorAll('.page').forEach(p => p.classList.remove('active'));
    document.querySelector('.page[data-page="' + page + '"]').classList.add('active');
    location.hash = page;
    if (page === 'health')   loadHealth();
    if (page === 'stats')    loadStats();
    if (page === 'queue')    loadQueue();
    if (page === 'storage')  loadStorage();
  });
});
if (location.hash) {
  const link = document.querySelector('.nav a[data-page="' + location.hash.slice(1) + '"]');
  if (link) link.click();
}

// ---- Token ----
function saveToken() {
  const v = document.getElementById('admin-token').value.trim();
  if (v) {
    localStorage.setItem('oroya_admin_token', v);
    document.getElementById('token-banner').classList.add('set');
    alert('Token saved.');
  }
}
const savedTok = localStorage.getItem('oroya_admin_token');
if (savedTok) {
  document.getElementById('admin-token').value = savedTok;
  document.getElementById('token-banner').classList.add('set');
}
function adminHeaders() {
  const t = localStorage.getItem('oroya_admin_token');
  return t ? { 'X-Admin-Token': t } : {};
}

// ---- Helpers ----
function copy(id) {
  const text = document.getElementById(id).textContent;
  navigator.clipboard.writeText(text).then(() => {
    const btns = document.querySelectorAll('button');
    btns.forEach(b => { if (b.textContent.includes(id) || true) {} });
  });
}
function bytes(n) {
  if (!n) return '0 B';
  const u = ['B','KB','MB','GB','TB'];
  let i = 0; while (n >= 1024 && i < u.length-1) { n /= 1024; i++; }
  return n.toFixed(1) + ' ' + u[i];
}
function fmt(j) { return JSON.stringify(j, null, 2); }

// ---- Connection page ----
fetch('/api/v1/info').then(r => r.json()).then(info => {
  const url = info.api_base_url;
  document.getElementById('env-line').textContent = 'PUBLIC_API_BASE_URL=' + url;
  document.getElementById('api-base').textContent = url;
  document.getElementById('test-url').textContent = url;
  let count = 0;
  Object.values(info.endpoints).forEach(g => count += Object.keys(g).length);
  document.getElementById('ep-count').textContent = count;
  document.getElementById('cors-count').textContent = (info.cors_allowed_origins || []).length;
  document.getElementById('cors-list').textContent = (info.cors_allowed_origins || []).join('\n');

  // Endpoints page
  let html = '';
  for (const [section, eps] of Object.entries(info.endpoints)) {
    html += '<div class="card"><h2 style="text-transform:capitalize;">' + section + '</h2><table><thead><tr><th>Method</th><th>Path</th><th>Description</th></tr></thead><tbody>';
    for (const [name, line] of Object.entries(eps)) {
      const m = line.match(/^(\w+)\s+(\S+)\s*(.*)$/);
      if (!m) continue;
      const lock = m[3].includes('auth') || m[3].includes('owner') ? '<span class="lock">🔒</span>' : '';
      html += '<tr><td><span class="method ' + m[1] + '">' + m[1] + '</span></td><td class="path">' + m[2] + lock + '</td><td class="muted">' + m[3] + '</td></tr>';
    }
    html += '</tbody></table></div>';
  }
  document.getElementById('endpoints-list').innerHTML = html;
});

// ---- Health page ----
function loadHealth() {
  fetch('/api/v1/admin/health').then(r => r.json()).then(j => {
    const services = [
      ['Database',  j.database],
      ['Storage',   j.storage],
      ['Worker',    j.worker]
    ];
    document.getElementById('health-cards').innerHTML = services.map(([name, ok]) =>
      '<div class="stat"><div class="label">' + name + '</div><div class="value"><span class="pill ' + (ok ? 'ok' : 'bad') + '">' + (ok ? 'UP' : 'DOWN') + '</span></div></div>'
    ).join('') +
    '<div class="stat"><div class="label">Overall</div><div class="value"><span class="pill ' + (j.status === 'ok' ? 'ok' : 'warn') + '">' + (j.status || '?').toUpperCase() + '</span></div></div>' +
    '<div class="stat"><div class="label">Uptime</div><div class="value" style="font-size:18px;">' + (j.uptime_s || 0) + 's</div></div>';
    document.getElementById('health-raw').textContent = fmt(j);
  });
}

// ---- Stats page ----
function loadStats() {
  fetch('/api/v1/admin/stats', { headers: adminHeaders() }).then(r => r.json()).then(j => {
    const d = j.data || {};
    const cards = [
      ['Users',           d.users_total],
      ['Videos',          d.videos_total],
      ['Processing',      d.videos_processing],
      ['Ready',           d.videos_ready],
      ['Failed',          d.videos_failed],
      ['Comments',        d.comments_total],
      ['Views',           d.views_total],
      ['Likes',           d.likes_total]
    ];
    document.getElementById('stats-cards').innerHTML = cards.map(([l, v]) =>
      '<div class="stat"><div class="label">' + l + '</div><div class="value">' + (v ?? 0) + '</div></div>'
    ).join('');
    document.getElementById('runtime-raw').textContent = fmt(j.runtime || {});
  });
}

// ---- Queue page ----
function loadQueue() {
  fetch('/api/v1/admin/worker-status', { headers: adminHeaders() }).then(r => r.json()).then(j => {
    const cards = [
      ['Status',       j.running ? '<span class="pill ok">RUNNING</span>' : '<span class="pill bad">STOPPED</span>'],
      ['Concurrency',  j.concurrency || 0],
      ['Pending',      j.pending || 0],
      ['Processing',   j.processing || 0],
      ['Completed',    j.completed || 0],
      ['Failed',       j.failed || 0]
    ];
    document.getElementById('queue-cards').innerHTML = cards.map(([l, v]) =>
      '<div class="stat"><div class="label">' + l + '</div><div class="value" style="font-size:22px;">' + v + '</div></div>'
    ).join('');
    document.getElementById('worker-raw').textContent = fmt(j);
  });
}

// ---- Storage page ----
function loadStorage() {
  fetch('/api/v1/admin/storage-status', { headers: adminHeaders() }).then(r => r.json()).then(j => {
    const buckets = j.buckets || {};
    const rows = Object.entries(buckets).map(([label, b]) =>
      '<tr><td><strong>' + (b.bucket || label) + '</strong></td><td class="num">' + (b.files || 0) + '</td><td class="num">' + bytes(b.size_bytes || 0) + '</td></tr>'
    ).join('');
    document.getElementById('storage-table').innerHTML = rows || '<tr><td colspan="3" class="muted">No data</td></tr>';
  });
}

// initial loads for visible page
loadHealth();
</script>
</body>
</html>`
