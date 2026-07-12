package ui

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"

	"github.com/gotk3/gotk3/glib"
	"goxfreerdp/internal/config"
)

// GetSocketPath returns the Unix domain socket path for single-instance IPC.
func GetSocketPath() string {
	runtimeDir := os.Getenv("XDG_RUNTIME_DIR")
	if runtimeDir == "" {
		runtimeDir = os.TempDir()
	}
	return filepath.Join(runtimeDir, "goxfreerdp.sock")
}

// TrySingleInstance checks if another instance is running.
// If it is, it sends the file path to that instance and returns true.
// Otherwise, it returns false.
func TrySingleInstance(initialRDPFile string) bool {
	socketPath := GetSocketPath()
	conn, err := net.Dial("unix", socketPath)
	if err != nil {
		// No other instance running or socket is dead
		return false
	}
	defer conn.Close()

	// Send the file path (or empty string/command)
	payload := initialRDPFile
	if payload == "" {
		payload = "ACTIVATE"
	}
	fmt.Fprintln(conn, payload)
	return true
}

// StartIPCServer starts a Unix socket server to listen for new connection requests.
func (app *AppUI) StartIPCServer() error {
	socketPath := GetSocketPath()

	// Remove leftover socket if any
	_ = os.Remove(socketPath)

	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		return fmt.Errorf("failed to listen on socket: %w", err)
	}

	go func() {
		// Clean up socket when listener stops
		defer listener.Close()
		defer os.Remove(socketPath)

		for {
			conn, err := listener.Accept()
			if err != nil {
				return
			}
			go app.handleIPCClient(conn)
		}
	}()

	return nil
}

func (app *AppUI) handleIPCClient(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	line, err := reader.ReadString('\n')
	if err != nil {
		return
	}

	cmd := strings.TrimSpace(line)
	if cmd == "" {
		return
	}

	glib.IdleAdd(func() {
		if app.Window != nil {
			app.Window.Present()
		}

		if cmd != "ACTIVATE" {
			var dummyServer config.ServerConfig
			dummyServer.HostIP = cmd
			dummyServer.Name = filepath.Base(cmd)

			if app.LogTextBuffer != nil {
				endIter := app.LogTextBuffer.GetEndIter()
				msg := fmt.Sprintf("[IPC] Received RDP connection request from another instance for file: %s\n", cmd)
				app.LogTextBuffer.Insert(endIter, msg)
			}

			app.RunConnectionWithAuthFallback(dummyServer)
		}
	})
}
