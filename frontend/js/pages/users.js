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
            <thead><tr><th>ID</th><th>Username</th><th>Email</th><th>Role</th><th>Created</th><th>Actions</th></tr></thead>
            <tbody>
              ${userList.map(u => `<tr>
                <td>${u.id}</td>
                <td>${u.username}</td>
                <td>${u.email}</td>
                <td>${u.Role ? u.Role.role_name : u.roles_id}</td>
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

  // ===== ROLES =====
  async function loadRoles() {
    const el = document.getElementById('tab-content');
    try {
      const roles = await App.api('/roles/');
      const list = Array.isArray(roles) ? roles : [];

      el.innerHTML = `
        <div class="toolbar">
          <button class="btn" id="new-role-btn">+ New Role</button>
        </div>
        <div class="table-wrap">
          <table>
            <thead><tr><th>ID</th><th>Name</th><th>Description</th><th>Permissions</th><th>Actions</th></tr></thead>
            <tbody>
              ${list.map(r => `<tr>
                <td>${r.id}</td>
                <td>${r.role_name}</td>
                <td class="text-muted">${r.description || '-'}</td>
                <td class="text-muted" style="max-width:200px;overflow:hidden;text-overflow:ellipsis">${r.permissions || '-'}</td>
                <td>
                  <button class="btn btn-sm btn-danger" onclick="deleteRole(${r.id})">Del</button>
                </td>
              </tr>`).join('')}
            </tbody>
          </table>
        </div>
      `;

      document.getElementById('new-role-btn').addEventListener('click', () => {
        App.modal('New Role', `
          <form id="role-form">
            <div class="form-group"><label>Role Name</label><input id="rf-name" required></div>
            <div class="form-group"><label>Description</label><input id="rf-desc"></div>
            <div class="form-group"><label>Permissions</label><textarea id="rf-perms" rows="3" placeholder="e.g. read,write,admin"></textarea></div>
            <button type="submit" class="btn btn-block">Create Role</button>
          </form>
        `);
        document.getElementById('role-form').addEventListener('submit', async e => {
          e.preventDefault();
          try {
            await App.api('/roles/', {
              method: 'POST',
              body: {
                role_name: document.getElementById('rf-name').value,
                description: document.getElementById('rf-desc').value,
                permissions: document.getElementById('rf-perms').value,
              }
            });
            App.closeModal();
            App.toast('Role created', 'success');
            loadRoles();
          } catch (err) { App.toast(err.message, 'error'); }
        });
      });
    } catch (err) { el.innerHTML = `<div class="empty-msg">${err.message}</div>`; }
  }

  window.deleteRole = async function(id) {
    if (!confirm('Delete this role?')) return;
    try {
      await App.api('/roles/' + id, { method: 'DELETE' });
      App.toast('Role deleted', 'success');
      loadRoles();
    } catch (err) { App.toast(err.message, 'error'); }
  };
});
