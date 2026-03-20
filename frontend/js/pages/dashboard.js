App.registerPage('dashboard', async () => {
  render('<div class="loading">Loading dashboard...</div>');

  const role = App.getRole();

  try {
    if (role === 'preparation') {
      await renderPreparationDashboard();
    } else if (role === 'accueil') {
      await renderAccueilDashboard();
    } else {
      await renderAdminDashboard();
    }
  } catch (err) {
    render(`<div class="empty-msg">Error loading dashboard: ${esc(err.message)}</div>`);
  }

  // ===== ADMIN DASHBOARD =====
  async function renderAdminDashboard() {
    const [products, categories, menus, customers, orders] = await Promise.all([
      App.api('/products/'),
      App.api('/categories/'),
      App.api('/menus/'),
      App.api('/customers/'),
      App.api('/orders/'),
    ]);

    const prodList = Array.isArray(products) ? products : [];
    const catList  = Array.isArray(categories) ? categories : [];
    const menuList = Array.isArray(menus) ? menus : [];
    const custList = Array.isArray(customers) ? customers : [];
    const ordList  = Array.isArray(orders) ? orders : [];

    const openOrders = ordList.filter(o => o.status === 'pending' || o.status === 'preparing');
    const recent = ordList.slice(0, 5);

    render(`
      <div class="stat-grid">
        <div class="stat-card">
          <div class="stat-value">${prodList.length}</div>
          <div class="stat-label">Products</div>
        </div>
        <div class="stat-card">
          <div class="stat-value">${catList.length}</div>
          <div class="stat-label">Categories</div>
        </div>
        <div class="stat-card">
          <div class="stat-value">${menuList.length}</div>
          <div class="stat-label">Menus</div>
        </div>
        <div class="stat-card">
          <div class="stat-value">${custList.length}</div>
          <div class="stat-label">Customers</div>
        </div>
        <div class="stat-card">
          <div class="stat-value">${openOrders.length}</div>
          <div class="stat-label">Open Orders</div>
        </div>
        <div class="stat-card">
          <div class="stat-value">${ordList.length}</div>
          <div class="stat-label">Total Orders</div>
        </div>
      </div>

      <div class="card">
        <div class="section-title">Recent Orders</div>
        ${recent.length === 0 ? '<p class="text-muted">No orders yet</p>' : `
        <div class="table-wrap">
          <table>
            <thead>
              <tr><th>#</th><th>Type</th><th>Status</th><th>Total</th><th>Created</th></tr>
            </thead>
            <tbody>
              ${recent.map(o => `
                <tr>
                  <td class="text-accent">#${o.id}</td>
                  <td>${esc(o.order_type)}</td>
                  <td>${statusBadge(o.status)}</td>
                  <td>${fmtPrice(o.total_price)}</td>
                  <td>${fmtDate(o.created_at)}</td>
                </tr>
              `).join('')}
            </tbody>
          </table>
        </div>`}
      </div>
    `);
  }

  // ===== PREPARATION DASHBOARD =====
  async function renderPreparationDashboard() {
    const orders = await App.api('/orders/');
    const ordList = Array.isArray(orders) ? orders : [];

    const prepOrders = ordList
      .filter(o => o.status === 'pending' || o.status === 'preparing')
      .sort((a, b) => {
        if (!a.scheduled_time && !b.scheduled_time) return 0;
        if (!a.scheduled_time) return 1;
        if (!b.scheduled_time) return -1;
        return new Date(a.scheduled_time) - new Date(b.scheduled_time);
      });

    const pending = prepOrders.filter(o => o.status === 'pending');
    const preparing = prepOrders.filter(o => o.status === 'preparing');

    render(`
      <div class="stat-grid">
        <div class="stat-card">
          <div class="stat-value">${pending.length}</div>
          <div class="stat-label">Pending</div>
        </div>
        <div class="stat-card">
          <div class="stat-value">${preparing.length}</div>
          <div class="stat-label">Preparing</div>
        </div>
      </div>

      <div class="kanban">
        <div class="kanban-col">
          <div class="kanban-col-header">Pending <span class="count">${pending.length}</span></div>
          <div class="kanban-col-body">
            ${pending.length === 0 ? '<p class="text-muted" style="text-align:center;font-size:12px;">No orders</p>' :
              pending.map(o => renderPrepCard(o)).join('')}
          </div>
        </div>
        <div class="kanban-col">
          <div class="kanban-col-header">Preparing <span class="count">${preparing.length}</span></div>
          <div class="kanban-col-body">
            ${preparing.length === 0 ? '<p class="text-muted" style="text-align:center;font-size:12px;">No orders</p>' :
              preparing.map(o => renderPrepCard(o)).join('')}
          </div>
        </div>
      </div>
    `);
  }

  function renderPrepCard(o) {
    const items = (o.order_items || []).map(it => {
      const name = it.product ? it.product.name : (it.menu ? it.menu.name : 'Item');
      return `${it.quantity}x ${esc(name)}`;
    }).join(', ');
    const notes = o.notes ? esc(o.notes.length > 60 ? o.notes.slice(0, 60) + '...' : o.notes) : '';
    const action = o.status === 'pending'
      ? `<button class="btn btn-sm btn-info" onclick="updateOrderStatus(${o.id},'preparing')">Start</button>`
      : `<button class="btn btn-sm btn-success" onclick="updateOrderStatus(${o.id},'prepared')">Ready</button>`;

    return `<div class="kanban-card">
      <div class="order-id">#${o.id}</div>
      <div class="order-meta">${esc(o.order_type)} | ${o.order_items ? o.order_items.length : 0} items</div>
      ${items ? `<div class="order-meta text-muted" style="font-size:11px;">${items}</div>` : ''}
      ${o.scheduled_time ? `<div class="order-meta text-muted" style="font-size:11px;">Scheduled: ${fmtDate(o.scheduled_time)}</div>` : ''}
      ${notes ? `<div class="order-meta text-muted" style="font-size:11px;font-style:italic;">${notes}</div>` : ''}
      <div class="order-actions">${action}</div>
    </div>`;
  }

  // ===== ACCUEIL DASHBOARD =====
  async function renderAccueilDashboard() {
    const [orders, customers] = await Promise.all([
      App.api('/orders/'),
      App.api('/customers/'),
    ]);
    const ordList = Array.isArray(orders) ? orders : [];
    const custList = Array.isArray(customers) ? customers : [];

    const pending = ordList.filter(o => o.status === 'pending').length;
    const preparing = ordList.filter(o => o.status === 'preparing').length;
    const prepared = ordList.filter(o => o.status === 'prepared').length;
    const deliveredToday = ordList.filter(o => {
      if (o.status !== 'delivered') return false;
      return new Date(o.updated_at).toDateString() === new Date().toDateString();
    }).length;

    const readyOrders = ordList.filter(o => o.status === 'prepared');

    render(`
      <div class="stat-grid">
        <div class="stat-card">
          <div class="stat-value">${pending}</div>
          <div class="stat-label">Pending</div>
        </div>
        <div class="stat-card">
          <div class="stat-value">${preparing}</div>
          <div class="stat-label">Preparing</div>
        </div>
        <div class="stat-card">
          <div class="stat-value">${prepared}</div>
          <div class="stat-label">Ready to Deliver</div>
        </div>
        <div class="stat-card">
          <div class="stat-value">${deliveredToday}</div>
          <div class="stat-label">Delivered Today</div>
        </div>
      </div>

      <div style="margin-bottom:20px;">
        <button class="btn" onclick="App.navigate('orders')">+ New Order</button>
      </div>

      ${readyOrders.length > 0 ? `
      <div class="card">
        <div class="section-title">Ready for Delivery</div>
        <div class="table-wrap">
          <table>
            <thead><tr><th>#</th><th>Type</th><th>Customer</th><th>Total</th><th>Actions</th></tr></thead>
            <tbody>
              ${readyOrders.map(o => `<tr>
                <td class="text-accent">#${o.id}</td>
                <td>${esc(o.order_type)}</td>
                <td>${o.customer ? esc(o.customer.name) : 'Walk-in'}</td>
                <td>${fmtPrice(o.total_price)}</td>
                <td><button class="btn btn-sm btn-success" onclick="updateOrderStatus(${o.id},'delivered')">Deliver</button></td>
              </tr>`).join('')}
            </tbody>
          </table>
        </div>
      </div>` : '<div class="card"><p class="text-muted">No orders ready for delivery</p></div>'}
    `);
  }

  // Global handlers for dashboard action buttons
  window.updateOrderStatus = async function(id, status) {
    try {
      await App.api('/orders/' + id + '/status', { method: 'PATCH', body: { status } });
      App.toast('Status updated to ' + status, 'success');
      App.route();
    } catch (err) { App.toast(err.message, 'error'); }
  };
});
