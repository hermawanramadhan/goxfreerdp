// 1. Copy to Clipboard Functionality
document.querySelectorAll('.copy-btn').forEach(button => {
    button.addEventListener('click', () => {
        const targetId = button.getAttribute('data-target');
        const targetElement = document.getElementById(targetId);
        if (!targetElement) return;

        let text = targetElement.innerText || targetElement.textContent;
        // Clean up leading comment indicators if any
        if (text.startsWith('#')) {
            // keep the text as is but copy-to-clipboard handles it
        }

        navigator.clipboard.writeText(text).then(() => {
            // Save original icon HTML
            const originalIcon = button.innerHTML;
            
            // Show checkmark icon
            button.innerHTML = `
                <svg viewBox="0 0 16 16" width="16" height="16">
                    <path fill="#10b981" fill-rule="evenodd" d="M13.78 4.22a.75.75 0 010 1.06l-7.25 7.25a.75.75 0 01-1.06 0L2.22 9.28a.75.75 0 011.06-1.06L6 10.94l6.72-6.72a.75.75 0 011.06 0z"></path>
                </svg>
            `;
            button.style.borderColor = '#10b981';

            // Reset icon after 2 seconds
            setTimeout(() => {
                button.innerHTML = originalIcon;
                button.style.borderColor = '';
            }, 2000);
        }).catch(err => {
            console.error('Failed to copy text: ', err);
        });
    });
});

// 2. Installation Tabs Switching
const installTabHeaders = document.querySelectorAll('.install-tab-hdr');
const installTabContents = document.querySelectorAll('.install-tab-content');

installTabHeaders.forEach(header => {
    header.addEventListener('click', () => {
        // Remove active class from all headers and contents
        installTabHeaders.forEach(hdr => hdr.classList.remove('active'));
        installTabContents.forEach(content => content.classList.remove('active'));

        // Add active class to clicked header
        header.classList.add('active');

        // Show target content
        const targetId = header.getAttribute('data-tab');
        const targetContent = document.getElementById(targetId);
        if (targetContent) {
            targetContent.classList.add('active');
        }
    });
});

// 3. Mock App Tab Switching
const tabBtnServers = document.getElementById('tab-btn-servers');
const tabBtnLogs = document.getElementById('tab-btn-logs');
const tabContentServers = document.getElementById('tab-content-servers');
const tabContentLogs = document.getElementById('tab-content-logs');

function switchAppTab(tabName) {
    if (tabName === 'servers') {
        tabBtnServers.classList.add('active');
        tabBtnLogs.classList.remove('active');
        tabContentServers.classList.add('active');
        tabContentLogs.classList.remove('active');
    } else {
        tabBtnLogs.classList.add('active');
        tabBtnServers.classList.remove('active');
        tabContentLogs.classList.add('active');
        tabContentServers.classList.remove('active');
    }
}

tabBtnServers.addEventListener('click', () => switchAppTab('servers'));
tabBtnLogs.addEventListener('click', () => switchAppTab('logs'));

// 4. Overlays & Dialogs Management
const dialogAbout = document.getElementById('dialog-about');
const dialogAddServer = document.getElementById('dialog-add-server');
const dialogFallback = document.getElementById('dialog-fallback');

function openDialog(dialog) {
    dialog.classList.add('active');
}

function closeDialog(dialog) {
    dialog.classList.remove('active');
}

// About Dialog
document.getElementById('btn-about-app').addEventListener('click', () => openDialog(dialogAbout));
document.getElementById('btn-close-about').addEventListener('click', () => closeDialog(dialogAbout));

// Add Server Dialog
document.getElementById('btn-add-server').addEventListener('click', () => {
    // Reset form fields
    document.getElementById('srv-name').value = '';
    document.getElementById('srv-host').value = '';
    document.getElementById('srv-user').value = '';
    openDialog(dialogAddServer);
});

const closeAddDialogElements = [
    document.getElementById('btn-close-add-dialog'),
    document.getElementById('btn-cancel-add')
];
closeAddDialogElements.forEach(btn => {
    btn.addEventListener('click', () => closeDialog(dialogAddServer));
});

// Advanced Drawer Toggle in Add Server Dialog
const drawerToggleBtn = document.getElementById('drawer-toggle-btn');
const advancedDrawerContent = document.getElementById('advanced-drawer-content');

drawerToggleBtn.addEventListener('click', () => {
    const isExpanded = drawerToggleBtn.classList.toggle('expanded');
    advancedDrawerContent.classList.toggle('expanded');
});

// Save Server Logic (Dynamic list addition in mockup)
document.getElementById('btn-save-server').addEventListener('click', () => {
    const name = document.getElementById('srv-name').value.trim() || 'Remote Server';
    const host = document.getElementById('srv-host').value.trim() || '192.168.1.1';
    const username = document.getElementById('srv-user').value.trim() || 'user';
    
    // Create new server row markup
    const listBox = document.querySelector('.gtk-list-box');
    const newRow = document.createElement('div');
    newRow.className = 'gtk-list-row';
    
    // Generate simple ID
    const rowId = 'row-' + Date.now();
    newRow.id = rowId;
    
    newRow.innerHTML = `
        <div class="gtk-row-left">
            <span class="gtk-server-icon">🖥️</span>
            <div class="gtk-server-details">
                <div class="gtk-server-name">${name}</div>
                <div class="gtk-server-badges">
                    <span class="gtk-host">${host}</span>
                    <span class="gtk-badge">👤 ${username}</span>
                    <span class="gtk-badge">⚙️ xfreerdp</span>
                </div>
            </div>
        </div>
        <button class="gtk-btn-connect" title="Connect to ${name}">
            <svg viewBox="0 0 16 16" width="12" height="12">
                <path fill="currentColor" d="M3 2v12l10-6z"></path>
            </svg>
        </button>
    `;
    
    // Bind connect handler to the new play button
    newRow.querySelector('.gtk-btn-connect').addEventListener('click', (e) => {
        e.stopPropagation();
        simulateConnection(name, host, username, 'xfreerdp');
    });
    
    // Append to list
    listBox.appendChild(newRow);
    
    // Close Dialog
    closeDialog(dialogAddServer);
});

