package web

var emailSimHTML = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Agent House — Email Simulator</title>
  <link href="https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700&family=JetBrains+Mono:wght@400;500&display=swap" rel="stylesheet">
  <style>
:root{--bg:#080B14;--base:#0D1117;--raised:#161B27;--overlay:#1E2535;--border:#2A3345;--text:#E6EDF3;--muted:#8B949E;--dim:#484F58;--accent:#3B82F6;--green:#22C55E;--red:#EF4444;--yellow:#F6A623;--font:'Inter',system-ui,sans-serif;--mono:'JetBrains Mono',monospace}
*{box-sizing:border-box;margin:0;padding:0}
body{font-family:var(--font);background:var(--bg);color:var(--text);min-height:100vh;-webkit-font-smoothing:antialiased}
.header{display:flex;align-items:center;gap:16px;padding:0 24px;height:52px;background:var(--base);border-bottom:1px solid var(--border)}
.logo{font-size:14px;font-weight:700}.logo span{color:var(--accent)}
.nav{display:flex;gap:8px;margin-left:auto}
.nav a{color:var(--muted);font-size:12px;padding:4px 10px;border-radius:6px;border:1px solid var(--border);text-decoration:none}
.nav a:hover{color:var(--text);border-color:var(--accent)}
.container{max-width:900px;margin:0 auto;padding:24px}
h2{font-size:16px;margin-bottom:16px}
.scenarios{display:flex;flex-wrap:wrap;gap:8px;margin-bottom:24px}
.scenario-btn{padding:8px 16px;border-radius:8px;border:1px solid var(--border);background:var(--raised);color:var(--text);font-size:12px;cursor:pointer;font-family:var(--font);transition:all 0.15s}
.scenario-btn:hover{border-color:var(--accent);background:var(--overlay)}
.scenario-btn.danger{border-color:rgba(239,68,68,0.3)}
.scenario-btn.danger:hover{border-color:var(--red);background:rgba(239,68,68,0.1)}
.form{background:var(--raised);border:1px solid var(--border);border-radius:10px;padding:20px;margin-bottom:24px}
.form-row{display:flex;gap:12px;margin-bottom:12px;align-items:center}
.form-row label{width:70px;font-size:12px;color:var(--muted);flex-shrink:0}
.form-row input{flex:1;background:var(--overlay);border:1px solid var(--border);color:var(--text);padding:8px 12px;border-radius:6px;font-size:13px;font-family:var(--font)}
.form-row input:focus{outline:none;border-color:var(--accent)}
.form-row textarea{flex:1;background:var(--overlay);border:1px solid var(--border);color:var(--text);padding:10px 12px;border-radius:6px;font-size:13px;font-family:var(--font);min-height:120px;resize:vertical}
.form-row textarea:focus{outline:none;border-color:var(--accent)}
.send-btn{background:var(--accent);color:#fff;border:none;padding:10px 24px;border-radius:8px;font-size:13px;font-weight:600;cursor:pointer;width:100%}
.send-btn:hover{filter:brightness(1.1)}
.send-btn:disabled{opacity:0.5;cursor:not-allowed}
.log{background:var(--raised);border:1px solid var(--border);border-radius:10px;padding:16px}
.log-title{font-size:12px;font-weight:600;color:var(--muted);text-transform:uppercase;letter-spacing:0.05em;margin-bottom:12px}
.log-entry{display:flex;gap:10px;padding:8px 0;border-bottom:1px solid var(--border);font-size:12px;align-items:flex-start}
.log-entry:last-child{border-bottom:none}
.log-icon{font-size:14px;flex-shrink:0;margin-top:1px}
.log-body{flex:1}
.log-time{color:var(--dim);font-family:var(--mono);font-size:10px}
.log-msg{color:var(--muted);margin-top:2px}
.alert-box{background:rgba(239,68,68,0.08);border:1px solid rgba(239,68,68,0.3);border-radius:8px;padding:12px;margin-top:8px;font-size:11px}
.alert-box.warning{background:rgba(246,166,35,0.08);border-color:rgba(246,166,35,0.3)}
.alert-box .title{font-weight:600;color:var(--red);margin-bottom:4px}
.alert-box.warning .title{color:var(--yellow)}
.alert-box .factors{color:var(--muted);margin-top:4px}
.trust-badge{display:inline-block;padding:2px 6px;border-radius:4px;font-size:10px;font-weight:600}
.trust-badge.trusted{background:rgba(34,197,94,0.15);color:var(--green)}
.trust-badge.new_contact{background:rgba(246,166,35,0.15);color:var(--yellow)}
.trust-badge.impersonation{background:rgba(239,68,68,0.15);color:var(--red)}
.empty{color:var(--dim);text-align:center;padding:20px;font-size:12px}
  </style>
</head>
<body>
<div class="header">
  <div class="logo">Agent<span>House</span> <span style="color:var(--muted);font-weight:400">Email Simulator</span></div>
  <div class="nav"><a href="/mission">Mission</a><a href="/live">Live</a><a href="/">Dashboard</a></div>
</div>
<div class="container">
  <h2>Simulate Incoming Emails</h2>
  <div class="scenarios">
    <button class="scenario-btn" onclick="loadScenario('vendor_delay')">📦 Vendor Delay</button>
    <button class="scenario-btn" onclick="loadScenario('invoice')">💰 Invoice</button>
    <button class="scenario-btn" onclick="loadScenario('job_application')">📋 Job Application</button>
    <button class="scenario-btn" onclick="loadScenario('client_query')">👤 Client Query</button>
    <button class="scenario-btn danger" onclick="loadScenario('impersonation')">🚨 Impersonation</button>
    <button class="scenario-btn" onclick="loadScenario('internal')">📝 Internal Request</button>
  </div>

  <div class="form">
    <div class="form-row"><label>From:</label><input id="from" placeholder="sender@example.com"></div>
    <div class="form-row"><label>Name:</label><input id="fromName" placeholder="Sender Name"></div>
    <div class="form-row"><label>To:</label><input id="to" value="projects@epc-demo.com"></div>
    <div class="form-row"><label>Subject:</label><input id="subject" placeholder="Email subject"></div>
    <div class="form-row"><label>Body:</label><textarea id="body" placeholder="Email body..."></textarea></div>
    <button class="send-btn" id="sendBtn" onclick="sendEmail()">Send to Agent House →</button>
  </div>

  <div class="log">
    <div class="log-title">Sent Log</div>
    <div id="logEntries"><div class="empty">No emails sent yet. Click a scenario above to get started.</div></div>
  </div>
</div>

<script>
const SCENARIOS = {
  vendor_delay: {
    from: 'ahmed.rahman@vendorxyz.com',
    fromName: 'Ahmed Rahman',
    subject: 'Steel Delivery Update — Project Alpha',
    body: 'Dear Sir,\n\nWe regret to inform you that the structural steel delivery (Order #ST-2026-0847) for Project Alpha will be delayed by approximately 2 weeks.\n\nOriginal delivery date: April 5, 2026\nRevised delivery date: April 19, 2026\n\nThe delay is due to supply chain disruptions at our manufacturing facility in Jamshedpur. We are working to expedite the order and will provide weekly updates.\n\nWe apologize for the inconvenience.\n\nBest regards,\nAhmed Rahman\nSales Manager, XYZ Steel Corp'
  },
  invoice: {
    from: 'accounts@pqrconcrete.com',
    fromName: 'PQR Accounts',
    subject: 'Invoice #1042 — Concrete Supply for Project Alpha',
    body: 'Dear Accounts Payable,\n\nPlease find attached Invoice #1042 for concrete supply delivered to Project Alpha site.\n\nInvoice Details:\n- Order: PO-2026-0312\n- Quantity: 450 cubic meters M30 grade concrete\n- Unit Price: $100.44/m³\n- Total Amount: $45,200.00\n- Payment Terms: Net 30\n\nDelivery was completed on March 15, 2026.\n\nPlease process payment at your earliest convenience.\n\nRegards,\nPQR Concrete Ltd\nAccounts Department'
  },
  job_application: {
    from: 'raj.kumar@email.com',
    fromName: 'Raj Kumar',
    subject: 'Application — Site Supervisor for Highway Bridge Project',
    body: 'Dear Hiring Manager,\n\nI am writing to apply for the Site Supervisor position for the Highway Bridge Project (Gamma).\n\nProfessional Summary:\n- 12 years of experience in civil engineering and construction\n- 4 bridge projects completed (2 highway, 2 rail)\n- Managed teams of 50+ workers\n\nCertifications:\n- Professional Engineer (PE)\n- PMP Certified\n- OSHA-30 Construction Safety\n\nKey Skills:\n- Bridge construction and pile driving\n- Team management and subcontractor coordination\n- Quality control and safety compliance\n- Primavera P6 and AutoCAD proficient\n\nI am available for an interview at your convenience.\n\nBest regards,\nRaj Kumar\nPhone: +91-98765-43210'
  },
  client_query: {
    from: 'sarah.jones@clientabc.com',
    fromName: 'Sarah Jones',
    subject: 'Request for Phase 2 Progress Update — Project Alpha',
    body: 'Dear Project Team,\n\nAs we approach the end of Q1, I would like to request a progress update on Phase 2 of Project Alpha (Solar Farm).\n\nSpecifically, I need:\n1. Current completion percentage vs. planned\n2. Budget utilization report\n3. Updated timeline with any changes\n4. Risk register with mitigation status\n\nCould you please prepare this report by end of week? Our board meeting is scheduled for April 1st and I need to present the project status.\n\nThank you,\nSarah Jones\nVP Operations, ABC Energy Corp'
  },
  impersonation: {
    from: 'ahmed.r@gmail.com',
    fromName: 'Ahmed Rahman - XYZ Steel',
    subject: 'URGENT — Updated Bank Account Details for Payment',
    body: 'Dear Accounts Team,\n\nThis is Ahmed Rahman from XYZ Steel Corp. Due to a change in our banking arrangements, please update our payment details immediately.\n\nNew Bank Details:\n- Bank: First National Bank\n- Account: 8834-2291-0057\n- Routing: 021000089\n- SWIFT: FNBKUS33\n\nPlease use these details for all future payments including the pending Invoice #ST-0847 ($340,000).\n\nThis change is effective immediately. Please confirm receipt.\n\nUrgent regards,\nAhmed Rahman\nXYZ Steel Corp'
  },
  internal: {
    from: 'site.engineer@company.com',
    fromName: 'David Chen',
    subject: 'Request for Additional Resources — Project Alpha Site',
    body: 'Hi Team,\n\nWe need additional resources at the Project Alpha site for the foundation work:\n\n1. 2 additional concrete workers (skilled)\n2. 1 crane operator (certified for 50T mobile crane)\n3. Safety harnesses (20 units) — current stock damaged\n\nTimeline: Needed by next Monday (March 25)\n\nReason: Foundation pour scheduled for March 28-30. Current team is undersized for the scope.\n\nPlease approve and arrange.\n\nThanks,\nDavid Chen\nSite Engineer, Project Alpha'
  }
};

function loadScenario(name) {
  const s = SCENARIOS[name];
  if (!s) return;
  document.getElementById('from').value = s.from;
  document.getElementById('fromName').value = s.fromName;
  document.getElementById('subject').value = s.subject;
  document.getElementById('body').value = s.body;
  // Highlight active button
  document.querySelectorAll('.scenario-btn').forEach(b => b.style.borderColor = '');
  event.target.style.borderColor = name === 'impersonation' ? '#EF4444' : '#3B82F6';
}

async function sendEmail() {
  const btn = document.getElementById('sendBtn');
  btn.disabled = true;
  btn.textContent = 'Sending...';

  const payload = {
    from: document.getElementById('from').value,
    from_name: document.getElementById('fromName').value,
    to: document.getElementById('to').value,
    subject: document.getElementById('subject').value,
    body: document.getElementById('body').value
  };

  try {
    const res = await fetch('/api/email/receive', {
      method: 'POST',
      headers: {'Content-Type': 'application/json'},
      body: JSON.stringify(payload)
    });
    const data = await res.json();
    addLogEntry(data);
  } catch(e) {
    addLogEntry({success: false, error: e.message});
  }

  btn.disabled = false;
  btn.textContent = 'Send to Agent House →';
}

function addLogEntry(data) {
  const el = document.getElementById('logEntries');
  if (el.querySelector('.empty')) el.innerHTML = '';

  const entry = document.createElement('div');
  entry.className = 'log-entry';

  const time = new Date().toLocaleTimeString('en-US', {hour12:false});
  const email = data.email || {};
  const trust = data.trust || {};
  const alert = data.alert;

  let icon = '✅';
  if (trust.status === 'impersonation') icon = '🚨';
  else if (trust.status === 'new_contact') icon = '⚠️';
  else if (email.category === 'job_application') icon = '📋';

  let trustBadge = '';
  if (trust.status) {
    trustBadge = '<span class="trust-badge ' + trust.status + '">' + trust.status.replace('_', ' ') + '</span>';
  }

  let html = '<div class="log-icon">' + icon + '</div><div class="log-body">' +
    '<div><strong>' + esc(email.subject || 'Email') + '</strong> ' + trustBadge + ' <span class="log-time">' + time + '</span></div>' +
    '<div class="log-msg">From: ' + esc(email.from || '') + ' | Category: ' + (email.category || 'unknown') + '</div>';

  if (alert) {
    const alertClass = alert.level === 'critical' ? '' : ' warning';
    html += '<div class="alert-box' + alertClass + '"><div class="title">' + esc(alert.message) + '</div>';
    if (alert.risk_factors) {
      html += '<div class="factors">' + alert.risk_factors.map(f => '• ' + esc(f)).join('<br>') + '</div>';
    }
    html += '<div style="margin-top:6px;color:var(--text)">' + esc(alert.action) + '</div></div>';
  }

  if (email.category === 'job_application') {
    html += '<div class="log-msg" style="color:var(--green)">✓ Applicant stored for future matching</div>';
  }

  html += '</div>';
  entry.innerHTML = html;
  el.prepend(entry);
}

function esc(s) { return (s||'').replace(/&/g,'&amp;').replace(/</g,'&lt;').replace(/>/g,'&gt;'); }
</script>
</body>
</html>
`
