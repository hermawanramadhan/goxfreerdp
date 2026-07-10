package rdp

import (
	"reflect"
	"testing"
	"goxfreerdp/internal/config"
)

func TestBuildArgs(t *testing.T) {
	settings := config.SettingsConfig{
		Engine:            "xfreerdp3",
		Host:              "default-host",
		Username:          "default-user",
		Password:          "default-pass",
		IgnoreCertificate: true,
		TLSSecLevel:       "3",
		Clipboard:         true,
		SecNLA:            true,
		Fullscreen:        true,
		DynamicRes:        false,
		Multimon:          true,
		Sound:             false,
		ShareHome:         false,
		FontSmoothing:     true,
		Wallpaper:         false,
		Themes:            false,
		CustomParams:      "/custom-setting",
	}

	server := config.ServerConfig{
		ID:                "srv-001",
		Name:              "Test Server",
		HostIP:            "10.0.0.1",
		Username:          "server-user",
		Engine:            "xfreerdp", // override engine to v2
		IgnoreCertificate: "no",       // override to disable ignore certificate
		TLSSecLevel:       "1",        // override TLS level
		Clipboard:         "default",  // use default settings (true)
		SecNLA:            "default",  // use default settings (true)
		Fullscreen:        "no",       // override to disable fullscreen
		DynamicRes:        "yes",      // override to enable dynamic resolution
		Multimon:          "default",  // use default (true)
		Sound:             "yes",      // override to enable sound
		ShareHome:         "default",  // use default (false)
		FontSmoothing:     "no",       // override to disable font smoothing
		Wallpaper:         "default",  // use default (false)
		Themes:            "default",  // use default (false)
		CustomParams:      "/custom-server",
	}

	expected := []string{
		"/v:10.0.0.1", // server host
		"/t:Test Server - GoXFreeRDP", // window title
		"/u:server-user", // server username
		"/p:default-pass", // settings password fallback
		"/tls-seclevel:1", // overridden tls level with engine xfreerdp syntax
		"+clipboard", // default clipboard fallback
		"+sec-nla", // default sec-nla fallback
		"/dynamic-resolution", // overridden dynamic resolution
		"/multimon", // default multimon fallback
		"/sound", // overridden sound
		"-fonts", // overridden font smoothing disabled
		"-wallpaper", // default wallpaper fallback
		"-themes", // default themes fallback
		"/custom-server", // server custom params
		"/custom-setting", // settings custom params
	}

	args := BuildArgs(settings, server, "")

	if !reflect.DeepEqual(args, expected) {
		t.Errorf("Args mismatch.\nExpected: %v\nGot:      %v", expected, args)
	}
}

func TestBuildArgsWithRdpFile(t *testing.T) {
	settings := config.SettingsConfig{
		Engine:            "xfreerdp",
		IgnoreCertificate: true,
		Clipboard:         true,
		SecNLA:            true,
	}

	server := config.ServerConfig{
		HostIP:    "/path/to/connection.rdp",
		Name:      "connection.rdp",
		Clipboard: "default",
		SecNLA:    "default",
	}

	expected := []string{
		"/path/to/connection.rdp",
		"/t:connection.rdp - GoXFreeRDP",
		"/cert:ignore",
		"+clipboard",
		"+sec-nla",
		"/size:85%",
		"-fonts",
		"-wallpaper",
		"-themes",
	}

	args := BuildArgs(settings, server, "")

	if !reflect.DeepEqual(args, expected) {
		t.Errorf("RDP file args mismatch.\nExpected: %v\nGot:      %v", expected, args)
	}
}

func TestBuildArgsWithCustomPort(t *testing.T) {
	settings := config.SettingsConfig{
		Engine:            "xfreerdp",
		IgnoreCertificate: true,
		Clipboard:         true,
		SecNLA:            true,
		Port:              "3389",
	}

	server := config.ServerConfig{
		HostIP:    "192.168.1.50",
		Port:      "13389",
		Name:      "Custom Port Server",
		Clipboard: "default",
		SecNLA:    "default",
	}

	expected := []string{
		"/v:192.168.1.50:13389",
		"/t:Custom Port Server - GoXFreeRDP",
		"/cert:ignore",
		"+clipboard",
		"+sec-nla",
		"-fonts",
		"-wallpaper",
		"-themes",
	}

	args := BuildArgs(settings, server, "")

	if !reflect.DeepEqual(args, expected) {
		t.Errorf("Custom port args mismatch.\nExpected: %v\nGot:      %v", expected, args)
	}
}

