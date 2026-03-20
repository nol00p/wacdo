/* ===== Global App Object ===== */
const App = {
  API: 'http://localhost:8000',

  // --- Auth ---
  getToken()  { return localStorage.getItem('wacdo_token'); },
  setToken(t) { localStorage.setItem('wacdo_token', t); },
  clearToken(){ localStorage.removeItem('wacdo_token'); },
  isLoggedIn(){ return !!this.getToken(); },

  // Decode JWT payload to read claims (UserID, RoleName)
  getTokenPayload() {
    const token = this.getToken();
    if (!token) return null;
    try {
      const payload = token.split('.')[1];
      return JSON.parse(atob(payload));
    } catch { return null; }
  },
  getRole() {
    const p = this.getTokenPayload();
    return p ? (p.RoleName || '') : '';
  },

  // --- API helper ---
  async api(path, opts = {}) {
    const url = this.API + path;
    const headers = { 'Content-Type': 'application/json' };
    const token = this.getToken();
    if (token) headers['Authorization'] = 'Bearer ' + token;

    const res = await fetch(url, {
      method: opts.method || 'GET',
      headers,
      body: opts.body ? JSON.stringify(opts.body) : undefined,
    });

    if (res.status === 401) {
      this.clearToken();
      this.toast('Session expired. Please log in again.', 'error');
      this.navigate('login');
      throw new Error('Session expired');
    }

    const text = await res.text();
    let data;
    try { data = JSON.parse(text); } catch { data = text; }

    if (!res.ok) {
      const msg = (data && typeof data === 'object' && data.error) ? data.error : 'Request failed';
      throw new Error(msg);
    }
    return data;
  },

  // --- Toast ---
  toast(msg, type = 'info') {
    const el = document.createElement('div');
    el.className = 'toast toast-' + type;
    el.textContent = msg;
    document.getElementById('toast-container').appendChild(el);
    setTimeout(() => el.remove(), 3000);
  },

  // --- Modal ---
  modal(title, contentHTML) {
    document.getElementById('modal-title').textContent = title;
    document.getElementById('modal-body').innerHTML = contentHTML;
    document.getElementById('modal-overlay').classList.remove('hidden');
  },
  closeModal() {
    document.getElementById('modal-overlay').classList.add('hidden');
  },

  // --- Router ---
  pages: {},
  registerPage(name, fn) { this.pages[name] = fn; },

  navigate(page) {
    window.location.hash = page;
  },

  async route() {
    const hash = (window.location.hash || '#login').slice(1).split('/')[0];

    if (!this.isLoggedIn() && hash !== 'login') {
      window.location.hash = 'login';
      return;
    }

    if (this.isLoggedIn() && hash === 'login') {
      window.location.hash = 'dashboard';
      return;
    }

    const shell = document.getElementById('app');
    const sidebar = document.getElementById('sidebar');
    const topbar = document.getElementById('topbar');

    if (hash === 'login') {
      sidebar.classList.add('hidden');
      topbar.classList.add('hidden');
      document.getElementById('content').style.padding = '0';
    } else {
      sidebar.classList.remove('hidden');
      topbar.classList.remove('hidden');
      document.getElementById('content').style.padding = '';

      // Show user role in topbar
      const role = this.getRole();
      document.getElementById('user-label').textContent = role ? role : '';

      // Hide sidebar links based on role
      const access = {
        admin:       ['dashboard', 'products', 'menus', 'orders', 'customers', 'users', 'privacy'],
        accueil:     ['dashboard', 'orders', 'customers', 'privacy'],
        preparation: ['dashboard', 'orders', 'privacy'],
      };
      const allowed = access[role] || [];
      document.querySelectorAll('#sidebar-nav a').forEach(a => {
        const page = a.dataset.page;
        a.style.display = (allowed.length === 0 || allowed.includes(page)) ? '' : 'none';
      });
    }

    // Update active nav
    document.querySelectorAll('#sidebar-nav a').forEach(a => {
      a.classList.toggle('active', a.dataset.page === hash);
    });

    // Page title
    document.getElementById('page-title').textContent =
      hash.charAt(0).toUpperCase() + hash.slice(1);

    const pageFn = this.pages[hash];
    if (pageFn) {
      await pageFn();
    } else {
      document.getElementById('content').innerHTML =
        '<div class="empty-msg">Page not found</div>';
    }
  },

  init() {
    // Modal close
    document.getElementById('modal-close').addEventListener('click', () => this.closeModal());
    document.getElementById('modal-overlay').addEventListener('click', e => {
      if (e.target === e.currentTarget) this.closeModal();
    });

    // Logout
    document.getElementById('logout-btn').addEventListener('click', () => {
      this.clearToken();
      this.navigate('login');
    });

    // Sidebar toggle (mobile)
    document.getElementById('sidebar-toggle').addEventListener('click', () => {
      document.getElementById('sidebar').classList.toggle('open');
    });

    // Close sidebar on nav click (mobile)
    document.querySelectorAll('#sidebar-nav a').forEach(a => {
      a.addEventListener('click', () => {
        document.getElementById('sidebar').classList.remove('open');
      });
    });

    // Route on hash change
    window.addEventListener('hashchange', () => this.route());
    this.route();
  }
};

/* ===== Helper: render HTML into content ===== */
function render(html) {
  document.getElementById('content').innerHTML = html;
}

/* ===== Helper: escape HTML to prevent XSS ===== */
function esc(s) {
  if (!s) return '';
  const div = document.createElement('div');
  div.textContent = String(s);
  return div.innerHTML;
}

/* ===== Helper: format price ===== */
function fmtPrice(v) {
  return (Number(v) || 0).toFixed(2) + ' €';
}

/* ===== Helper: format date ===== */
function fmtDate(d) {
  if (!d) return '-';
  return new Date(d).toLocaleString();
}

/* ===== Helper: status badge ===== */
function statusBadge(s) {
  return `<span class="badge badge-${s}">${s}</span>`;
}

/* ===== Helper: availability badge ===== */
function availBadge(v) {
  return v
    ? '<span class="badge badge-available">available</span>'
    : '<span class="badge badge-unavailable">unavailable</span>';
}

/* ===== Boot ===== */
document.addEventListener('DOMContentLoaded', () => App.init());
