.PHONY: all build install install-only uninstall clean

PREFIX ?= $(HOME)/.local
BINDIR = $(PREFIX)/bin
SHAREDIR = $(PREFIX)/share
APPDIR = $(SHAREDIR)/applications
MIMEDIR = $(SHAREDIR)/mime/packages

TARGET = goxfreerdp
SRC = cmd/goxfreerdp/main.go

all: build

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS = -ldflags "-X 'goxfreerdp/internal/ui.Version=$(VERSION)'"

build:
	@echo "Building $(TARGET)..."
	go build $(LDFLAGS) -o $(TARGET) $(SRC)

install: build install-only

install-only:
	@echo "Installing binary to $(BINDIR)..."
	mkdir -p $(BINDIR)
	cp $(TARGET) $(BINDIR)/$(TARGET)
	
	@echo "Installing MIME package..."
	mkdir -p $(MIMEDIR)
	cp resources/goxfreerdp.xml $(MIMEDIR)/goxfreerdp.xml
	update-mime-database $(SHAREDIR)/mime
	
	@echo "Installing desktop entry to $(APPDIR)..."
	mkdir -p $(APPDIR)
	sed "s|@BINDIR@|$(BINDIR)|g" resources/goxfreerdp.desktop.template > $(APPDIR)/goxfreerdp.desktop
	update-desktop-database $(APPDIR)
	
	@echo "Registering default MIME handler..."
	xdg-mime default goxfreerdp.desktop application/x-rdp
	
	@echo "Installation complete!"

uninstall:
	@echo "Removing binary..."
	rm -f $(BINDIR)/$(TARGET)
	
	@echo "Removing MIME package..."
	rm -f $(MIMEDIR)/goxfreerdp.xml
	update-mime-database $(SHAREDIR)/mime
	
	@echo "Removing desktop entry..."
	rm -f $(APPDIR)/goxfreerdp.desktop
	update-desktop-database $(APPDIR)
	
	@echo "Uninstallation complete!"

clean:
	@echo "Cleaning build artifacts..."
	rm -f $(TARGET)
