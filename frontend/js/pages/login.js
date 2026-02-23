App.registerPage('login', async () => {
  document.getElementById('content').innerHTML = '';

  const wrap = document.createElement('div');
  wrap.className = 'login-wrap';
  wrap.innerHTML = `
    <div class="login-box">
      <h2>WacDo Admin</h2>
      <p class="sub">Sign in to your account</p>
      <form id="login-form">
        <div class="form-group">
          <label>Email</label>
          <input type="email" id="login-email" required placeholder="admin@wacdo.com">
        </div>
        <div class="form-group">
          <label>Password</label>
          <input type="password" id="login-pass" required placeholder="Enter password">
        </div>
        <button type="submit" class="btn btn-block" id="login-submit">Login</button>
      </form>
    </div>
  `;

  // Replace content area with full-screen login
  document.getElementById('content').appendChild(wrap);

  document.getElementById('login-form').addEventListener('submit', async e => {
    e.preventDefault();
    const btn = document.getElementById('login-submit');
    btn.disabled = true;
    btn.textContent = 'Signing in...';

    try {
      const token = await App.api('/users/login', {
        method: 'POST',
        body: {
          email: document.getElementById('login-email').value,
          password: document.getElementById('login-pass').value,
        }
      });
      // Backend returns raw JSON string: c.JSON(200, tokenString)
      // await res.json() gives the string directly
      App.setToken(token);
      App.toast('Login successful', 'success');
      App.navigate('dashboard');
    } catch (err) {
      App.toast(err.message, 'error');
      btn.disabled = false;
      btn.textContent = 'Login';
    }
  });
});
