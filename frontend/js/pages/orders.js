App.registerPage('orders', async () => {
  render('<div class="empty-msg">Loading orders...</div>');

  let allProducts = [], allMenus = [], allCustomers = [], allOptions = {};

  try {
    [allProducts, allMenus, allCustomers] = await Promise.all([
      App.api('/products/').then(r => Array.isArray(r) ? r : []),
      App.api('/menus/').then(r => Array.isArray(r) ? r : []),
      App.api('/customers/').then(r => Array.isArray(r) ? r : []),
    ]);
  } catch {}

  const STATUSES = ['pending', 'preparing', 'prepared', 'delivered', 'cancelled'];

  render(`
    <div class="toolbar">
      <button class="btn" id="new-order-btn">+ New Order</button>
      <div class="tabs" style="border:none;margin:0;">
        <button class="tab-btn active" data-filter="">All</button>
        ${STATUSES.map(s => `<button class="tab-btn" data-filter="${s}">${s}</button>`).join('')}
      </div>
    </div>
    <div id="orders-view"></div>
  `);

  document.getElementById('new-order-btn').addEventListener('click', showNewOrderForm);
  document.querySelectorAll('.toolbar .tab-btn').forEach(btn => {
    btn.addEventListener('click', () => {
      document.querySelectorAll('.toolbar .tab-btn').forEach(b => b.classList.remove('active'));
      btn.classList.add('active');
      loadOrders(btn.dataset.filter);
    });
  });

  loadOrders('');

  async function loadOrders(statusFilter) {
    const el = document.getElementById('orders-view');
    el.innerHTML = 'Loading...';
    try {
      const url = statusFilter ? '/orders/?status=' + statusFilter : '/orders/';
      const orders = await App.api(url);
      const list = Array.isArray(orders) ? orders : [];

      if (!statusFilter) {
        // Kanban view
        el.innerHTML = `<div class="kanban">
          ${STATUSES.map(s => {
            const col = list.filter(o => o.status === s);
            return `<div class="kanban-col">
              <div class="kanban-col-header">${s} <span class="count">${col.length}</span></div>
              <div class="kanban-col-body">
                ${col.length === 0 ? '<p class="text-muted" style="text-align:center;font-size:12px;">No orders</p>' :
                  col.map(o => renderKanbanCard(o)).join('')}
              </div>
            </div>`;
          }).join('')}
        </div>`;
      } else {
        // Table view
        el.innerHTML = `<div class="table-wrap"><table>
          <thead><tr><th>#</th><th>Type</th><th>Customer</th><th>Status</th><th>Total</th><th>Items</th><th>Created</th><th>Actions</th></tr></thead>
          <tbody>
            ${list.map(o => `<tr>
              <td class="text-accent">#${o.id}</td>
              <td>${o.order_type}</td>
              <td>${o.customer ? o.customer.name : '-'}</td>
              <td>${statusBadge(o.status)}</td>
              <td>${fmtPrice(o.total_price)}</td>
              <td>${o.order_items ? o.order_items.length : 0}</td>
              <td>${fmtDate(o.created_at)}</td>
              <td class="inline-flex">${orderActionButtons(o)}</td>
            </tr>`).join('')}
          </tbody>
        </table></div>`;
      }
    } catch (err) { el.innerHTML = `<div class="empty-msg">${err.message}</div>`; }
  }

  function renderKanbanCard(o) {
    const notesSnippet = o.notes ? (o.notes.length > 50 ? o.notes.slice(0, 50) + '...' : o.notes) : '';
    return `<div class="kanban-card">
      <div class="order-id">#${o.id}</div>
      <div class="order-meta">${o.order_type} | ${fmtPrice(o.total_price)}</div>
      <div class="order-meta">${o.customer ? o.customer.name : 'Walk-in'}</div>
      <div class="order-meta">${o.order_items ? o.order_items.length : 0} items</div>
      ${notesSnippet ? `<div class="order-meta text-muted" style="font-size:11px;font-style:italic;">${notesSnippet}</div>` : ''}
      <div class="order-actions">${orderActionButtons(o)}</div>
    </div>`;
  }

  function orderActionButtons(o) {
    const btns = [];
    if (o.status === 'pending') {
      btns.push(`<button class="btn btn-sm btn-info" onclick="updateOrderStatus(${o.id},'preparing')">Prepare</button>`);
      btns.push(`<button class="btn btn-sm btn-danger" onclick="cancelOrder(${o.id})">Cancel</button>`);
    } else if (o.status === 'preparing') {
      btns.push(`<button class="btn btn-sm btn-success" onclick="updateOrderStatus(${o.id},'prepared')">Ready</button>`);
    } else if (o.status === 'prepared') {
      btns.push(`<button class="btn btn-sm btn-success" onclick="updateOrderStatus(${o.id},'delivered')">Deliver</button>`);
    }
    btns.push(`<button class="btn btn-sm btn-outline" onclick="viewOrderDetail(${o.id})">View</button>`);
    return btns.join('');
  }

  window.updateOrderStatus = async function(id, status) {
    try {
      await App.api('/orders/' + id + '/status', { method: 'PATCH', body: { status } });
      App.toast('Status updated to ' + status, 'success');
      const activeFilter = document.querySelector('.toolbar .tab-btn.active').dataset.filter;
      loadOrders(activeFilter);
    } catch (err) { App.toast(err.message, 'error'); }
  };

  window.cancelOrder = async function(id) {
    if (!confirm('Cancel this order?')) return;
    try {
      await App.api('/orders/' + id + '/cancel', { method: 'PATCH' });
      App.toast('Order cancelled', 'success');
      const activeFilter = document.querySelector('.toolbar .tab-btn.active').dataset.filter;
      loadOrders(activeFilter);
    } catch (err) { App.toast(err.message, 'error'); }
  };

  window.viewOrderDetail = async function(id) {
    try {
      const o = await App.api('/orders/' + id);
      App.modal('Order #' + o.id, `
        <div class="mb-16">
          <p><strong>Type:</strong> ${o.order_type}</p>
          <p><strong>Status:</strong> ${statusBadge(o.status)}</p>
          <p><strong>Customer:</strong> ${o.customer ? o.customer.name : 'Walk-in'}</p>
          <p><strong>Notes:</strong> ${o.notes || '-'}</p>
          <p><strong>Total:</strong> <span class="text-accent">${fmtPrice(o.total_price)}</span></p>
          <p><strong>Created:</strong> ${fmtDate(o.created_at)}</p>
        </div>
        <div class="section-title">Items</div>
        <table class="sub-table">
          <thead><tr><th>Item</th><th>Qty</th><th>Unit Price</th><th>Total</th></tr></thead>
          <tbody>
            ${(o.order_items || []).map(it => {
              let itemName;
              if (it.product_id) {
                const prod = allProducts.find(p => p.id === it.product_id);
                itemName = prod ? prod.name : 'Product #' + it.product_id;
              } else {
                const menu = allMenus.find(m => m.id === it.menu_id);
                itemName = menu ? menu.name : 'Menu #' + it.menu_id;
              }
              const opts = (it.order_item_options || []).map(o => o.option_value ? o.option_value.value : '').filter(Boolean);
              const optsStr = opts.length ? '<div class="text-muted" style="font-size:11px;">' + opts.join(', ') + '</div>' : '';
              return `<tr>
              <td>${itemName}${optsStr}</td>
              <td>${it.quantity}</td>
              <td>${fmtPrice(it.unit_price)}</td>
              <td>${fmtPrice(it.item_total)}</td>
            </tr>`;
            }).join('')}
          </tbody>
        </table>
      `);
    } catch (err) { App.toast(err.message, 'error'); }
  };

  // ===== NEW ORDER FORM =====
  function showNewOrderForm() {
    let items = [];
    let nextItemId = 0;

    App.modal('New Order', `
      <form id="order-form">
        <div class="form-row">
          <div class="form-group">
            <label>Customer (optional)</label>
            <select id="of-customer">
              <option value="">Walk-in</option>
              ${allCustomers.map(c => `<option value="${c.id}">${c.name}</option>`).join('')}
            </select>
          </div>
          <div class="form-group">
            <label>Order Type</label>
            <div class="radio-group">
              <label><input type="radio" name="order_type" value="counter" checked> Counter</label>
              <label><input type="radio" name="order_type" value="phone"> Phone</label>
            </div>
          </div>
        </div>
        <div class="form-group"><label>Notes</label><input id="of-notes" placeholder="Special instructions..."></div>

        <div class="section-title mt-16">Items</div>
        <div id="order-items"></div>
        <button type="button" class="btn btn-outline mb-16" id="add-item-btn">+ Add Item</button>

        <div class="price-preview" id="price-preview">Total: 0.00 €</div>
        <button type="submit" class="btn btn-block mt-12">Create Order</button>
      </form>
    `);

    document.getElementById('add-item-btn').addEventListener('click', addItem);
    document.getElementById('order-form').addEventListener('submit', submitOrder);

    addItem(); // start with one item

    function addItem() {
      const id = nextItemId++;
      items.push({ id, type: 'product', product_id: null, menu_id: null, quantity: 1, option_values: [] });

      const row = document.createElement('div');
      row.className = 'item-row';
      row.id = 'item-' + id;
      row.innerHTML = `
        <div class="form-group">
          <label>Type</label>
          <select data-field="type" onchange="orderItemTypeChange(${id}, this.value)">
            <option value="product">Product</option>
            <option value="menu">Menu</option>
          </select>
        </div>
        <div class="form-group grow">
          <label>Item</label>
          <select data-field="item_id" onchange="orderItemSelected(${id}, this.value)">
            <option value="">Select...</option>
            ${allProducts.map(p => `<option value="${p.id}">${p.name} (${fmtPrice(p.price)})</option>`).join('')}
          </select>
        </div>
        <div class="form-group">
          <label>Qty</label>
          <input type="number" min="1" value="1" style="width:60px" data-field="quantity" onchange="orderItemQtyChange(${id}, this.value)">
        </div>
        <div class="form-group" id="item-options-${id}" style="width:100%"></div>
        <button type="button" class="btn btn-sm btn-danger" onclick="removeOrderItem(${id})">X</button>
      `;
      document.getElementById('order-items').appendChild(row);
    }

    window.orderItemTypeChange = function(id, type) {
      const item = items.find(i => i.id === id);
      if (!item) return;
      item.type = type;
      item.product_id = null;
      item.menu_id = null;
      item.option_values = [];
      const sel = document.querySelector(`#item-${id} [data-field="item_id"]`);
      const source = type === 'product' ? allProducts : allMenus;
      sel.innerHTML = `<option value="">Select...</option>
        ${source.map(x => `<option value="${x.id}">${x.name} (${fmtPrice(x.price)})</option>`).join('')}`;
      document.getElementById('item-options-' + id).innerHTML = '';
      updatePrice();
    };

    window.orderItemSelected = async function(id, val) {
      const item = items.find(i => i.id === id);
      if (!item) return;
      if (item.type === 'product') {
        item.product_id = val ? Number(val) : null;
        item.menu_id = null;
        // Load options for this product
        if (val) await loadItemOptions(id, val);
        else document.getElementById('item-options-' + id).innerHTML = '';
      } else {
        item.menu_id = val ? Number(val) : null;
        item.product_id = null;
        document.getElementById('item-options-' + id).innerHTML = '';
      }
      updatePrice();
    };

    async function loadItemOptions(itemId, productId) {
      const el = document.getElementById('item-options-' + itemId);
      try {
        const opts = await App.api('/options/product/' + productId);
        const optList = Array.isArray(opts) ? opts : [];
        if (optList.length === 0) { el.innerHTML = ''; return; }

        let html = '<label class="text-muted" style="font-size:12px">Options:</label><div style="display:flex;flex-wrap:wrap;gap:6px;margin-top:4px">';
        for (const opt of optList) {
          let vals;
          try { vals = await App.api('/options/' + opt.id + '/values/'); } catch { vals = []; }
          const valList = Array.isArray(vals) ? vals : [];
          for (const v of valList) {
            html += `<label style="font-size:12px;cursor:pointer;display:flex;align-items:center;gap:4px;">
              <input type="checkbox" value="${v.id}" data-item-id="${itemId}" data-price="${v.option_price}" onchange="toggleItemOption(${itemId})">
              ${v.value} ${v.option_price > 0 ? '(+' + fmtPrice(v.option_price) + ')' : ''}
            </label>`;
          }
        }
        html += '</div>';
        el.innerHTML = html;
      } catch { el.innerHTML = ''; }
    }

    window.toggleItemOption = function(itemId) {
      const item = items.find(i => i.id === itemId);
      if (!item) return;
      const checkboxes = document.querySelectorAll(`#item-options-${itemId} input[type="checkbox"]:checked`);
      item.option_values = Array.from(checkboxes).map(cb => ({
        option_value_id: Number(cb.value),
        price: Number(cb.dataset.price)
      }));
      updatePrice();
    };

    window.orderItemQtyChange = function(id, qty) {
      const item = items.find(i => i.id === id);
      if (item) item.quantity = Number(qty) || 1;
      updatePrice();
    };

    window.removeOrderItem = function(id) {
      items = items.filter(i => i.id !== id);
      const row = document.getElementById('item-' + id);
      if (row) row.remove();
      updatePrice();
    };

    function updatePrice() {
      let total = 0;
      for (const item of items) {
        let unitPrice = 0;
        if (item.type === 'product' && item.product_id) {
          const p = allProducts.find(x => x.id === item.product_id);
          if (p) unitPrice = p.price;
        } else if (item.type === 'menu' && item.menu_id) {
          const m = allMenus.find(x => x.id === item.menu_id);
          if (m) unitPrice = m.price;
        }
        const optPrice = (item.option_values || []).reduce((s, v) => s + (v.price || 0), 0);
        total += (unitPrice + optPrice) * item.quantity;
      }
      const el = document.getElementById('price-preview');
      if (el) el.textContent = 'Total: ' + fmtPrice(total);
    }

    async function submitOrder(e) {
      e.preventDefault();
      const custVal = document.getElementById('of-customer').value;
      const orderType = document.querySelector('input[name="order_type"]:checked').value;
      const notes = document.getElementById('of-notes').value;

      const orderItems = items.filter(i => i.product_id || i.menu_id).map(i => {
        const obj = { quantity: i.quantity };
        if (i.type === 'product') obj.product_id = i.product_id;
        else obj.menu_id = i.menu_id;
        if (i.option_values && i.option_values.length > 0) {
          obj.order_item_options = i.option_values.map(v => ({ option_value_id: v.option_value_id }));
        }
        return obj;
      });

      if (orderItems.length === 0) return App.toast('Add at least one item', 'error');

      const body = {
        order_type: orderType,
        notes,
        order_items: orderItems,
      };
      if (custVal) body.customer_id = Number(custVal);

      try {
        await App.api('/orders/', { method: 'POST', body });
        App.closeModal();
        App.toast('Order created', 'success');
        const activeFilter = document.querySelector('.toolbar .tab-btn.active').dataset.filter;
        loadOrders(activeFilter);
      } catch (err) { App.toast(err.message, 'error'); }
    }
  }
});
