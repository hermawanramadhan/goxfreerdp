#!/bin/bash
# GoXFreeRDP Installation Script
# Verifies system requirements and installs the app to user space (~/.local)

set -e

echo "==========================================="
echo "        GoXFreeRDP Installer"
echo "==========================================="

# Detect package manager
PM=""
if [ -f /etc/debian_version ] || command -v apt-get &>/dev/null; then
  PM="apt"
elif [ -f /etc/fedora-release ] || [ -f /etc/redhat-release ] || command -v dnf &>/dev/null; then
  PM="dnf"
elif [ -f /etc/arch-release ] || command -v pacman &>/dev/null; then
  PM="pacman"
fi

MISSING_PKGS=()

# Check Go
if ! command -v go &>/dev/null; then
  if [ "$PM" = "apt" ]; then MISSING_PKGS+=("golang-go"); fi
  if [ "$PM" = "dnf" ]; then MISSING_PKGS+=("golang"); fi
  if [ "$PM" = "pacman" ]; then MISSING_PKGS+=("go"); fi
fi

# Check pkg-config
if ! command -v pkg-config &>/dev/null; then
  if [ "$PM" = "apt" ]; then MISSING_PKGS+=("pkg-config"); fi
  if [ "$PM" = "dnf" ]; then MISSING_PKGS+=("pkg-config"); fi
  if [ "$PM" = "pacman" ]; then MISSING_PKGS+=("pkgconf"); fi
fi

# Check GTK 3 Dev Headers
if command -v pkg-config &>/dev/null; then
  if ! pkg-config --exists gtk+-3.0 2>/dev/null; then
    if [ "$PM" = "apt" ]; then MISSING_PKGS+=("libgtk-3-dev"); fi
    if [ "$PM" = "dnf" ]; then MISSING_PKGS+=("gtk3-devel"); fi
    if [ "$PM" = "pacman" ]; then MISSING_PKGS+=("gtk3"); fi
  fi
else
  if [ "$PM" = "apt" ]; then MISSING_PKGS+=("libgtk-3-dev"); fi
  if [ "$PM" = "dnf" ]; then MISSING_PKGS+=("gtk3-devel"); fi
  if [ "$PM" = "pacman" ]; then MISSING_PKGS+=("gtk3"); fi
fi

# Check FreeRDP
if ! command -v xfreerdp &>/dev/null && ! command -v xfreerdp3 &>/dev/null; then
  if [ "$PM" = "apt" ]; then MISSING_PKGS+=("freerdp2-x11"); fi
  if [ "$PM" = "dnf" ]; then MISSING_PKGS+=("freerdp"); fi
  if [ "$PM" = "pacman" ]; then MISSING_PKGS+=("freerdp"); fi
fi

# Check xdg-mime
if ! command -v xdg-mime &>/dev/null; then
  if [ "$PM" = "apt" ]; then MISSING_PKGS+=("xdg-utils"); fi
  if [ "$PM" = "dnf" ]; then MISSING_PKGS+=("xdg-utils"); fi
  if [ "$PM" = "pacman" ]; then MISSING_PKGS+=("xdg-utils"); fi
fi

if [ ${#MISSING_PKGS[@]} -gt 0 ]; then
  echo "The following missing system dependencies are required to build and run GoXFreeRDP:"
  for pkg in "${MISSING_PKGS[@]}"; do
    echo "  - $pkg"
  done
  echo ""
  if [ -z "$PM" ]; then
    echo "Error: Could not auto-detect package manager. Please install the packages listed above manually."
    exit 1
  fi

  read -p "Would you like to install them automatically now? (y/n) [y]: " INSTALL_DEP
  INSTALL_DEP=${INSTALL_DEP:-y}
  if [ "$INSTALL_DEP" = "y" ] || [ "$INSTALL_DEP" = "Y" ]; then
    echo "Installing missing dependencies (using sudo)..."
    if [ "$PM" = "apt" ]; then
      sudo apt-get update
      sudo apt-get install -y "${MISSING_PKGS[@]}"
    elif [ "$PM" = "dnf" ]; then
      sudo dnf install -y "${MISSING_PKGS[@]}"
    elif [ "$PM" = "pacman" ]; then
      sudo pacman -S --noconfirm "${MISSING_PKGS[@]}"
    fi
  else
    echo "Aborting. Missing dependencies are required to compile GoXFreeRDP."
    exit 1
  fi
fi

# Execute Makefile installation
echo ""
echo "Select installation target:"
echo " 1) Local user only (installs to ~/.local, does not require root privileges)"
echo " 2) All users (installs to /usr/local, requires sudo/root privileges)"
echo ""
read -p "Choose option [1]: " INSTALL_OPT
INSTALL_OPT=${INSTALL_OPT:-1}

echo "Compiling and installing GoXFreeRDP..."
if [ "$INSTALL_OPT" = "2" ]; then
  echo "Installing system-wide for all users (using sudo)..."
  sudo make install PREFIX=/usr/local
else
  echo "Installing for local user..."
  make install PREFIX="$HOME/.local"
fi

echo ""
echo "==========================================="
echo " 🎉 GoXFreeRDP Installed Successfully!"
echo "==========================================="
echo " - You can start the app from your application menu or terminal by running: goxfreerdp"
echo " - You can now open .rdp files directly from your File Explorer by double-clicking them."
echo ""
