App.registerPage('products', async () => {
  render(`
    <div class="tabs">
      <button class="tab-btn active" data-tab="categories">Categories</button>
      <button class="tab-btn" data-tab="products">Products</button>
      <button class="tab-btn" data-tab="options">Options</button>
    </div>
    <div id="tab-content"></div>
  `);

  const tabBtns = document.querySelectorAll('.tab-btn');
  tabBtns.forEach(btn => btn.addEventListener('click', () => {
    tabBtns.forEach(b => b.classList.remove('active'));
    btn.classList.add('active');
    loadTab(btn.dataset.tab);
  }));

  loadTab('categories');

  async function loadTab(tab) {
    if (tab === 'categories') await loadCategories();
    else if (tab === 'products') await loadProducts();
    else if (tab === 'options') await loadOptions();
  }

  // ===== CATEGORIES =====
  async function loadCategories() {
    const el = document.getElementById('tab-content');
    try {
      const cats = await App.api('/categories/');
      const list = Array.isArray(cats) ? cats : [];
      el.innerHTML = `
        <div class="toolbar">
          <button class="btn" onclick="showCategoryForm()">+ New Category</button>
        </div>
        <div class="table-wrap">
          <table>
            <thead><tr><th>ID</th><th>Name</th><th>Description</th><th>Order</th><th>Actions</th></tr></thead>
            <tbody>
              ${list.map(c => `<tr>
                <td>${c.id}</td>
                <td>${esc(c.name)}</td>
                <td class="text-muted">${esc(c.description) || '-'}</td>
                <td>${c.display_order || 0}</td>
                <td class="inline-flex">
                  <button class="btn btn-sm" onclick="showCategoryForm(${c.id})">Edit</button>
                  <button class="btn btn-sm btn-danger" onclick="deleteCategory(${c.id})">Del</button>
                </td>
              </tr>`).join('')}
            </tbody>
          </table>
        </div>
      `;
    } catch (err) { el.innerHTML = `<div class="empty-msg">${err.message}</div>`; }
  }

  window.showCategoryForm = async function(id) {
    let cat = { name: '', description: '', display_order: 0, image_url: '' };
    if (id) {
      try { cat = await App.api('/categories/' + id); } catch { return App.toast('Failed to load', 'error'); }
    }
    App.modal(id ? 'Edit Category' : 'New Category', `
      <form id="cat-form">
        <div class="form-group"><label>Name</label><input id="cf-name" value="${cat.name}" required></div>
        <div class="form-group"><label>Description</label><input id="cf-desc" value="${cat.description || ''}"></div>
        <div class="form-row">
          <div class="form-group"><label>Display Order</label><input type="number" id="cf-order" value="${cat.display_order || 0}"></div>
          <div class="form-group"><label>Image URL</label><input id="cf-img" value="${cat.image_url || ''}"></div>
        </div>
        <button type="submit" class="btn btn-block">${id ? 'Update' : 'Create'}</button>
      </form>
    `);
    document.getElementById('cat-form').addEventListener('submit', async e => {
      e.preventDefault();
      try {
        await App.api('/categories/' + (id || ''), {
          method: id ? 'PUT' : 'POST',
          body: {
            name: document.getElementById('cf-name').value,
            description: document.getElementById('cf-desc').value,
            display_order: Number(document.getElementById('cf-order').value),
            image_url: document.getElementById('cf-img').value,
          }
        });
        App.closeModal();
        App.toast(id ? 'Category updated' : 'Category created', 'success');
        loadCategories();
      } catch (err) { App.toast(err.message, 'error'); }
    });
  };

  window.deleteCategory = async function(id) {
    if (!confirm('Delete this category?')) return;
    try {
      await App.api('/categories/' + id, { method: 'DELETE' });
      App.toast('Deleted', 'success');
      loadCategories();
    } catch (err) { App.toast(err.message, 'error'); }
  };

  // ===== PRODUCTS =====
  async function loadProducts() {
    const el = document.getElementById('tab-content');
    try {
      const [prods, cats] = await Promise.all([App.api('/products/'), App.api('/categories/')]);
      const list = Array.isArray(prods) ? prods : [];
      const catList = Array.isArray(cats) ? cats : [];

      el.innerHTML = `
        <div class="toolbar">
          <button class="btn" onclick="showProductForm()">+ New Product</button>
          <select id="prod-cat-filter">
            <option value="">All Categories</option>
            ${catList.map(c => `<option value="${c.id}">${c.name}</option>`).join('')}
          </select>
        </div>
        <div class="table-wrap">
          <table>
            <thead><tr><th>ID</th><th>Name</th><th>Category</th><th>Price</th><th>Stock</th><th>Available</th><th>Actions</th></tr></thead>
            <tbody id="prod-tbody">
              ${renderProductRows(list)}
            </tbody>
          </table>
        </div>
      `;

      document.getElementById('prod-cat-filter').addEventListener('change', async e => {
        const catId = e.target.value;
        let filtered;
        if (catId) {
          try { filtered = await App.api('/products/category/' + catId); } catch { filtered = []; }
        } else {
          filtered = list;
        }
        document.getElementById('prod-tbody').innerHTML = renderProductRows(Array.isArray(filtered) ? filtered : []);
      });
    } catch (err) { el.innerHTML = `<div class="empty-msg">${err.message}</div>`; }
  }

  function renderProductRows(list) {
    return list.map(p => `<tr>
      <td>${p.id}</td>
      <td>${p.image_url ? `<img src="${esc(p.image_url)}" style="width:36px;height:36px;object-fit:cover;border-radius:4px;vertical-align:middle;margin-right:8px;">` : ''}${esc(p.name)}</td>
      <td>${p.category ? esc(p.category.name) : p.category_id}</td>
      <td>${fmtPrice(p.price)}</td>
      <td>
        <input type="number" value="${p.stock_quantity}" style="width:60px;padding:4px;background:var(--bg-input);border:1px solid var(--border);border-radius:4px;color:var(--text);"
          onchange="updateStock(${p.id}, this.value)">
      </td>
      <td>
        <label class="toggle">
          <input type="checkbox" ${p.is_available ? 'checked' : ''} onchange="toggleProductAvail(${p.id})">
          <span class="toggle-slider"></span>
        </label>
      </td>
      <td class="inline-flex">
        <button class="btn btn-sm" onclick="showProductForm(${p.id})">Edit</button>
        <button class="btn btn-sm btn-danger" onclick="deleteProduct(${p.id})">Del</button>
      </td>
    </tr>`).join('');
  }

  window.toggleProductAvail = async function(id) {
    try {
      await App.api('/products/' + id + '/availability', { method: 'PATCH' });
      App.toast('Availability toggled', 'success');
    } catch (err) { App.toast(err.message, 'error'); loadProducts(); }
  };

  window.updateStock = async function(id, qty) {
    try {
      await App.api('/products/' + id + '/stock', { method: 'PATCH', body: { stock_quantity: Number(qty) } });
      App.toast('Stock updated', 'success');
    } catch (err) { App.toast(err.message, 'error'); }
  };

  window.showProductForm = async function(id) {
    const cats = await App.api('/categories/');
    const catList = Array.isArray(cats) ? cats : [];
    let prod = { name: '', description: '', price: 0, category_id: '', stock_quantity: 0, preparation_time: 0, image_url: '', is_available: true };
    if (id) {
      try { prod = await App.api('/products/' + id); } catch { return App.toast('Failed to load', 'error'); }
    }
    App.modal(id ? 'Edit Product' : 'New Product', `
      <form id="prod-form">
        <div class="form-group"><label>Name</label><input id="pf-name" value="${prod.name}" required></div>
        <div class="form-group"><label>Description</label><textarea id="pf-desc" rows="2">${prod.description || ''}</textarea></div>
        <div class="form-row">
          <div class="form-group"><label>Category</label>
            <select id="pf-cat" required>
              <option value="">Select...</option>
              ${catList.map(c => `<option value="${c.id}" ${c.id === prod.category_id ? 'selected' : ''}>${c.name}</option>`).join('')}
            </select>
          </div>
          <div class="form-group"><label>Price</label><input type="number" step="0.01" id="pf-price" value="${prod.price}" required></div>
        </div>
        <div class="form-row">
          <div class="form-group"><label>Stock</label><input type="number" id="pf-stock" value="${prod.stock_quantity}"></div>
          <div class="form-group"><label>Prep Time (min)</label><input type="number" id="pf-prep" value="${prod.preparation_time || 0}"></div>
        </div>
        <div class="form-group"><label>Image URL</label><input id="pf-img" value="${prod.image_url || ''}" placeholder="https://..."></div>
        <button type="submit" class="btn btn-block">${id ? 'Update' : 'Create'}</button>
      </form>
    `);
    document.getElementById('prod-form').addEventListener('submit', async e => {
      e.preventDefault();
      try {
        await App.api('/products/' + (id || ''), {
          method: id ? 'PUT' : 'POST',
          body: {
            name: document.getElementById('pf-name').value,
            description: document.getElementById('pf-desc').value,
            category_id: Number(document.getElementById('pf-cat').value),
            price: Number(document.getElementById('pf-price').value),
            stock_quantity: Number(document.getElementById('pf-stock').value),
            preparation_time: Number(document.getElementById('pf-prep').value),
            image_url: document.getElementById('pf-img').value,
            is_available: true,
          }
        });
        App.closeModal();
        App.toast(id ? 'Product updated' : 'Product created', 'success');
        loadProducts();
      } catch (err) { App.toast(err.message, 'error'); }
    });
  };

  window.deleteProduct = async function(id) {
    if (!confirm('Delete this product?')) return;
    try {
      await App.api('/products/' + id, { method: 'DELETE' });
      App.toast('Deleted', 'success');
      loadProducts();
    } catch (err) { App.toast(err.message, 'error'); }
  };

  // ===== OPTIONS =====
  async function loadOptions() {
    const el = document.getElementById('tab-content');
    try {
      const prods = await App.api('/products/');
      const list = Array.isArray(prods) ? prods : [];

      el.innerHTML = `
        <div class="toolbar">
          <select id="opt-prod-sel">
            <option value="">Select a product</option>
            ${list.map(p => `<option value="${p.id}">${p.name}</option>`).join('')}
          </select>
          <button class="btn" id="add-option-btn" disabled>+ New Option</button>
        </div>
        <div id="options-content"><p class="text-muted">Select a product to view options</p></div>
      `;

      const sel = document.getElementById('opt-prod-sel');
      const addBtn = document.getElementById('add-option-btn');
      sel.addEventListener('change', () => {
        addBtn.disabled = !sel.value;
        if (sel.value) loadProductOptions(sel.value);
      });
      addBtn.addEventListener('click', () => { if (sel.value) showOptionForm(sel.value); });
    } catch (err) { el.innerHTML = `<div class="empty-msg">${err.message}</div>`; }
  }

  async function loadProductOptions(productId) {
    const el = document.getElementById('options-content');
    try {
      const opts = await App.api('/options/product/' + productId);
      const list = Array.isArray(opts) ? opts : [];
      if (list.length === 0) {
        el.innerHTML = '<p class="text-muted mt-12">No options for this product</p>';
        return;
      }
      el.innerHTML = list.map(o => `
        <div class="card mb-8">
          <div class="flex-between mb-8">
            <div>
              <strong>${esc(o.name)}</strong>
              <span class="text-muted" style="margin-left:8px">${o.is_unique} | ${o.is_required ? 'Required' : 'Optional'}</span>
            </div>
            <div class="inline-flex">
              <button class="btn btn-sm" onclick="showOptionForm(${productId}, ${o.id})">Edit</button>
              <button class="btn btn-sm btn-danger" onclick="deleteOption(${o.id}, ${productId})">Del</button>
              <button class="btn btn-sm btn-info" onclick="showValueForm(${o.id}, ${productId})">+ Value</button>
            </div>
          </div>
          <div id="vals-${o.id}">Loading values...</div>
        </div>
      `).join('');

      // Load values for each option
      for (const o of list) loadOptionValues(o.id, productId);
    } catch (err) { el.innerHTML = `<div class="empty-msg">${err.message}</div>`; }
  }

  async function loadOptionValues(optionId, productId) {
    const el = document.getElementById('vals-' + optionId);
    if (!el) return;
    try {
      const vals = await App.api('/options/' + optionId + '/values/');
      const list = Array.isArray(vals) ? vals : [];
      if (list.length === 0) { el.innerHTML = '<p class="text-muted">No values</p>'; return; }
      el.innerHTML = `<table class="sub-table"><thead><tr><th>Value</th><th>Price</th><th>Actions</th></tr></thead><tbody>
        ${list.map(v => `<tr>
          <td>${esc(v.value)}</td>
          <td>${fmtPrice(v.option_price)}</td>
          <td class="inline-flex">
            <button class="btn btn-sm" onclick="showValueForm(${optionId}, ${productId}, ${v.id})">Edit</button>
            <button class="btn btn-sm btn-danger" onclick="deleteValue(${v.id}, ${optionId}, ${productId})">Del</button>
          </td>
        </tr>`).join('')}
      </tbody></table>`;
    } catch { el.innerHTML = '<p class="text-muted">Failed to load values</p>'; }
  }

  window.showOptionForm = async function(productId, optId) {
    let opt = { name: '', is_unique: 'single', is_required: false };
    if (optId) {
      try { opt = await App.api('/options/' + optId); } catch { return App.toast('Failed to load', 'error'); }
    }
    App.modal(optId ? 'Edit Option' : 'New Option', `
      <form id="opt-form">
        <div class="form-group"><label>Name</label><input id="of-name" value="${opt.name}" required></div>
        <div class="form-row">
          <div class="form-group"><label>Type</label>
            <select id="of-unique">
              <option value="single" ${opt.is_unique === 'single' ? 'selected' : ''}>Single choice</option>
              <option value="multiple" ${opt.is_unique === 'multiple' ? 'selected' : ''}>Multiple choice</option>
            </select>
          </div>
          <div class="form-group"><label>Required</label>
            <select id="of-req">
              <option value="false" ${!opt.is_required ? 'selected' : ''}>No</option>
              <option value="true" ${opt.is_required ? 'selected' : ''}>Yes</option>
            </select>
          </div>
        </div>
        <button type="submit" class="btn btn-block">${optId ? 'Update' : 'Create'}</button>
      </form>
    `);
    document.getElementById('opt-form').addEventListener('submit', async e => {
      e.preventDefault();
      try {
        await App.api('/options/' + (optId || ''), {
          method: optId ? 'PUT' : 'POST',
          body: {
            product_id: Number(productId),
            name: document.getElementById('of-name').value,
            is_unique: document.getElementById('of-unique').value,
            is_required: document.getElementById('of-req').value === 'true',
          }
        });
        App.closeModal();
        App.toast(optId ? 'Option updated' : 'Option created', 'success');
        loadProductOptions(productId);
      } catch (err) { App.toast(err.message, 'error'); }
    });
  };

  window.deleteOption = async function(id, productId) {
    if (!confirm('Delete this option?')) return;
    try {
      await App.api('/options/' + id, { method: 'DELETE' });
      App.toast('Deleted', 'success');
      loadProductOptions(productId);
    } catch (err) { App.toast(err.message, 'error'); }
  };

  window.showValueForm = async function(optionId, productId, valId) {
    let val = { value: '', option_price: 0 };
    if (valId) {
      try { val = await App.api('/options/values/' + valId); } catch { return App.toast('Failed to load', 'error'); }
    }
    App.modal(valId ? 'Edit Value' : 'New Value', `
      <form id="val-form">
        <div class="form-group"><label>Value</label><input id="vf-val" value="${val.value}" required></div>
        <div class="form-group"><label>Extra Price</label><input type="number" step="0.01" id="vf-price" value="${val.option_price}"></div>
        <button type="submit" class="btn btn-block">${valId ? 'Update' : 'Create'}</button>
      </form>
    `);
    document.getElementById('val-form').addEventListener('submit', async e => {
      e.preventDefault();
      const path = valId ? '/options/values/' + valId : '/options/' + optionId + '/values/';
      const valueObj = {
        option_id: Number(optionId),
        value: document.getElementById('vf-val').value,
        option_price: Number(document.getElementById('vf-price').value),
      };
      try {
        await App.api(path, {
          method: valId ? 'PUT' : 'POST',
          // Backend expects an array for creation, single object for update
          body: valId ? valueObj : [valueObj],
        });
        App.closeModal();
        App.toast(valId ? 'Value updated' : 'Value created', 'success');
        loadProductOptions(productId);
      } catch (err) { App.toast(err.message, 'error'); }
    });
  };

  window.deleteValue = async function(id, optionId, productId) {
    if (!confirm('Delete this value?')) return;
    try {
      await App.api('/options/values/' + id, { method: 'DELETE' });
      App.toast('Deleted', 'success');
      loadProductOptions(productId);
    } catch (err) { App.toast(err.message, 'error'); }
  };
});
