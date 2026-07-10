package ui

import (
	"fmt"

	"github.com/gotk3/gotk3/gtk"

	"goxfreerdp/internal/config"
)

// showSettingsDialog opens the settings configuration dialog.
func showSettingsDialog(app *AppUI) {
	dialog, err := gtk.DialogNew()
	if err != nil {
		fmt.Printf("Failed to create Settings dialog: %v\n", err)
		return
	}
	defer dialog.Destroy()

	dialog.SetTitle("Settings")
	dialog.SetTransientFor(app.Window)
	dialog.SetModal(true)
	dialog.SetDefaultSize(550, 600)

	dialog.AddButton("Cancel", gtk.RESPONSE_CANCEL)
	dialog.AddButton("Save", gtk.RESPONSE_OK)

	contentArea, err := dialog.GetContentArea()
	if err != nil {
		fmt.Printf("Failed to get dialog content area: %v\n", err)
		return
	}

	// Create a scrolled window inside content area
	scroll, err := gtk.ScrolledWindowNew(nil, nil)
	if err != nil {
		return
	}
	scroll.SetPolicy(gtk.POLICY_NEVER, gtk.POLICY_AUTOMATIC)

	grid, err := gtk.GridNew()
	if err != nil {
		return
	}
	grid.SetColumnSpacing(15)
	grid.SetRowSpacing(10)
	grid.SetMarginStart(15)
	grid.SetMarginEnd(15)
	grid.SetMarginTop(15)
	grid.SetMarginBottom(15)
	scroll.Add(grid)

	r := 0

	// 1. Default Connection Options Header
	headerConnection, _ := createHeaderLabel("Default Connection Options")
	grid.Attach(headerConnection, 0, r, 2, 1)
	r++

	lblEngine, _ := gtk.LabelNew("RDP Engine:")
	lblEngine.SetHAlign(gtk.ALIGN_START)
	comboEngine, _ := gtk.ComboBoxTextNew()
	comboEngine.Append("xfreerdp", "xfreerdp (v2)")
	comboEngine.Append("xfreerdp3", "xfreerdp3 (v3)")
	comboEngine.SetHExpand(true)
	grid.Attach(lblEngine, 0, r, 1, 1)
	grid.Attach(comboEngine, 1, r, 1, 1)
	r++

	chkIgnoreCert, _ := gtk.CheckButtonNewWithLabel("Ignore Certificate errors")
	chkIgnoreCert.SetHAlign(gtk.ALIGN_START)
	grid.Attach(chkIgnoreCert, 0, r, 2, 1)
	r++

	lblTLS, _ := gtk.LabelNew("TLS Security Level:")
	lblTLS.SetHAlign(gtk.ALIGN_START)
	comboTLS, _ := gtk.ComboBoxTextNew()
	comboTLS.SetHExpand(true)
	levels := []string{"default", "0", "1", "2", "3", "4", "5"}
	levelLabels := []string{
		"Default (1)",
		"0 (Minimum security / SSL CTX level 0)",
		"1 (Default OpenSSL level 1)",
		"2 (OpenSSL level 2 - 112 bit security)",
		"3 (OpenSSL level 3 - 128 bit security)",
		"4 (OpenSSL level 4 - 192 bit security)",
		"5 (OpenSSL level 5 - 256 bit security)",
	}
	for i, lbl := range levelLabels {
		comboTLS.Append(levels[i], lbl)
	}
	grid.Attach(lblTLS, 0, r, 1, 1)
	grid.Attach(comboTLS, 1, r, 1, 1)
	r++

	lblPort, _ := gtk.LabelNew("Default Port:")
	lblPort.SetHAlign(gtk.ALIGN_START)
	entryPort, _ := gtk.EntryNew()
	entryPort.SetPlaceholderText("e.g. 3389")
	entryPort.SetHExpand(true)
	grid.Attach(lblPort, 0, r, 1, 1)
	grid.Attach(entryPort, 1, r, 1, 1)
	r++

	// 2. Redirections Header
	headerRedir, _ := createHeaderLabel("Redirections & Integration")
	grid.Attach(headerRedir, 0, r, 2, 1)
	r++

	chkClipboard, _ := gtk.CheckButtonNewWithLabel("Enable clipboard redirection")
	chkClipboard.SetHAlign(gtk.ALIGN_START)
	grid.Attach(chkClipboard, 0, r, 2, 1)
	r++

	chkSecNLA, _ := gtk.CheckButtonNewWithLabel("Enable Network Level Authentication (NLA)")
	chkSecNLA.SetHAlign(gtk.ALIGN_START)
	grid.Attach(chkSecNLA, 0, r, 2, 1)
	r++

	chkSound, _ := gtk.CheckButtonNewWithLabel("Redirect audio playback")
	chkSound.SetHAlign(gtk.ALIGN_START)
	grid.Attach(chkSound, 0, r, 2, 1)
	r++

	chkDrive, _ := gtk.CheckButtonNewWithLabel("Share local home directory")
	chkDrive.SetHAlign(gtk.ALIGN_START)
	grid.Attach(chkDrive, 0, r, 2, 1)
	r++

	// 3. Display Settings Header
	headerDisplay, _ := createHeaderLabel("Display Settings")
	grid.Attach(headerDisplay, 0, r, 2, 1)
	r++

	chkFullscreen, _ := gtk.CheckButtonNewWithLabel("Start in fullscreen")
	chkFullscreen.SetHAlign(gtk.ALIGN_START)
	grid.Attach(chkFullscreen, 0, r, 2, 1)
	r++

	chkDynamicRes, _ := gtk.CheckButtonNewWithLabel("Dynamic resolution adjustment")
	chkDynamicRes.SetHAlign(gtk.ALIGN_START)
	grid.Attach(chkDynamicRes, 0, r, 2, 1)
	r++

	chkMultimon, _ := gtk.CheckButtonNewWithLabel("Use multiple monitors")
	chkMultimon.SetHAlign(gtk.ALIGN_START)
	grid.Attach(chkMultimon, 0, r, 2, 1)
	r++

	// 4. Performance & Experience Header
	headerPerf, _ := createHeaderLabel("Performance & Experience")
	grid.Attach(headerPerf, 0, r, 2, 1)
	r++

	chkFonts, _ := gtk.CheckButtonNewWithLabel("Enable font smoothing")
	chkFonts.SetHAlign(gtk.ALIGN_START)
	grid.Attach(chkFonts, 0, r, 2, 1)
	r++

	chkWallpaper, _ := gtk.CheckButtonNewWithLabel("Enable desktop wallpaper")
	chkWallpaper.SetHAlign(gtk.ALIGN_START)
	grid.Attach(chkWallpaper, 0, r, 2, 1)
	r++

	chkThemes, _ := gtk.CheckButtonNewWithLabel("Enable desktop themes")
	chkThemes.SetHAlign(gtk.ALIGN_START)
	grid.Attach(chkThemes, 0, r, 2, 1)
	r++

	// 5. Advanced Settings Header
	headerAdvanced, _ := createHeaderLabel("Advanced Settings")
	grid.Attach(headerAdvanced, 0, r, 2, 1)
	r++

	lblCustom, _ := gtk.LabelNew("Custom Parameters:")
	lblCustom.SetHAlign(gtk.ALIGN_START)
	entryCustom, _ := gtk.EntryNew()
	entryCustom.SetPlaceholderText("e.g. /gdi:hw /network:lan")
	entryCustom.SetHExpand(true)
	grid.Attach(lblCustom, 0, r, 1, 1)
	grid.Attach(entryCustom, 1, r, 1, 1)
	r++

	lblLogLevel, _ := gtk.LabelNew("Log Level:")
	lblLogLevel.SetHAlign(gtk.ALIGN_START)
	comboLogLevel, _ := gtk.ComboBoxTextNew()
	comboLogLevel.SetHExpand(true)
	logLevels := []string{"default", "OFF", "FATAL", "ERROR", "WARN", "INFO", "DEBUG", "TRACE"}
	logLevelLabels := []string{
		"Default",
		"OFF (No logs)",
		"FATAL (Fatal errors only)",
		"ERROR (Errors only)",
		"WARN (Warnings and errors)",
		"INFO (Information, warnings, errors)",
		"DEBUG (Detailed debug logs)",
		"TRACE (Maximum diagnostic logs)",
	}
	for i, lbl := range logLevelLabels {
		comboLogLevel.Append(logLevels[i], lbl)
	}
	grid.Attach(lblLogLevel, 0, r, 1, 1)
	grid.Attach(comboLogLevel, 1, r, 1, 1)
	r++

	// Populate values from config
	chkIgnoreCert.SetActive(app.Config.Settings.IgnoreCertificate)
	chkClipboard.SetActive(app.Config.Settings.Clipboard)
	chkSecNLA.SetActive(app.Config.Settings.SecNLA)
	chkSound.SetActive(app.Config.Settings.Sound)
	chkDrive.SetActive(app.Config.Settings.ShareHome)
	chkFullscreen.SetActive(app.Config.Settings.Fullscreen)
	chkDynamicRes.SetActive(app.Config.Settings.DynamicRes)
	chkMultimon.SetActive(app.Config.Settings.Multimon)
	chkFonts.SetActive(app.Config.Settings.FontSmoothing)
	chkWallpaper.SetActive(app.Config.Settings.Wallpaper)
	chkThemes.SetActive(app.Config.Settings.Themes)
	entryCustom.SetText(app.Config.Settings.CustomParams)
	entryPort.SetText(app.Config.Settings.Port)

	logLevelIdx := 0
	for i, val := range logLevels {
		if val == app.Config.Settings.LogLevel {
			logLevelIdx = i
			break
		}
	}
	comboLogLevel.SetActive(logLevelIdx)

	if app.Config.Settings.Engine == "xfreerdp3" {
		comboEngine.SetActive(1)
	} else {
		comboEngine.SetActive(0)
	}

	tlsIdx := 0
	for i, val := range levels {
		if val == app.Config.Settings.TLSSecLevel {
			tlsIdx = i
			break
		}
	}
	comboTLS.SetActive(tlsIdx)

	contentArea.PackStart(scroll, true, true, 0)
	dialog.ShowAll()

	response := dialog.Run()
	if response == gtk.RESPONSE_OK {
		if comboEngine.GetActive() == 1 {
			app.Config.Settings.Engine = "xfreerdp3"
		} else {
			app.Config.Settings.Engine = "xfreerdp"
		}

		app.Config.Settings.IgnoreCertificate = chkIgnoreCert.GetActive()

		activeIdx := comboTLS.GetActive()
		if activeIdx >= 0 && activeIdx < len(levels) {
			app.Config.Settings.TLSSecLevel = levels[activeIdx]
		} else {
			app.Config.Settings.TLSSecLevel = "default"
		}

		app.Config.Settings.Clipboard = chkClipboard.GetActive()
		app.Config.Settings.SecNLA = chkSecNLA.GetActive()
		app.Config.Settings.Sound = chkSound.GetActive()
		app.Config.Settings.ShareHome = chkDrive.GetActive()
		app.Config.Settings.Fullscreen = chkFullscreen.GetActive()
		app.Config.Settings.DynamicRes = chkDynamicRes.GetActive()
		app.Config.Settings.Multimon = chkMultimon.GetActive()
		app.Config.Settings.FontSmoothing = chkFonts.GetActive()
		app.Config.Settings.Wallpaper = chkWallpaper.GetActive()
		app.Config.Settings.Themes = chkThemes.GetActive()
		app.Config.Settings.CustomParams, _ = entryCustom.GetText()
		app.Config.Settings.Port, _ = entryPort.GetText()

		activeLogIdx := comboLogLevel.GetActive()
		if activeLogIdx >= 0 && activeLogIdx < len(logLevels) {
			app.Config.Settings.LogLevel = logLevels[activeLogIdx]
		} else {
			app.Config.Settings.LogLevel = "default"
		}

		// Save modifications to config.json
		err = config.SaveConfig(*app.Config)
		if err != nil {
			fmt.Printf("[Config] Failed to save settings: %v\n", err)
		} else {
			fmt.Println("[Config] Settings updated successfully.")
		}
	}
}
