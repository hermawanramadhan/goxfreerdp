// Default Settings State (mirrors internal/config app.go defaults)
let appSettings = {
    engine: 'xfreerdp3',
    ignoreCertificate: true,
    tlsLevel: 'default',
    port: '3389',
    clipboard: true,
    nla: true,
    sound: true,
    shareHome: false,
    fullscreen: false,
    dynamicRes: true,
    multimon: false,
    fontSmoothing: true,
    wallpaper: false,
    themes: false,
    customParams: '',
    logLevel: 'default'
};

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
const dialogSettings = document.getElementById('dialog-settings');

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

// Settings Dialog Trigger
document.getElementById('btn-settings').addEventListener('click', () => {
    // Load config values into UI fields
    document.getElementById('cfg-engine').value = appSettings.engine;
    document.getElementById('cfg-ignore-cert').checked = appSettings.ignoreCertificate;
    document.getElementById('cfg-tls-level').value = appSettings.tlsLevel;
    document.getElementById('cfg-port').value = appSettings.port;
    
    document.getElementById('cfg-clipboard').checked = appSettings.clipboard;
    document.getElementById('cfg-nla').checked = appSettings.nla;
    document.getElementById('cfg-sound').checked = appSettings.sound;
    document.getElementById('cfg-drive').checked = appSettings.shareHome;
    
    document.getElementById('cfg-fullscreen').checked = appSettings.fullscreen;
    document.getElementById('cfg-dynamic-res').checked = appSettings.dynamicRes;
    document.getElementById('cfg-multimon').checked = appSettings.multimon;
    
    document.getElementById('cfg-fonts').checked = appSettings.fontSmoothing;
    document.getElementById('cfg-wallpaper').checked = appSettings.wallpaper;
    document.getElementById('cfg-themes').checked = appSettings.themes;
    
    document.getElementById('cfg-custom-params').value = appSettings.customParams;
    document.getElementById('cfg-log-level').value = appSettings.logLevel;
    
    openDialog(dialogSettings);
});

// Close Settings Dialog buttons
const closeSettingsElements = [
    document.getElementById('btn-close-settings-dialog'),
    document.getElementById('btn-cancel-settings')
];
closeSettingsElements.forEach(btn => {
    btn.addEventListener('click', () => closeDialog(dialogSettings));
});

// Save Settings Dialog logic
document.getElementById('btn-save-settings').addEventListener('click', () => {
    // Save UI values to config
    appSettings.engine = document.getElementById('cfg-engine').value;
    appSettings.ignoreCertificate = document.getElementById('cfg-ignore-cert').checked;
    appSettings.tlsLevel = document.getElementById('cfg-tls-level').value;
    appSettings.port = document.getElementById('cfg-port').value.trim() || '3389';
    
    appSettings.clipboard = document.getElementById('cfg-clipboard').checked;
    appSettings.nla = document.getElementById('cfg-nla').checked;
    appSettings.sound = document.getElementById('cfg-sound').checked;
    appSettings.shareHome = document.getElementById('cfg-drive').checked;
    
    appSettings.fullscreen = document.getElementById('cfg-fullscreen').checked;
    appSettings.dynamicRes = document.getElementById('cfg-dynamic-res').checked;
    appSettings.multimon = document.getElementById('cfg-multimon').checked;
    
    appSettings.fontSmoothing = document.getElementById('cfg-fonts').checked;
    appSettings.wallpaper = document.getElementById('cfg-wallpaper').checked;
    appSettings.themes = document.getElementById('cfg-themes').checked;
    
    appSettings.customParams = document.getElementById('cfg-custom-params').value.trim();
    appSettings.logLevel = document.getElementById('cfg-log-level').value;
    
    closeDialog(dialogSettings);
    
    // Quick notification banner log printout
    switchAppTab('logs');
    clearSimulatedLogs();
    appendLogLine(`[GUI] Default Settings Saved: Engine: ${appSettings.engine}, Port: ${appSettings.port}, Clipboard: ${appSettings.clipboard ? 'ON' : 'OFF'}`, 'gui');
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
    
    // Fallback to configured global engine if default or blank
    const activeEngine = (engine && engine !== 'default') ? engine : appSettings.engine;
    const activePort = appSettings.port !== '3389' ? `:${appSettings.port}` : '';
    
    // Compile RDP flags dynamically from settings
    let cmdFlags = `/v:${host}${activePort}`;
    if (username) cmdFlags += ` /u:${username}`;
    if (appSettings.ignoreCertificate) cmdFlags += ` /cert:ignore`;
    if (appSettings.tlsLevel !== 'default') cmdFlags += ` /tls-seclevel:${appSettings.tlsLevel}`;
    if (appSettings.nla) cmdFlags += ` /sec:nla`;
    if (appSettings.clipboard) cmdFlags += ` +clipboard`;
    if (appSettings.sound) cmdFlags += ` /sound:sys:alsa`;
    if (appSettings.shareHome) cmdFlags += ` /drive:home,/home/${username || 'user'}`;
    if (appSettings.fullscreen) cmdFlags += ` /f`;
    if (appSettings.dynamicRes) cmdFlags += ` /dynamic-resolution`;
    if (appSettings.multimon) cmdFlags += ` /multimon`;
    if (appSettings.fontSmoothing) cmdFlags += ` +fonts`;
    if (appSettings.wallpaper) cmdFlags += ` +wallpaper`;
    if (appSettings.themes) cmdFlags += ` +window-drag`;
    if (appSettings.logLevel !== 'default') cmdFlags += ` /log-level:${appSettings.logLevel}`;
    if (appSettings.customParams) cmdFlags += ` ${appSettings.customParams}`;
    
    // Queue lines
    scheduleLogLine(`[GUI] Engine resolved: ${activeEngine}`, 'gui', 500);
    scheduleLogLine(`[DEBUG] Executing command: ${activeEngine} ${cmdFlags}`, 'cmd', 1000);
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
