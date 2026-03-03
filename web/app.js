// === State ===
let authToken = '';
let currentHospital = '';

// === View Management ===
function showView(viewId) {
    document.querySelectorAll('.view').forEach(v => v.classList.remove('active'));
    document.getElementById(viewId).classList.add('active');
}

// === Tab Switching ===
function switchTab(tab) {
    document.querySelectorAll('.tab').forEach(t => t.classList.remove('active'));
    document.querySelector(`[data-tab="${tab}"]`).classList.add('active');

    document.querySelectorAll('.auth-form').forEach(f => f.classList.remove('active'));
    document.getElementById(`${tab}-form`).classList.add('active');

    hideMessage('auth-message');
}

// === Messages ===
function showMessage(elementId, text, type) {
    const el = document.getElementById(elementId);
    el.textContent = text;
    el.className = `message ${type}`;
}

function hideMessage(elementId) {
    const el = document.getElementById(elementId);
    el.className = 'message';
}

// === API Helper ===
async function apiCall(url, options = {}) {
    if (authToken) {
        options.headers = { ...options.headers, 'Authorization': `Bearer ${authToken}` };
    }
    const response = await fetch(url, options);
    const data = await response.json();
    return { ok: response.ok, status: response.status, data };
}

// === Auth Handlers ===
async function handleRegister(e) {
    e.preventDefault();
    const btn = document.getElementById('register-btn');
    btn.classList.add('loading');

    try {
        const body = {
            username: document.getElementById('reg-username').value,
            password: document.getElementById('reg-password').value,
            hospital: document.getElementById('reg-hospital').value,
        };

        const res = await apiCall('/staff/create', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(body),
        });

        if (res.ok) {
            showMessage('auth-message', '✓ Account created successfully! You can now login.', 'success');
            setTimeout(() => switchTab('login'), 1500);
        } else {
            showMessage('auth-message', res.data.error || 'Registration failed', 'error');
        }
    } catch (err) {
        showMessage('auth-message', 'Network error. Is the server running?', 'error');
    } finally {
        btn.classList.remove('loading');
    }
}

async function handleLogin(e) {
    e.preventDefault();
    const btn = document.getElementById('login-btn');
    btn.classList.add('loading');

    try {
        const body = {
            username: document.getElementById('login-username').value,
            password: document.getElementById('login-password').value,
            hospital: document.getElementById('login-hospital').value,
        };

        const res = await apiCall('/staff/login', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(body),
        });

        if (res.ok && res.data.token) {
            authToken = res.data.token;
            currentHospital = body.hospital;
            document.getElementById('nav-hospital').textContent = currentHospital;
            showView('dashboard-view');
        } else {
            showMessage('auth-message', res.data.error || 'Invalid credentials', 'error');
        }
    } catch (err) {
        showMessage('auth-message', 'Network error. Is the server running?', 'error');
    } finally {
        btn.classList.remove('loading');
    }
}

function handleLogout() {
    authToken = '';
    currentHospital = '';
    showView('auth-view');
    hideMessage('auth-message');
    // Clear forms
    document.getElementById('login-form').reset();
    document.getElementById('register-form').reset();
}

// === Patient Search ===
async function handleSearch(e) {
    e.preventDefault();
    hideMessage('search-message');

    const filters = {
        national_id: document.getElementById('s-national-id').value,
        passport_id: document.getElementById('s-passport-id').value,
        first_name: document.getElementById('s-first-name').value,
        middle_name: document.getElementById('s-middle-name').value,
        last_name: document.getElementById('s-last-name').value,
        date_of_birth: document.getElementById('s-dob').value,
        phone_number: document.getElementById('s-phone').value,
        email: document.getElementById('s-email').value,
    };

    // Build query string from non-empty filters
    const params = new URLSearchParams();
    Object.entries(filters).forEach(([key, value]) => {
        if (value.trim()) params.append(key, value.trim());
    });

    try {
        const res = await apiCall(`/patient/search?${params.toString()}`);

        if (res.ok) {
            renderResults(res.data.data || []);
        } else if (res.status === 401) {
            showMessage('search-message', 'Session expired. Please login again.', 'error');
            setTimeout(handleLogout, 2000);
        } else {
            showMessage('search-message', res.data.error || 'Search failed', 'error');
        }
    } catch (err) {
        showMessage('search-message', 'Network error. Is the server running?', 'error');
    }
}

function renderResults(patients) {
    const section = document.getElementById('results-section');
    const container = document.getElementById('results-container');
    const countEl = document.getElementById('results-count');

    section.style.display = 'block';
    countEl.textContent = `${patients.length} patient${patients.length !== 1 ? 's' : ''} found`;

    if (patients.length === 0) {
        container.innerHTML = `
            <div class="no-results glass">
                <div class="no-results-icon">📋</div>
                <p>No patients found matching your criteria.</p>
                <p style="font-size:13px;margin-top:8px;color:var(--text-muted)">
                    Try adjusting your search filters or use fewer criteria.
                </p>
            </div>`;
        return;
    }

    container.innerHTML = patients.map((p, i) => `
        <div class="patient-card" style="animation-delay:${i * 0.1}s">
            <div class="patient-card-header">
                <div>
                    <div class="patient-name">${esc(p.first_name_en || '')} ${esc(p.middle_name_en || '')} ${esc(p.last_name_en || '')}</div>
                    <div class="patient-name-th">${esc(p.first_name_th || '')} ${esc(p.middle_name_th || '')} ${esc(p.last_name_th || '')}</div>
                </div>
                <span class="badge">${esc(p.gender || 'N/A')}</span>
            </div>
            <div class="patient-details">
                ${detail('Patient HN', p.patient_hn)}
                ${detail('National ID', p.national_id)}
                ${detail('Passport ID', p.passport_id)}
                ${detail('Date of Birth', p.date_of_birth)}
                ${detail('Phone', p.phone_number)}
                ${detail('Email', p.email)}
                ${detail('Hospital', p.hospital_name)}
            </div>
        </div>
    `).join('');
}

function detail(label, value) {
    return `<div class="detail-item">
        <span class="detail-label">${label}</span>
        <span class="detail-value">${esc(value || '—')}</span>
    </div>`;
}

function esc(str) {
    const div = document.createElement('div');
    div.textContent = str;
    return div.innerHTML;
}

function clearSearch() {
    document.getElementById('search-form').reset();
    document.getElementById('results-section').style.display = 'none';
    hideMessage('search-message');
}