// Settings & Open RDP mock popups
document.getElementById('btn-settings').addEventListener('click', () => {
    alert('GoXFreeRDP Settings: Adaptive Dark Mode toggle, default engine parameters, and general application configuration shortcuts.');
});

document.getElementById('btn-open-rdp').addEventListener('click', () => {
    alert('Opens your local File Chooser to double-click any .rdp configuration file. Integrates connection configuration directly.');
});

// 5. Connection Simulator & Credential Fallback Demo
const logConsole = document.getElementById('log-console');
let simulatedTimeoutIds = [];
let pendingConnection = null;

function appendLogLine(text, className = '') {
    const line = document.createElement('div');
    line.className = `log-line ${className}`;
    line.textContent = text;
    logConsole.appendChild(line);
    logConsole.scrollTop = logConsole.scrollHeight;
}

function clearSimulatedLogs() {
    // Clear timeouts
    simulatedTimeoutIds.forEach(id => clearTimeout(id));
    simulatedTimeoutIds = [];
    logConsole.innerHTML = '';
}

function simulateConnection(name, host, username, engine) {
    // Switch to log tab and clear current logs
    switchAppTab('logs');
    clearSimulatedLogs();
    
    appendLogLine(`[GUI] Starting RDP connection to ${name} (${host})...`, 'gui');
    
    // Queue lines
    scheduleLogLine(`[GUI] Engine override: ${engine}`, 'gui', 500);
    scheduleLogLine(`[DEBUG] Executing command: ${engine} /v:${host} /u:${username} /gfx:on +clipboard +fonts /sound:sys:alsa /network:lan`, 'cmd', 1000);
    scheduleLogLine(`[INFO] Connection status: Connecting to RDP host ${host}:3389...`, 'info', 1800);
    scheduleLogLine(`[INFO] SSL established. Cipher suite: TLS_AES_256_GCM_SHA384`, 'info', 2300);
    
    // Trigger Credential Fallback Popup Mockup
    simulatedTimeoutIds.push(setTimeout(() => {
        appendLogLine(`[WARNING] Authentication requested credentials (ERRCONNECT_PASSWORD_EXPIRED or empty credentials).`, 'error');
        appendLogLine(`[GUI] Intercepting authentication challenge. Spawning credential dialog fallback...`, 'gui');
        
        // Save connection state
        pendingConnection = { name, host, username, engine };
        
        // Show fallback modal
        document.getElementById('fallback-srv-name').textContent = name;
        document.getElementById('fallback-username').value = username;
        document.getElementById('fallback-password').value = '';
        openDialog(dialogFallback);
        
    }, 3200));
}

function scheduleLogLine(text, className, delay) {
    const timeoutId = setTimeout(() => {
        appendLogLine(text, className);
    }, delay);
    simulatedTimeoutIds.push(timeoutId);
}

// Fallback Dialog Buttons
document.getElementById('btn-cancel-fallback').addEventListener('click', () => {
    closeDialog(dialogFallback);
    appendLogLine(`[GUI] Credential dialog cancelled by user.`, 'gui');
    appendLogLine(`[ERROR] Authentication failed. Exiting RDP session.`, 'error');
    appendLogLine(`[GUI] RDP process terminated with exit code 131.`, 'gui');
    pendingConnection = null;
});

document.getElementById('btn-submit-fallback').addEventListener('click', () => {
    const password = document.getElementById('fallback-password').value;
    const finalUser = document.getElementById('fallback-username').value;
    
    closeDialog(dialogFallback);
    
    if (!password) {
        appendLogLine(`[GUI] Resumed connection, but password input was left empty.`, 'gui');
        appendLogLine(`[ERROR] Logon failed: empty credentials provided.`, 'error');
        appendLogLine(`[GUI] RDP process terminated with exit code 12.`, 'gui');
        return;
    }
    
    appendLogLine(`[GUI] Credentials captured securely. Re-executing wrapper with fallbacks...`, 'gui');
    appendLogLine(`[INFO] Authenticating as ${finalUser}...`, 'info');
    
    scheduleLogLine(`[INFO] Credential challenge accepted. Local credentials validated.`, 'info', 800);
    scheduleLogLine(`[INFO] Loading virtual channels: clipboard, rdpsnd, sound.`, 'info', 1300);
    scheduleLogLine(`[INFO] Connected successfully! Session is now active.`, 'info', 1800);
    scheduleLogLine(`[GUI] Connection monitored. Press Ctrl+Alt+Enter to exit.`, 'gui', 2200);
    
    pendingConnection = null;
});
