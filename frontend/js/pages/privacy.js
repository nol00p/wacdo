App.registerPage('privacy', async () => {
  render(`
    <div style="max-width:720px;">
      <h2>Privacy Policy & GDPR Notice</h2>

      <div class="card mb-16">
        <h3>Data We Collect</h3>
        <p>WacDo collects and stores the following personal data to operate the ordering system:</p>
        <ul>
          <li><strong>Staff users:</strong> username, email address, encrypted password, assigned role.</li>
          <li><strong>Customers:</strong> name, phone number, email address.</li>
          <li><strong>Orders:</strong> order details, items, notes, timestamps, and the staff member who created the order.</li>
        </ul>
      </div>

      <div class="card mb-16">
        <h3>How We Use Your Data</h3>
        <ul>
          <li>Customer information is used solely to manage and fulfill orders.</li>
          <li>Staff credentials are used for authentication and access control.</li>
          <li>Order history is kept for operational tracking and reporting.</li>
        </ul>
      </div>

      <div class="card mb-16">
        <h3>Data Sharing</h3>
        <p>Personal data is <strong>not shared with third parties</strong>. All data remains within the WacDo system and is only accessible to authorized staff members based on their role.</p>
      </div>

      <div class="card mb-16">
        <h3>Your Rights (GDPR)</h3>
        <p>In accordance with the General Data Protection Regulation, individuals have the right to:</p>
        <ul>
          <li><strong>Right of access</strong> &mdash; Request a copy of your personal data held in the system.</li>
          <li><strong>Right to rectification</strong> &mdash; Request correction of inaccurate data.</li>
          <li><strong>Right to erasure</strong> &mdash; Request deletion of your personal data.</li>
        </ul>
        <p>Staff members with appropriate permissions can view, edit, and delete customer records directly from the <a href="#customers">Customers</a> page to fulfill these requests.</p>
      </div>

      <div class="card mb-16">
        <h3>Data Retention</h3>
        <p>Customer data and order history are retained as long as they are needed for operational purposes. Customers may request deletion at any time by contacting staff.</p>
      </div>

      <div class="card mb-16">
        <h3>Data Security</h3>
        <ul>
          <li>Passwords are hashed using bcrypt and never stored in plain text.</li>
          <li>API access is protected by JWT authentication.</li>
          <li>Role-based access control restricts data access to authorized personnel.</li>
          <li>Security headers (CSP, XSS filter, frame protection) are enforced.</li>
        </ul>
      </div>
    </div>
  `);
});
