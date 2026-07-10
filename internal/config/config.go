package config

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// SettingsConfig holds default connection parameters and engine settings.
type SettingsConfig struct {
	Engine            string `json:"engine"`
	Host              string `json:"host"`
	Port              string `json:"port"`
	Username          string `json:"username"`
	Password          string `json:"password"`
	IgnoreCertificate bool   `json:"ignore_certificate"`
	TLSSecLevel       string `json:"tls_seclevel"`
	Clipboard         bool   `json:"clipboard"`
	SecNLA            bool   `json:"sec_nla"`
	Fullscreen        bool   `json:"fullscreen"`
	DynamicRes        bool   `json:"dynamic_res"`
	Multimon          bool   `json:"multimon"`
	Sound             bool   `json:"sound"`
	ShareHome         bool   `json:"share_home"`
	FontSmoothing     bool   `json:"font_smoothing"`
	Wallpaper         bool   `json:"wallpaper"`
	Themes            bool   `json:"themes"`
	CustomParams      string `json:"custom_params"`
	LogLevel          string `json:"log_level"`
}

// ServerConfig holds information for a single RDP server connection.
type ServerConfig struct {
	ID                string `json:"id"`
	Name              string `json:"name"`
	HostIP            string `json:"host_ip"`
	Port              string `json:"port"`
	Username          string `json:"username"`
	Password          string `json:"password"`
	Engine            string `json:"engine"`
	IgnoreCertificate string `json:"ignore_certificate"`
	TLSSecLevel       string `json:"tls_seclevel"`
	Clipboard         string `json:"clipboard"`
	SecNLA            string `json:"sec_nla"`
	Fullscreen        string `json:"fullscreen"`
	DynamicRes        string `json:"dynamic_res"`
	Multimon          string `json:"multimon"`
	Sound             string `json:"sound"`
	ShareHome         string `json:"share_home"`
	FontSmoothing     string `json:"font_smoothing"`
	Wallpaper         string `json:"wallpaper"`
	Themes            string `json:"themes"`
	CustomParams      string `json:"custom_params"`
}

// AppConfig represents the main application configuration structure.
type AppConfig struct {
	Settings SettingsConfig `json:"settings"`
	Servers  []ServerConfig `json:"servers"`
}

// EnsureConfigDir ensures the config directory exists.
func EnsureConfigDir() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user config directory: %w", err)
	}

	dirPath := filepath.Join(configDir, "goxfreerdp")
	err = os.MkdirAll(dirPath, 0755)
	if err != nil {
		return "", fmt.Errorf("failed to create config directory: %w", err)
	}

	return dirPath, nil
}

// GetConfigFilePath returns the path to config.json.
func GetConfigFilePath() (string, error) {
	dirPath, err := EnsureConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dirPath, "config.json"), nil
}

// LoadConfig reads config.json and decodes its content.
func LoadConfig() (AppConfig, error) {
	var cfg AppConfig
	// Set default values in case file doesn't exist
	cfg.Settings.Engine = "xfreerdp"
	cfg.Settings.IgnoreCertificate = true
	cfg.Settings.TLSSecLevel = "default"
	cfg.Settings.Clipboard = true
	cfg.Settings.DynamicRes = true
	cfg.Settings.FontSmoothing = true
	cfg.Settings.LogLevel = "default"
	cfg.Servers = make([]ServerConfig, 0)

	filePath, err := GetConfigFilePath()
	if err != nil {
		return cfg, err
	}

	file, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return cfg, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return cfg, fmt.Errorf("failed to decode JSON config: %w", err)
	}

	if cfg.Servers == nil {
		cfg.Servers = make([]ServerConfig, 0)
	}

	// Decode passwords from Base64
	for i, s := range cfg.Servers {
		if s.Password != "" {
			decoded, err := base64.StdEncoding.DecodeString(s.Password)
			if err == nil {
				cfg.Servers[i].Password = string(decoded)
			}
		}
	}

	return cfg, nil
}

// SaveConfig saves the configuration structure.
func SaveConfig(cfg AppConfig) error {
	filePath, err := GetConfigFilePath()
	if err != nil {
		return err
	}

	if cfg.Servers == nil {
		cfg.Servers = make([]ServerConfig, 0)
	}

	// Encode passwords to Base64
	cfgCopy := cfg
	cfgCopy.Servers = make([]ServerConfig, len(cfg.Servers))
	for i, s := range cfg.Servers {
		cfgCopy.Servers[i] = s
		if s.Password != "" {
			cfgCopy.Servers[i].Password = base64.StdEncoding.EncodeToString([]byte(s.Password))
		}
	}

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(cfgCopy); err != nil {
		return fmt.Errorf("failed to encode JSON config: %w", err)
	}

	return nil
}