func TestBuildArgsWithNLADisabled(t *testing.T) {
	settings := config.SettingsConfig{
		Engine:    "xfreerdp",
		Clipboard: true,
		SecNLA:    true, // default is true
	}

	server := config.ServerConfig{
		HostIP:    "192.168.1.50",
		Name:      "No NLA Server",
		Clipboard: "default",
		SecNLA:    "no", // override to disable
	}

	expected := []string{
		"/v:192.168.1.50",
		"/t:No NLA Server - GoXFreeRDP",
		"+clipboard",
		"-sec-nla", // disabled NLA
		"-fonts",
		"-wallpaper",
		"-themes",
	}

	args := BuildArgs(settings, server, "")

	if !reflect.DeepEqual(args, expected) {
		t.Errorf("Disabled NLA args mismatch.\nExpected: %v\nGot:      %v", expected, args)
	}
}

func TestBuildArgsWithXFreeRDP3NLA(t *testing.T) {
	settings := config.SettingsConfig{
		Engine:    "xfreerdp3",
		Clipboard: true,
		SecNLA:    true,
	}

	server := config.ServerConfig{
		HostIP:    "192.168.1.50",
		Name:      "NLA Server xfreerdp3",
		Clipboard: "default",
		SecNLA:    "default",
	}

	expected := []string{
		"/v:192.168.1.50",
		"/t:NLA Server xfreerdp3 - GoXFreeRDP",
		"/u:",
		"/p:",
		"+clipboard",
		"/sec:nla:on",
		"-fonts",
		"-wallpaper",
		"-themes",
	}

	args := BuildArgs(settings, server, "")

	if !reflect.DeepEqual(args, expected) {
		t.Errorf("xfreerdp3 NLA args mismatch.\nExpected: %v\nGot:      %v", expected, args)
	}
}

func TestBuildArgsWithXFreeRDP3NLADisabled(t *testing.T) {
	settings := config.SettingsConfig{
		Engine:    "xfreerdp3",
		Clipboard: true,
		SecNLA:    true,
	}

	server := config.ServerConfig{
		HostIP:    "192.168.1.50",
		Name:      "No NLA Server xfreerdp3",
		Clipboard: "default",
		SecNLA:    "no",
	}

	expected := []string{
		"/v:192.168.1.50",
		"/t:No NLA Server xfreerdp3 - GoXFreeRDP",
		"/u:",
		"/p:",
		"+clipboard",
		"/sec:nla:off",
		"-fonts",
		"-wallpaper",
		"-themes",
	}

	args := BuildArgs(settings, server, "")

	if !reflect.DeepEqual(args, expected) {
		t.Errorf("xfreerdp3 disabled NLA args mismatch.\nExpected: %v\nGot:      %v", expected, args)
	}
}

func TestBuildArgsWithLogLevel(t *testing.T) {
	settings := config.SettingsConfig{
		Engine:    "xfreerdp",
		Clipboard: true,
		SecNLA:    true,
		LogLevel:  "WARN",
	}

	server := config.ServerConfig{
		HostIP:    "192.168.1.50",
		Name:      "Log Level Server",
		Clipboard: "default",
		SecNLA:    "default",
	}

	expected := []string{
		"/log-level:WARN",
		"/v:192.168.1.50",
		"/t:Log Level Server - GoXFreeRDP",
		"+clipboard",
		"+sec-nla",
		"-fonts",
		"-wallpaper",
		"-themes",
	}

	args := BuildArgs(settings, server, "")

	if !reflect.DeepEqual(args, expected) {
		t.Errorf("Log level args mismatch.\nExpected: %v\nGot:      %v", expected, args)
	}
}
