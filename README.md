# GoXFreeRDP

GoXFreeRDP is a modern GTK3 GUI wrapper for `xfreerdp` (FreeRDP v2/v3) on Linux. It provides a simple, clean, and responsive user interface to manage remote desktop connections and seamlessly integrates with your system's file manager to open `.rdp` files on double-click.

---

## Features

- **Connection Management:** Save and organize multiple RDP servers with custom configurations.
- **Adaptive Dark Mode:** Automatically detects and syncs with your desktop's light or dark mode theme in real-time.
- **Credential Fallback:** Intercepts authentication errors and prompts you with a GTK dialog to input your password securely.
- **Advanced Parameter Overrides:** Fine-tune global and per-server overrides for clipboard sharing, audio redirection, fullscreen, dynamic resolution, multi-monitor, font smoothing, and custom command-line flags.
- **Collapsible Layout:** Keeps dialogs clean by grouping advanced options inside a collapsible drawer.
- **Interactive Connection Logs:** A built-in terminal-like console to view the output of RDP sessions in real-time for troubleshooting.
- **File Explorer Integration:** Registers MIME types to let you double-click `.rdp` files in file managers like Nautilus, Dolphin, or Thunar to launch connections instantly.

---

## Prerequisites

Before building or installing, ensure you have the required build tools and libraries installed.

> [!NOTE]
> **Go** and **GTK3 Development Headers** are only required during the installation process because the application is compiled/built from source on your machine. Once installed, the compiled binary runs natively without requiring them.

### Build Dependencies:
- **Go** (1.18 or higher)
- **pkg-config**
- **GTK 3 Development Headers**

### Runtime Dependency:
- **FreeRDP** (either `xfreerdp` or `xfreerdp3` binary in your PATH)

### Installing Prerequisites:

- **Ubuntu / Debian / Linux Mint:**
  ```bash
  sudo apt update
  sudo apt install golang-go libgtk-3-dev pkg-config freerdp2-x11
  ```

- **Fedora / RHEL:**
  ```bash
  sudo dnf install golang gtk3-devel pkg-config freerdp
  ```

- **Arch Linux:**
  ```bash
  sudo pacman -S go gtk3 pkgconf freerdp
  ```

---

## Installation

### 1. Automated Script
We provide an interactive installer script that lets you choose between:
1. **Downloading a precompiled release binary** from GitHub (Recommended - fast, does not require compilation tools like Go or GTK3 development headers).
2. **Compiling from source** (which will automatically detect and prompt to install any missing development dependencies like Go, GTK3 Dev Headers, etc. using your package manager).

The script also prompts you to install either for your **local user only** (does not require root privileges) or **system-wide for all users** (using `sudo`):
```bash
./install.sh
```

### 2. Manual Compilation & Installation (using Makefile)
To build the application without installing:
```bash
make
```

To install the application in user-space (`~/.local/bin`):
```bash
make install
```

To install the application system-wide (requires root privileges):
```bash
sudo make install PREFIX=/usr/local
```

---

## File Association

GoXFreeRDP registers itself as the default application to open `.rdp` files. When you install using `make install` or `./install.sh`:
1. It registers the `application/x-rdp` MIME type.
2. It installs the desktop application launcher (`goxfreerdp.desktop`).
3. It associates `.rdp` extensions with GoXFreeRDP.

You can now double-click any `.rdp` file in your desktop file manager to connect directly to the remote server.

---

## Uninstallation

To remove GoXFreeRDP and clean up desktop integrations:

- If installed to user-space:
  ```bash
  make uninstall
  ```

- If installed system-wide:
  ```bash
  sudo make uninstall PREFIX=/usr/local
  ```

---

## License

This project is licensed under the MIT License - see the LICENSE file for details.
