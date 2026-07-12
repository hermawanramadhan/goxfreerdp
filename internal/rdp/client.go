package rdp

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"goxfreerdp/internal/config"
)

// RDPError holds execution error and stderr stream contents.
type RDPError struct {
	Err    error
	Stderr string
}

func (e *RDPError) Error() string {
	return fmt.Sprintf("%v: %s", e.Err, e.Stderr)
}

// parseBoolOverride parses a per-server string override, falling back to a default boolean value if empty.
func parseBoolOverride(val string, defaultVal bool) bool {
	val = strings.ToLower(strings.TrimSpace(val))
	if val == "" || val == "default" {
		return defaultVal
	}
	return val == "yes" || val == "true" || val == "1" || val == "+" || val == "enable" || val == "enabled"
}

// parseStringOverride parses a per-server string override, falling back to a default string value if empty.
func parseStringOverride(val string, defaultVal string) string {
	val = strings.TrimSpace(val)
	if val == "" || strings.ToLower(val) == "default" {
		return defaultVal
	}
	return val
}

// BuildArgs maps configurations to CLI parameters for the selected engine.
func BuildArgs(settings config.SettingsConfig, server config.ServerConfig, overrideHost string) []string {
	var args []string

	// Log Level
	logLevel := settings.LogLevel
	if logLevel != "" && logLevel != "default" {
		args = append(args, "/log-level:"+logLevel)
	}

	// Target host selection: server specific IP or overridden path, falling back to default settings Host
	host := server.HostIP
	if overrideHost != "" {
		host = overrideHost
	}
	if host == "" {
		host = settings.Host
	}

	port := server.Port
	if port == "" {
		port = settings.Port
	}

	if host != "" {
		if strings.HasSuffix(strings.ToLower(host), ".rdp") {
			args = append(args, host)
		} else {
			if port != "" {
				args = append(args, "/v:"+host+":"+port)
			} else {
				args = append(args, "/v:"+host)
			}
		}
	}

	// Window Title
	title := server.Name
	if title == "" {
		title = host
	}
	if title != "" {
		args = append(args, "/t:"+title+" - GoXFreeRDP")
	}

	// Engine
	engine := parseStringOverride(server.Engine, settings.Engine)

	// Username: server username overrides default settings username
	username := server.Username
	if username == "" {
		username = settings.Username
	}

	if username != "" {
		args = append(args, "/u:"+username)
	} else if engine == "xfreerdp3" {
		args = append(args, "/u:")
	}

	// Password: server password overrides default settings password
	password := server.Password
	if password == "" {
		password = settings.Password
	}

	if password != "" {
		args = append(args, "/p:"+password)
	} else if engine == "xfreerdp3" {
		args = append(args, "/p:")
	}

	// Ignore Certificate
	ignoreCert := parseBoolOverride(server.IgnoreCertificate, settings.IgnoreCertificate)
	if ignoreCert {
		args = append(args, "/cert:ignore")
	}

	// TLS security level parameter handling depending on the engine
	tlsLevel := parseStringOverride(server.TLSSecLevel, settings.TLSSecLevel)
	if tlsLevel != "" && tlsLevel != "default" {
		if engine == "xfreerdp3" {
			args = append(args, "/tls:seclevel:"+tlsLevel)
		} else {
			args = append(args, "/tls-seclevel:"+tlsLevel)
		}
	}

	// Clipboard: Server specific clipboard settings overrides global settings if saved.
	useClipboard := parseBoolOverride(server.Clipboard, settings.Clipboard)
	if useClipboard {
		args = append(args, "+clipboard")
	} else {
		args = append(args, "-clipboard")
	}

	// Network Level Authentication (NLA)
	useNLA := parseBoolOverride(server.SecNLA, settings.SecNLA)
	if useNLA {
		if engine == "xfreerdp3" {
			args = append(args, "/sec:nla:on")
		} else {
			args = append(args, "+sec-nla")
		}
	} else {
		if engine == "xfreerdp3" {
			args = append(args, "/sec:nla:off")
		} else {
			args = append(args, "-sec-nla")
		}
	}

	// Fullscreen
	fullscreen := parseBoolOverride(server.Fullscreen, settings.Fullscreen)
	if fullscreen {
		args = append(args, "/f")
	} else {
		// If fullscreen is false and we are loading a .rdp file, we must pass /size to override
		// the RDP file's native fullscreen setting (screen mode id:i:2).
		if strings.HasSuffix(strings.ToLower(host), ".rdp") {
			args = append(args, "/size:85%")
		}
	}

	// Dynamic Resolution
	dynamicRes := parseBoolOverride(server.DynamicRes, settings.DynamicRes)
	if dynamicRes {
		args = append(args, "/dynamic-resolution")
	}

	// Multimon
	multimon := parseBoolOverride(server.Multimon, settings.Multimon)
	if multimon {
		args = append(args, "/multimon")
	}

	// Sound
	sound := parseBoolOverride(server.Sound, settings.Sound)
	if sound {
		args = append(args, "/sound")
	}

	// Share Home
	shareHome := parseBoolOverride(server.ShareHome, settings.ShareHome)
	if shareHome {
		homeDir, err := os.UserHomeDir()
		if err == nil {
			args = append(args, "/drive:home,"+homeDir)
		}
	}

	// Font Smoothing
	fontSmoothing := parseBoolOverride(server.FontSmoothing, settings.FontSmoothing)
	if fontSmoothing {
		args = append(args, "+fonts")
	} else {
		args = append(args, "-fonts")
	}

	// Wallpaper
	wallpaper := parseBoolOverride(server.Wallpaper, settings.Wallpaper)
	if wallpaper {
		args = append(args, "+wallpaper")
	} else {
		args = append(args, "-wallpaper")
	}

	// Themes
	themes := parseBoolOverride(server.Themes, settings.Themes)
	if themes {
		args = append(args, "+themes")
	} else {
		args = append(args, "-themes")
	}

	// Server specific custom parameters
	if server.CustomParams != "" {
		args = append(args, strings.Fields(server.CustomParams)...)
	}

	// Settings custom parameters
	if settings.CustomParams != "" {
		args = append(args, strings.Fields(settings.CustomParams)...)
	}

	return args
}

