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

Depending on whether you use the precompiled binary or compile from source, you need different packages installed on your system.

### 1. If using the Precompiled Binary (Recommended)
You only need the **runtime** dependencies. Golang and development headers are **NOT** required:
* **FreeRDP** (either `xfreerdp` or `xfreerdp3` binary in your PATH)
* **xdg-utils** (for MIME file association)

Install runtime dependencies:
* **Ubuntu / Debian / Linux Mint:**
  ```bash
  sudo apt update
  sudo apt install freerdp2-x11 xdg-utils
  ```
* **Fedora / RHEL:**
  ```bash
  sudo dnf install freerdp xdg-utils
  ```
* **Arch Linux:**
  ```bash
  sudo pacman -S freerdp xdg-utils
  ```

### 2. If Compiling from Source
You need both the **build tools** (only required during compilation) and **runtime** dependencies:
* **Go** (1.18 or higher)
* **pkg-config**
* **GTK 3 Development Headers**
* **FreeRDP** & **xdg-utils** (runtime)

Install build + runtime dependencies:
* **Ubuntu / Debian / Linux Mint:**
  ```bash
  sudo apt update
  sudo apt install golang-go libgtk-3-dev pkg-config freerdp2-x11 xdg-utils
  ```
* **Fedora / RHEL:**
  ```bash
  sudo dnf install golang gtk3-devel pkg-config freerdp xdg-utils
  ```
* **Arch Linux:**
  ```bash
  sudo pacman -S go gtk3 pkgconf freerdp xdg-utils
  ```

---

## Installation

### 1. Automated Script (Recommended)
We provide an interactive installer script that automatically detects and prompts you to install any missing runtime/development dependencies (such as FreeRDP, xdg-utils, Go, or GTK3 headers) using your package manager. 

It lets you choose between:
1. **Downloading a precompiled release binary** from GitHub (Recommended - fast, does not require compilation tools like Go or GTK3 development headers).
2. **Compiling from source** (builds GoXFreeRDP locally on your machine).

It also configures the desktop menu entries, icons, and `.rdp` file associations:
```bash
git clone https://github.com/hermawanramadhan/goxfreerdp.git
cd goxfreerdp
chmod +x install.sh
./install.sh
```

### 2. Manual Compilation & Installation (using Makefile)
If you prefer compiling the application manually from source:

* **Build only** (binary remains in the current directory):
  ```bash
  make
  ```

* **Install for local user only** (installs to `~/.local/bin`):
  ```bash
  make install
  ```

* **Install system-wide for all users** (requires root privileges):
  ```bash
  sudo make install PREFIX=/usr/local
  ```

### 3. Manual Installation of Precompiled Binary
If you downloaded a precompiled binary (`goxfreerdp-linux-amd64` or `goxfreerdp-linux-arm64`) from the GitHub Releases page:

1. Rename the downloaded binary to **`goxfreerdp`** and place it in the root of the cloned repository.
2. Grant executable permissions to the binary:
   ```bash
   chmod +x goxfreerdp
   ```
3. Install the binary along with the desktop launcher, icons, and MIME associations:
   * **For local user only**:
     ```bash
     make install-only
     ```
   * **System-wide for all users**:
     ```bash
     sudo make install-only PREFIX=/usr/local
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
