#!/bin/bash
# GoXFreeRDP Installation Script
# Verifies system requirements and installs the app to user space (~/.local)

set -e

echo "==========================================="
echo "        GoXFreeRDP Installer"
echo "==========================================="

# Dependency checking helper
check_dep() {
  if ! command -v "$1" &>/dev/null; then
    echo "Error: '$1' is required but not installed." >&2
    return 1
  fi
}

# Check build dependencies
echo "Checking dependencies..."
check_dep go || { echo "Please install Go: https://go.dev/doc/install"; exit 1; }
check_dep pkg-config || { echo "Please install pkg-config using your package manager."; exit 1; }
check_dep xdg-mime || { echo "xdg-mime is required for desktop file association."; exit 1; }

# Check GTK 3 development headers
if ! pkg-config --exists gtk+-3.0; then
  echo "" >&2
  echo "Error: GTK 3 development headers (gtk+-3.0) were not found." >&2
  echo "Please install GTK 3 development files:" >&2
  echo " - Debian/Ubuntu: sudo apt install libgtk-3-dev" >&2
  echo " - Fedora/CentOS: sudo dnf install gtk3-devel" >&2
  echo " - Arch Linux:    sudo pacman -S gtk3" >&2
  exit 1
fi

# Check RDP execution engine (xfreerdp)
if ! command -v xfreerdp &>/dev/null && ! command -v xfreerdp3 &>/dev/null; then
  echo "" >&2
  echo "Warning: FreeRDP ('xfreerdp' or 'xfreerdp3') was not found in your PATH." >&2
  echo "You must install FreeRDP to launch RDP connections." >&2
  echo " - Debian/Ubuntu: sudo apt install freerdp2-x11" >&2
  echo " - Fedora/CentOS: sudo dnf install freerdp" >&2
  echo " - Arch Linux:    sudo pacman -S freerdp" >&2
  echo "" >&2
fi

# Execute Makefile installation
echo "Compiling and installing GoXFreeRDP..."
make install

echo ""
echo "==========================================="
echo " 🎉 GoXFreeRDP Installed Successfully!"
echo "==========================================="
echo " - You can start the app from your application menu or terminal by running: goxfreerdp"
echo " - You can now open .rdp files directly from your File Explorer by double-clicking them."
echo ""
