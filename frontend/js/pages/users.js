App.registerPage('users', async () => {
  render(`
    <div class="tabs">
      <button class="tab-btn active" data-tab="users">Users</button>
      <button class="tab-btn" data-tab="roles">Roles</button>
    </div>
    <div id="tab-content"></div>
  `);

  const tabBtns = document.querySelectorAll('.tab-btn');
  tabBtns.forEach(btn => btn.addEventListener('click', () => {
    tabBtns.forEach(b => b.classList.remove('active'));
    btn.classList.add('active');
    loadTab(btn.dataset.tab);
  }));

  loadTab('users');

  function loadTab(tab) {
    if (tab === 'users') loadUsers();
    else loadRoles();
  }

  // ===== USERS =====
  async function loadUsers() {
    const el = document.getElementById('tab-content');
    try {
      const [users, roles] = await Promise.all([
        App.api('/users/'),
        App.api('/roles/'),
      ]);
      const userList = Array.isArray(users) ? users : [];
      const roleList = Array.isArray(roles) ? roles : [];

      el.innerHTML = `
        <div class="toolbar">
          <button class="btn" id="new-user-btn">+ New User</button>
        </div>
        <div class="table-wrap">
          <table>
            <thead><tr><th>ID</th><th>Username</th><th>Email</th><th>Role</th><th>Active</th><th>Created</th><th>Actions</th></tr></thead>
            <tbody>
              ${userList.map(u => `<tr>
                <td>${u.id}</td>
                <td>${esc(u.username)}</td>
                <td>${esc(u.email)}</td>
                <td>${u.Role ? esc(u.Role.role_name) : u.roles_id}</td>
                <td>
                  <label class="toggle">
                    <input type="checkbox" ${u.is_active ? 'checked' : ''} onchange="toggleUserStatus(${u.id})">
                    <span class="toggle-slider"></span>
                  </label>
                </td>
                <td>${fmtDate(u.created_at)}</td>
                <td>
                  <button class="btn btn-sm btn-danger" onclick="deleteUser(${u.id})">Del</button>
                </td>
              </tr>`).join('')}
            </tbody>
          </table>
        </div>
      `;

      document.getElementById('new-user-btn').addEventListener('click', () => {
        App.modal('New User', `
          <form id="user-form">
            <div class="form-group"><label>Username</label><input id="uf-name" required></div>
            <div class="form-group"><label>Email</label><input type="email" id="uf-email" required></div>
            <div class="form-group"><label>Password</label><input type="password" id="uf-pass" required minlength="6"></div>
            <div class="form-group"><label>Role</label>
              <select id="uf-role" required>
                <option value="">Select role...</option>
                ${roleList.map(r => `<option value="${r.id}">${r.role_name}</option>`).join('')}
              </select>
            </div>
            <button type="submit" class="btn btn-block">Create User</button>
          </form>
        `);
        document.getElementById('user-form').addEventListener('submit', async e => {
          e.preventDefault();
          try {
            await App.api('/users/', {
              method: 'POST',
              body: {
                username: document.getElementById('uf-name').value,
                email: document.getElementById('uf-email').value,
                password: document.getElementById('uf-pass').value,
                roles_id: Number(document.getElementById('uf-role').value),
              }
            });
            App.closeModal();
            App.toast('User created', 'success');
            loadUsers();
          } catch (err) { App.toast(err.message, 'error'); }
        });
      });
    } catch (err) { el.innerHTML = `<div class="empty-msg">${err.message}</div>`; }
  }

  window.deleteUser = async function(id) {
    if (!confirm('Delete this user?')) return;
    try {
      await App.api('/users/' + id, { method: 'DELETE' });
      App.toast('User deleted', 'success');
      loadUsers();
    } catch (err) { App.toast(err.message, 'error'); }
  };

  window.toggleUserStatus = async function(id) {
    try {
      await App.api('/users/' + id + '/status', { method: 'PATCH' });
      App.toast('User status updated', 'success');
    } catch (err) { App.toast(err.message, 'error'); loadUsers(); }
  };

  // ===== ROLES & PERMISSIONS =====
  function loadRoles() {
    const el = document.getElementById('tab-content');

    const permissions = [
      { resource: 'Users',      actions: 'Create, View, Delete, Activate/Deactivate', admin: true, accueil: false, preparation: false },
      { resource: 'Roles',      actions: 'View',                                      admin: true, accueil: false, preparation: false },
      { resource: 'Products',   actions: 'View',                                      admin: true, accueil: true,  preparation: true },
      { resource: 'Products',   actions: 'Create, Edit, Delete, Stock, Availability', admin: true, accueil: false, preparation: false },
      { resource: 'Categories', actions: 'View',                                      admin: true, accueil: true,  preparation: true },
      { resource: 'Categories', actions: 'Create, Edit, Delete',                      admin: true, accueil: false, preparation: false },
      { resource: 'Menus',      actions: 'View',                                      admin: true, accueil: true,  preparation: true },
      { resource: 'Menus',      actions: 'Create, Edit, Delete, Availability',        admin: true, accueil: false, preparation: false },
      { resource: 'Customers',  actions: 'Create, View, Edit, Delete',                admin: true, accueil: true,  preparation: false },
      { resource: 'Orders',     actions: 'View',                                      admin: true, accueil: true,  preparation: true },
      { resource: 'Orders',     actions: 'Create, Cancel',                            admin: true, accueil: true,  preparation: false },
      { resource: 'Orders',     actions: 'Update Status',                             admin: true, accueil: true,  preparation: true },
    ];

    const check = v => v ? '<span class="text-accent">&#10003;</span>' : '<span class="text-muted">&#10007;</span>';

    el.innerHTML = `
      <div class="section-title mt-16">Roles & Permissions</div>
      <p class="text-muted mb-16">Roles are predefined. The table below shows what each role can access.</p>
      <div class="table-wrap">
        <table>
          <thead><tr><th>Resource</th><th>Actions</th><th>Admin</th><th>Accueil</th><th>Preparation</th></tr></thead>
          <tbody>
            ${permissions.map(p => `<tr>
              <td><strong>${p.resource}</strong></td>
              <td class="text-muted">${p.actions}</td>
              <td style="text-align:center">${check(p.admin)}</td>
              <td style="text-align:center">${check(p.accueil)}</td>
              <td style="text-align:center">${check(p.preparation)}</td>
            </tr>`).join('')}
          </tbody>
        </table>
      </div>
    `;
  }
});
