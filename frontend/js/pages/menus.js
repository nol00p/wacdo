App.registerPage('menus', async () => {
  render('<div class="empty-msg">Loading menus...</div>');

  let allProducts = [];
  try { allProducts = await App.api('/products/'); if (!Array.isArray(allProducts)) allProducts = []; } catch {}

  await loadMenus();

  async function loadMenus() {
    try {
      const menus = await App.api('/menus/');
      const list = Array.isArray(menus) ? menus : [];

      render(`
        <div class="toolbar">
          <button class="btn" id="new-menu-btn">+ New Menu</button>
        </div>
        <div class="table-wrap">
          <table>
            <thead><tr><th>ID</th><th>Name</th><th>Description</th><th>Price</th><th>Available</th><th>Actions</th></tr></thead>
            <tbody id="menu-tbody">
              ${list.map(m => `
                <tr>
                  <td>${m.id}</td>
                  <td>${m.name}</td>
                  <td class="text-muted">${m.description || '-'}</td>
                  <td>${fmtPrice(m.price)}</td>
                  <td>
                    <label class="toggle">
                      <input type="checkbox" ${m.is_available ? 'checked' : ''} onchange="toggleMenuAvail(${m.id})">
                      <span class="toggle-slider"></span>
                    </label>
                  </td>
                  <td class="inline-flex">
                    <button class="btn btn-sm btn-info" onclick="expandMenu(${m.id}, this)">Products</button>
                    <button class="btn btn-sm" onclick="showMenuForm(${m.id})">Edit</button>
                    <button class="btn btn-sm btn-danger" onclick="deleteMenu(${m.id})">Del</button>
                  </td>
                </tr>
                <tr class="expand-row hidden" id="menu-expand-${m.id}">
                  <td colspan="6"><div id="menu-products-${m.id}"></div></td>
                </tr>
              `).join('')}
            </tbody>
          </table>
        </div>
      `);

      document.getElementById('new-menu-btn').addEventListener('click', () => showMenuForm());
    } catch (err) { render(`<div class="empty-msg">${err.message}</div>`); }
  }

  window.toggleMenuAvail = async function(id) {
    try {
      await App.api('/menus/' + id + '/availability', { method: 'PATCH' });
      App.toast('Availability toggled', 'success');
    } catch (err) { App.toast(err.message, 'error'); loadMenus(); }
  };

  window.showMenuForm = async function(id) {
    let menu = { name: '', description: '', price: 0 };
    if (id) {
      try { menu = await App.api('/menus/' + id); } catch { return App.toast('Failed to load', 'error'); }
    }
    App.modal(id ? 'Edit Menu' : 'New Menu', `
      <form id="menu-form">
        <div class="form-group"><label>Name</label><input id="mf-name" value="${menu.name}" required></div>
        <div class="form-group"><label>Description</label><input id="mf-desc" value="${menu.description || ''}"></div>
        <div class="form-group"><label>Price</label><input type="number" step="0.01" id="mf-price" value="${menu.price}" required></div>
        <button type="submit" class="btn btn-block">${id ? 'Update' : 'Create'}</button>
      </form>
    `);
    document.getElementById('menu-form').addEventListener('submit', async e => {
      e.preventDefault();
      try {
        await App.api('/menus/' + (id || ''), {
          method: id ? 'PUT' : 'POST',
          body: {
            name: document.getElementById('mf-name').value,
            description: document.getElementById('mf-desc').value,
            price: Number(document.getElementById('mf-price').value),
          }
        });
        App.closeModal();
        App.toast(id ? 'Menu updated' : 'Menu created', 'success');
        loadMenus();
      } catch (err) { App.toast(err.message, 'error'); }
    });
  };

  window.deleteMenu = async function(id) {
    if (!confirm('Delete this menu?')) return;
    try {
      await App.api('/menus/' + id, { method: 'DELETE' });
      App.toast('Deleted', 'success');
      loadMenus();
    } catch (err) { App.toast(err.message, 'error'); }
  };

  window.expandMenu = async function(id, btn) {
    const row = document.getElementById('menu-expand-' + id);
    const isHidden = row.classList.contains('hidden');
    // collapse all
    document.querySelectorAll('.expand-row').forEach(r => r.classList.add('hidden'));
    if (isHidden) {
      row.classList.remove('hidden');
      loadMenuProducts(id);
    }
  };

  async function loadMenuProducts(menuId) {
    const el = document.getElementById('menu-products-' + menuId);
    el.innerHTML = 'Loading...';
    try {
      const items = await App.api('/menus/' + menuId + '/products/');
      const list = Array.isArray(items) ? items : [];

      el.innerHTML = `
        <div class="flex-between mb-8">
          <strong>Menu Products</strong>
          <div class="inline-flex">
            <select id="mp-prod-${menuId}">
              <option value="">Select product</option>
              ${allProducts.map(p => `<option value="${p.id}">${p.name}</option>`).join('')}
            </select>
            <input type="number" id="mp-qty-${menuId}" value="1" min="1" style="width:50px;padding:4px;background:var(--bg-input);border:1px solid var(--border);border-radius:4px;color:var(--text);" placeholder="Qty">
            <input type="number" id="mp-order-${menuId}" value="0" min="0" style="width:50px;padding:4px;background:var(--bg-input);border:1px solid var(--border);border-radius:4px;color:var(--text);" placeholder="Order">
            <label style="display:flex;align-items:center;gap:4px;font-size:12px;color:var(--text-muted);cursor:pointer;"><input type="checkbox" id="mp-opt-${menuId}"> Optional</label>
            <button class="btn btn-sm" onclick="addMenuProduct(${menuId})">Add</button>
          </div>
        </div>
        ${list.length === 0 ? '<p class="text-muted">No products in this menu</p>' : `
          <table class="sub-table">
            <thead><tr><th>Product</th><th>Qty</th><th>Optional</th><th>Actions</th></tr></thead>
            <tbody>
              ${list.map(mp => {
                const prod = allProducts.find(p => p.id === mp.product_id);
                const prodLabel = prod ? prod.name + ' (' + fmtPrice(prod.price) + ')' : 'Product #' + mp.product_id;
                return `<tr>
                <td>${prodLabel}</td>
                <td>${mp.quantity}</td>
                <td>${mp.is_optional ? 'Yes' : 'No'}</td>
                <td><button class="btn btn-sm btn-danger" onclick="removeMenuProduct(${mp.id}, ${menuId})">Remove</button></td>
              </tr>`;
              }).join('')}
            </tbody>
          </table>`}
      `;
    } catch (err) { el.innerHTML = `<p class="text-muted">${err.message}</p>`; }
  }

  window.addMenuProduct = async function(menuId) {
    const prodId = document.getElementById('mp-prod-' + menuId).value;
    const qty = document.getElementById('mp-qty-' + menuId).value;
    const displayOrder = document.getElementById('mp-order-' + menuId).value;
    const isOptional = document.getElementById('mp-opt-' + menuId).checked;
    if (!prodId) return App.toast('Select a product', 'error');
    try {
      await App.api('/menus/' + menuId + '/products/', {
        method: 'POST',
        body: { product_id: Number(prodId), quantity: Number(qty), display_order: Number(displayOrder), is_optional: isOptional }
      });
      App.toast('Product added', 'success');
      loadMenuProducts(menuId);
    } catch (err) { App.toast(err.message, 'error'); }
  };

  window.removeMenuProduct = async function(mpId, menuId) {
    if (!confirm('Remove this product?')) return;
    try {
      await App.api('/menus/products/' + mpId, { method: 'DELETE' });
      App.toast('Removed', 'success');
      loadMenuProducts(menuId);
    } catch (err) { App.toast(err.message, 'error'); }
  };
});