// CensorArgs replaces any "/p:<password>" arguments with "/p:******" for logging security.
func CensorArgs(args []string) []string {
	censored := make([]string, len(args))
	for i, arg := range args {
		if strings.HasPrefix(arg, "/p:") {
			censored[i] = "/p:******"
		} else {
			censored[i] = arg
		}
	}
	return censored
}

// LaunchRDP executes the xfreerdp/xfreerdp3 process asynchronously.
func LaunchRDP(settings config.SettingsConfig, server config.ServerConfig, overrideHost string, logWriter io.Writer) chan error {
	done := make(chan error, 1)

	go func() {
		args := BuildArgs(settings, server, overrideHost)

		// If we are loading a .rdp file, we rewrite the .rdp file to ensure the screen mode
		// matches the requested fullscreen setting (1 for windowed, 2 for fullscreen).
		fullscreen := parseBoolOverride(server.Fullscreen, settings.Fullscreen)
		for i, arg := range args {
			if strings.HasSuffix(strings.ToLower(arg), ".rdp") {
				modifiedPath, err := createTempModifiedRdpFile(arg, fullscreen)
				if err == nil {
					args[i] = modifiedPath
				} else {
					fmt.Printf("[RDP Launch] Warning: Failed to create temp modified RDP file: %v\n", err)
				}
				break
			}
		}

		engine := parseStringOverride(server.Engine, settings.Engine)
		if engine == "" {
			engine = "xfreerdp"
		}

		censored := CensorArgs(args)
		fmt.Printf("[RDP Launch] Executing: %s %s\n", engine, strings.Join(censored, " "))

		cmd := exec.Command(engine, args...)

		// Capture stderr to analyze for logon/authentication failures
		var stderrBuf bytes.Buffer

		if logWriter != nil {
			cmd.Stdout = io.MultiWriter(os.Stdout, logWriter)
			cmd.Stderr = io.MultiWriter(os.Stderr, logWriter, &stderrBuf)
		} else {
			cmd.Stdout = os.Stdout
			cmd.Stderr = io.MultiWriter(os.Stderr, &stderrBuf)
		}

		err := cmd.Start()
		if err != nil {
			fmt.Printf("[RDP Error] Failed to start RDP engine (%s): %v\n", engine, err)
			done <- &RDPError{Err: err, Stderr: err.Error()}
			close(done)
			return
		}

		err = cmd.Wait()
		if err != nil {
			fmt.Printf("[RDP Status] RDP session exited with error/status: %v\n", err)
			done <- &RDPError{Err: err, Stderr: stderrBuf.String()}
		} else {
			fmt.Println("[RDP Status] RDP session completed successfully.")
			done <- nil
		}
		close(done)
	}()

	return done
}

// ParseRdpFile parses a standard RDP text file for the 'full address:s:' entry.
func ParseRdpFile(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		if strings.HasPrefix(strings.ToLower(line), "full address:s:") {
			val := line[len("full address:s:"):]
			return strings.TrimSpace(val), nil
		}
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return "", fmt.Errorf("could not find 'full address:s:' field in .rdp file")
}

// createTempModifiedRdpFile reads an RDP file, alters screen mode id:i: to 1 (windowed) or 2 (fullscreen) based on the fullscreen flag, and writes a temp copy.
func createTempModifiedRdpFile(originalPath string, fullscreen bool) (string, error) {
	data, err := os.ReadFile(originalPath)
	if err != nil {
		return "", err
	}

	content := string(data)
	lines := strings.Split(content, "\n")
	hasScreenMode := false

	targetMode := "1"
	if fullscreen {
		targetMode = "2"
	}

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(strings.ToLower(trimmed), "screen mode id:i:") {
			lines[i] = "screen mode id:i:" + targetMode
			hasScreenMode = true
		}
	}

	if !hasScreenMode {
		lines = append(lines, "screen mode id:i:"+targetMode)
	}

	modifiedContent := strings.Join(lines, "\n")

	tempFile := filepath.Join(os.TempDir(), "goxfreerdp_temp.rdp")
	err = os.WriteFile(tempFile, []byte(modifiedContent), 0644)
	if err != nil {
		return "", err
	}

	return tempFile, nil
}
