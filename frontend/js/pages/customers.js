App.registerPage('customers', async () => {
  render('<div class="empty-msg">Loading customers...</div>');

  let allCustomers = [];
  await loadCustomers();

  async function loadCustomers() {
    try {
      const data = await App.api('/customers/');
      allCustomers = Array.isArray(data) ? data : [];
      renderList(allCustomers);
    } catch (err) { render(`<div class="empty-msg">${err.message}</div>`); }
  }

  function renderList(list) {
    render(`
      <div class="toolbar">
        <button class="btn" id="new-cust-btn">+ New Customer</button>
        <input type="search" id="cust-search" placeholder="Search customers...">
      </div>
      <div class="table-wrap">
        <table>
          <thead><tr><th>ID</th><th>Name</th><th>Phone</th><th>Email</th><th>Created</th><th>Actions</th></tr></thead>
          <tbody id="cust-tbody">
            ${list.map(c => `<tr>
              <td>${c.id}</td>
              <td>${c.name}</td>
              <td>${c.phone || '-'}</td>
              <td>${c.email || '-'}</td>
              <td>${fmtDate(c.created_at)}</td>
              <td class="inline-flex">
                <button class="btn btn-sm btn-info" onclick="viewCustOrders(${c.id}, '${c.name}')">Orders</button>
                <button class="btn btn-sm" onclick="showCustForm(${c.id})">Edit</button>
                <button class="btn btn-sm btn-danger" onclick="deleteCust(${c.id})">Del</button>
              </td>
            </tr>`).join('')}
          </tbody>
        </table>
      </div>
    `);

    document.getElementById('new-cust-btn').addEventListener('click', () => showCustForm());
    document.getElementById('cust-search').addEventListener('input', e => {
      const q = e.target.value.toLowerCase();
      const filtered = allCustomers.filter(c =>
        c.name.toLowerCase().includes(q) ||
        (c.email || '').toLowerCase().includes(q) ||
        (c.phone || '').includes(q)
      );
      document.getElementById('cust-tbody').innerHTML = filtered.map(c => `<tr>
        <td>${c.id}</td>
        <td>${c.name}</td>
        <td>${c.phone || '-'}</td>
        <td>${c.email || '-'}</td>
        <td>${fmtDate(c.created_at)}</td>
        <td class="inline-flex">
          <button class="btn btn-sm btn-info" onclick="viewCustOrders(${c.id}, '${c.name}')">Orders</button>
          <button class="btn btn-sm" onclick="showCustForm(${c.id})">Edit</button>
          <button class="btn btn-sm btn-danger" onclick="deleteCust(${c.id})">Del</button>
        </td>
      </tr>`).join('');
    });
  }

  window.showCustForm = async function(id) {
    let cust = { name: '', phone: '', email: '' };
    if (id) {
      try { cust = await App.api('/customers/' + id); } catch { return App.toast('Failed to load', 'error'); }
    }
    App.modal(id ? 'Edit Customer' : 'New Customer', `
      <form id="cust-form">
        <div class="form-group"><label>Name</label><input id="cf-name" value="${cust.name}" required></div>
        <div class="form-row">
          <div class="form-group"><label>Phone</label><input id="cf-phone" value="${cust.phone || ''}"></div>
          <div class="form-group"><label>Email</label><input type="email" id="cf-email" value="${cust.email || ''}"></div>
        </div>
        <button type="submit" class="btn btn-block">${id ? 'Update' : 'Create'}</button>
      </form>
    `);
    document.getElementById('cust-form').addEventListener('submit', async e => {
      e.preventDefault();
      try {
        await App.api('/customers/' + (id || ''), {
          method: id ? 'PUT' : 'POST',
          body: {
            name: document.getElementById('cf-name').value,
            phone: document.getElementById('cf-phone').value,
            email: document.getElementById('cf-email').value,
          }
        });
        App.closeModal();
        App.toast(id ? 'Customer updated' : 'Customer created', 'success');
        loadCustomers();
      } catch (err) { App.toast(err.message, 'error'); }
    });
  };

  window.deleteCust = async function(id) {
    if (!confirm('Delete this customer?')) return;
    try {
      await App.api('/customers/' + id, { method: 'DELETE' });
      App.toast('Deleted', 'success');
      loadCustomers();
    } catch (err) { App.toast(err.message, 'error'); }
  };

  window.viewCustOrders = async function(id, name) {
    try {
      const orders = await App.api('/customers/' + id + '/orders');
      const list = Array.isArray(orders) ? orders : [];
      App.modal('Orders for ' + name, list.length === 0 ? '<p class="text-muted">No orders</p>' : `
        <table class="sub-table">
          <thead><tr><th>#</th><th>Type</th><th>Status</th><th>Total</th><th>Date</th></tr></thead>
          <tbody>
            ${list.map(o => `<tr>
              <td class="text-accent">#${o.id}</td>
              <td>${o.order_type}</td>
              <td>${statusBadge(o.status)}</td>
              <td>${fmtPrice(o.total_price)}</td>
              <td>${fmtDate(o.created_at)}</td>
            </tr>`).join('')}
          </tbody>
        </table>
      `);
    } catch (err) { App.toast(err.message, 'error'); }
  };
});
