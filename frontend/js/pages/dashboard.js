App.registerPage('dashboard', async () => {
  render('<div class="empty-msg">Loading dashboard...</div>');

  try {
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
                  <td>${o.order_type}</td>
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
  } catch (err) {
    render(`<div class="empty-msg">Error loading dashboard: ${err.message}</div>`);
  }
});
