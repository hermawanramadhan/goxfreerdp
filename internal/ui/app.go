package ui

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"

	"goxfreerdp/internal/config"
	"goxfreerdp/internal/rdp"
)

type AppUI struct {
	Config        *config.AppConfig
	LogTextBuffer *gtk.TextBuffer
	LogTextView   *gtk.TextView
	Notebook      *gtk.Notebook
	Window        *gtk.Window
	ListBox       *gtk.ListBox
}

// GtkLogWriter implements io.Writer and logs output directly into a scrollable GtkTextView.
type GtkLogWriter struct {
	textBuffer *gtk.TextBuffer
	textView   *gtk.TextView
}

func (w *GtkLogWriter) Write(p []byte) (n int, err error) {
	text := string(p)
	os.Stdout.Write(p)

	glib.IdleAdd(func() {
		endIter := w.textBuffer.GetEndIter()
		w.textBuffer.Insert(endIter, text)

		endIter = w.textBuffer.GetEndIter()
		w.textBuffer.PlaceCursor(endIter)
		w.textView.ScrollToMark(w.textBuffer.GetInsert(), 0.0, true, 0.0, 1.0)
	})

	return len(p), nil
}

// PopulateServerList populates the GtkListBox with servers from the configuration.
func (app *AppUI) PopulateServerList() error {
	children := app.ListBox.GetChildren()
	for l := children; l != nil; l = l.Next() {
		if widget, ok := l.Data().(gtk.IWidget); ok {
			app.ListBox.Remove(widget)
		}
	}

	for _, s := range app.Config.Servers {
		row, err := gtk.ListBoxRowNew()
		if err != nil {
			return fmt.Errorf("failed to create ListBoxRow: %w", err)
		}

		box, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 12)
		if err != nil {
			return fmt.Errorf("failed to create Box: %w", err)
		}

		box.SetMarginStart(16)
		box.SetMarginEnd(16)
		box.SetMarginTop(12)
		box.SetMarginBottom(12)

		// 1. Server Icon
		serverImg, err := gtk.ImageNewFromIconName("network-server-symbolic", gtk.ICON_SIZE_DND)
		if err == nil {
			serverImg.SetMarginEnd(8)
			box.PackStart(serverImg, false, false, 0)
		}

		// 2. Text layout Box (Vertical)
		labelBox, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 4)
		if err != nil {
			return fmt.Errorf("failed to create label Box: %w", err)
		}

		lblName, err := gtk.LabelNew(s.Name)
		if err != nil {
			return fmt.Errorf("failed to create Name Label: %w", err)
		}
		lblName.SetHAlign(gtk.ALIGN_START)
		if styleCtx, err := lblName.GetStyleContext(); err == nil {
			styleCtx.AddClass("server-title")
		}
		labelBox.PackStart(lblName, false, false, 0)

		// Details row (Horizontal)
		detailsBox, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 6)
		if err != nil {
			return fmt.Errorf("failed to create details Box: %w", err)
		}

		lblHost, err := gtk.LabelNew(s.HostIP)
		if err != nil {
			return fmt.Errorf("failed to create Host Label: %w", err)
		}
		lblHost.SetHAlign(gtk.ALIGN_START)
		if styleCtx, err := lblHost.GetStyleContext(); err == nil {
			styleCtx.AddClass("server-subtitle")
		}
		detailsBox.PackStart(lblHost, false, false, 0)

		if s.Username != "" {
			badgeUser, err := gtk.LabelNew("👤 " + s.Username)
			if err == nil {
				badgeUser.SetHAlign(gtk.ALIGN_START)
				if styleCtx, err := badgeUser.GetStyleContext(); err == nil {
					styleCtx.AddClass("badge")
				}
				detailsBox.PackStart(badgeUser, false, false, 0)
			}
		}

		engineName := s.Engine
		if engineName == "" || engineName == "default" {
			engineName = app.Config.Settings.Engine
		}
		if engineName == "" {
			engineName = "xfreerdp"
		}
		badgeEngine, err := gtk.LabelNew("⚙️ " + engineName)
		if err == nil {
			badgeEngine.SetHAlign(gtk.ALIGN_START)
			if styleCtx, err := badgeEngine.GetStyleContext(); err == nil {
				styleCtx.AddClass("badge")
			}
			detailsBox.PackStart(badgeEngine, false, false, 0)
		}

		labelBox.PackStart(detailsBox, false, false, 0)

		// 3. Connect Button
		btnConnect, err := gtk.ButtonNew()
		if err != nil {
			return fmt.Errorf("failed to create Connect Button: %w", err)
		}
		btnConnect.SetHAlign(gtk.ALIGN_END)
		btnConnect.SetTooltipText("Connect to " + s.Name)

		img, err := gtk.ImageNewFromIconName("media-playback-start", gtk.ICON_SIZE_BUTTON)
		if err == nil {
			btnConnect.SetImage(img)
			btnConnect.SetAlwaysShowImage(true)
		}

		if styleCtxConnect, err := btnConnect.GetStyleContext(); err == nil {
			styleCtxConnect.AddClass("btn-connect")
		}

		box.PackStart(labelBox, true, true, 0)
		box.PackEnd(btnConnect, false, false, 0)

		eventBox, err := gtk.EventBoxNew()
		if err != nil {
			return fmt.Errorf("failed to create EventBox: %w", err)
		}

		eventBox.Add(box)
		row.Add(eventBox)
		app.ListBox.Add(row)

		server := s
		btnConnect.Connect("clicked", func() {
			if strings.TrimSpace(server.HostIP) == "" {
				showErrorDialog(app.Window, "Host / IP is required to connect.")
				return
			}
			fmt.Printf("[GUI] Starting RDP connection to %s (%s)...\n", server.Name, server.HostIP)
			app.RunConnectionWithAuthFallback(server)
		})

		btnConnect.Connect("button-press-event", func(btn *gtk.Button, event *gdk.Event) bool {
			eventButton := gdk.EventButtonNewFromEvent(event)
			if eventButton.Button() == 3 {
				return true
			}
			return false
		})

		eventBox.Connect("button-press-event", func(eb *gtk.EventBox, event *gdk.Event) bool {
			eventButton := gdk.EventButtonNewFromEvent(event)
			if eventButton.Button() == 1 {
				if eventButton.Type() == gdk.EVENT_2BUTTON_PRESS {
					if strings.TrimSpace(server.HostIP) == "" {
						showErrorDialog(app.Window, "Host / IP is required to connect.")
						return true
					}
					fmt.Printf("[GUI] Starting RDP connection to %s (%s)...\n", server.Name, server.HostIP)
					app.RunConnectionWithAuthFallback(server)
					return true
				}
			} else if eventButton.Button() == 3 {
				showContextMenu(app, server.ID, eventButton)
				return true
			}
			return false
		})
	}

	app.ListBox.ShowAll()
	return nil
}

// BuildMainWindow constructs the main GTK application window.
func (app *AppUI) BuildMainWindow() error {
	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		return fmt.Errorf("failed to create main Window: %w", err)
	}
	win.SetDefaultSize(600, 485)
	win.SetPosition(gtk.WIN_POS_CENTER)
	win.SetTitle("GoXFreeRDP")
	app.Window = win

	notebook, err := gtk.NotebookNew()
	if err != nil {
		return fmt.Errorf("failed to create Notebook: %w", err)
	}
	win.Add(notebook)
	app.Notebook = notebook

	aboutBtn, err := gtk.ButtonNew()
	if err == nil {
		img, err := gtk.ImageNewFromIconName("help-about-symbolic", gtk.ICON_SIZE_MENU)
		if err == nil {
			aboutBtn.SetImage(img)
			aboutBtn.SetAlwaysShowImage(true)
		}
		aboutBtn.SetRelief(gtk.RELIEF_NONE)
		aboutBtn.SetTooltipText("About GoXFreeRDP")
		aboutBtn.Connect("clicked", func() {
			showAboutDialog(app.Window)
		})
		notebook.SetActionWidget(aboutBtn, gtk.PACK_END)
		aboutBtn.ShowAll()
	}

	serversBox, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 10)
	if err != nil {
		return err
	}

	scroll, err := gtk.ScrolledWindowNew(nil, nil)
	if err != nil {
		return fmt.Errorf("failed to create ScrolledWindow: %w", err)
	}
	scroll.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)
	scroll.SetMarginStart(12)
	scroll.SetMarginEnd(12)
	scroll.SetMarginTop(12)
	scroll.SetMarginBottom(0)

	listBox, err := gtk.ListBoxNew()
	if err != nil {
		return fmt.Errorf("failed to create ListBox: %w", err)
	}
	listBox.SetSelectionMode(gtk.SELECTION_NONE)
	app.ListBox = listBox

	scroll.Add(listBox)
	serversBox.PackStart(scroll, true, true, 0)

	actionBar, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 10)
	if err != nil {
		return fmt.Errorf("failed to create Action Box: %w", err)
	}
	actionBar.SetMarginStart(12)
	actionBar.SetMarginEnd(12)
	actionBar.SetMarginTop(10)
	actionBar.SetMarginBottom(12)

	addBtn, err := gtk.ButtonNew()
	if err != nil {
		return fmt.Errorf("failed to create Add Server button: %w", err)
	}
	addBtn.SetLabel("Add Server")
	if img, err := gtk.ImageNewFromIconName("list-add-symbolic", gtk.ICON_SIZE_BUTTON); err == nil {
		addBtn.SetImage(img)
		addBtn.SetAlwaysShowImage(true)
	}

	openRdpBtn, err := gtk.ButtonNew()
	if err != nil {
		return fmt.Errorf("failed to create Open .rdp button: %w", err)
	}
	openRdpBtn.SetLabel("Open .rdp")
	if img, err := gtk.ImageNewFromIconName("document-open-symbolic", gtk.ICON_SIZE_BUTTON); err == nil {
		openRdpBtn.SetImage(img)
		openRdpBtn.SetAlwaysShowImage(true)
	}

	settingsBtn, err := gtk.ButtonNew()
	if err != nil {
		return fmt.Errorf("failed to create Settings button: %w", err)
	}
	settingsBtn.SetLabel("Settings")
	if img, err := gtk.ImageNewFromIconName("preferences-system-symbolic", gtk.ICON_SIZE_BUTTON); err == nil {
		settingsBtn.SetImage(img)
		settingsBtn.SetAlwaysShowImage(true)
	}

	actionBar.PackStart(addBtn, false, false, 0)
	actionBar.PackEnd(settingsBtn, false, false, 0)
	actionBar.PackEnd(openRdpBtn, false, false, 0)
	serversBox.PackEnd(actionBar, false, false, 0)

	logScroll, err := gtk.ScrolledWindowNew(nil, nil)
	if err != nil {
		return err
	}
	logScroll.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)

	logTextView, err := gtk.TextViewNew()
	if err != nil {
		return err
	}
	logTextView.SetEditable(false)
	logTextView.SetCursorVisible(false)
	logTextView.SetWrapMode(gtk.WRAP_WORD_CHAR)

	styleCtx, err := logTextView.GetStyleContext()
	if err == nil {
		styleCtx.AddClass("monospace")
	}

	logTextBuffer, err := logTextView.GetBuffer()
	if err != nil {
		return err
	}
	logScroll.Add(logTextView)

	app.LogTextBuffer = logTextBuffer
	app.LogTextView = logTextView

	lblTabServers, _ := gtk.LabelNew("Servers")
	notebook.AppendPage(serversBox, lblTabServers)

	lblTabLogs, _ := gtk.LabelNew("Connection Logs")
	notebook.AppendPage(logScroll, lblTabLogs)

	settingsBtn.Connect("clicked", func() {
		showSettingsDialog(app)
	})

	// About button click is handled inline during creation

	addBtn.Connect("clicked", func() {
		showAddServerDialog(app)
	})

	openRdpBtn.Connect("clicked", func() {
		fileChooser, err := gtk.FileChooserDialogNewWith2Buttons(
			"Open RDP File",
			win,
			gtk.FILE_CHOOSER_ACTION_OPEN,
			"Cancel",
			gtk.RESPONSE_CANCEL,
			"Open",
			gtk.RESPONSE_ACCEPT,
		)
		if err != nil {
			fmt.Printf("Failed to create FileChooserDialog: %v\n", err)
			return
		}
		defer fileChooser.Destroy()

		filter, err := gtk.FileFilterNew()
		if err == nil {
			filter.SetName("RDP Files (*.rdp)")
			filter.AddPattern("*.rdp")
			fileChooser.AddFilter(filter)
		}

		response := fileChooser.Run()
		if response == gtk.RESPONSE_ACCEPT {
			filePath := fileChooser.GetFilename()
			var dummyServer config.ServerConfig
			dummyServer.HostIP = filePath
			dummyServer.Name = filepath.Base(filePath)

			fmt.Printf("[GUI] Starting RDP connection to RDP file %s...\n", filePath)
			app.RunConnectionWithAuthFallback(dummyServer)
		}
	})

	err = app.PopulateServerList()
	if err != nil {
		return fmt.Errorf("failed to load server list to UI: %w", err)
	}

	win.Connect("destroy", func() {
		_ = os.Remove(GetSocketPath())
		gtk.MainQuit()
	})

	return nil
}

// RunConnectionWithAuthFallback runs the connection, intercepting authentication failures to prompt for passwords.
func (app *AppUI) RunConnectionWithAuthFallback(server config.ServerConfig) {
	var logWriter io.Writer
	if app.LogTextBuffer != nil && app.LogTextView != nil {
		glib.IdleAdd(func() {
			endIter := app.LogTextBuffer.GetEndIter()
			separator := fmt.Sprintf("\n\n--------------------------------------------------\n--- Session: %s ---\n--------------------------------------------------\n", time.Now().Format("2006-01-02 15:04:05"))
			app.LogTextBuffer.Insert(endIter, separator)
		})

		logWriter = &GtkLogWriter{
			textBuffer: app.LogTextBuffer,
			textView:   app.LogTextView,
		}
	}

	if app.Notebook != nil {
		glib.IdleAdd(func() {
			app.Notebook.SetCurrentPage(1)
		})
	}

	if app.LogTextBuffer != nil {
		msg1 := fmt.Sprintf("[GUI] Starting RDP connection to %s (%s)...\n", server.Name, server.HostIP)
		args := rdp.BuildArgs(app.Config.Settings, server, "")
		
		engine := server.Engine
		if engine == "" || engine == "default" {
			engine = app.Config.Settings.Engine
		}
		if engine == "" {
			engine = "xfreerdp"
		}
		
		censoredArgs := rdp.CensorArgs(args)
		msg2 := fmt.Sprintf("[RDP Launch] Executing: %s %s\n\n", engine, strings.Join(censoredArgs, " "))

		glib.IdleAdd(func() {
			endIter := app.LogTextBuffer.GetEndIter()
			app.LogTextBuffer.Insert(endIter, msg1+msg2)
		})
	}

	errChan := rdp.LaunchRDP(app.Config.Settings, server, "", logWriter)
	go func() {
		err := <-errChan
		if err != nil {
			isAuthFail := false
			if rdpErr, ok := err.(*rdp.RDPError); ok {
				msg := strings.ToLower(rdpErr.Stderr)
				if strings.Contains(msg, "auth") ||
					strings.Contains(msg, "logon") ||
					strings.Contains(msg, "password") ||
					strings.Contains(msg, "credentials") ||
					strings.Contains(msg, "0xc000006d") ||
					strings.Contains(msg, "0xc0000022") ||
					strings.Contains(msg, "logon failure") {
					isAuthFail = true
				}
			}

			glib.IdleAdd(func() {
				if isAuthFail {
					password, ok := promptPasswordDialog(app.Window, server.Name)
					if ok {
						server.Password = password
						app.RunConnectionWithAuthFallback(server)
					}
				} else {
					showErrorDialog(app.Window, fmt.Sprintf("Failed to launch RDP connection: %v", err))
				}
			})
		} else {
			if app.LogTextBuffer != nil {
				glib.IdleAdd(func() {
					endIter := app.LogTextBuffer.GetEndIter()
					app.LogTextBuffer.Insert(endIter, "\n[GoXFreeRDP] Connection closed.\n")
				})
			}
		}
	}()
}

func SetupCSS() {
	cssProvider, err := gtk.CssProviderNew()
	if err == nil {
		css := `
			window {
				font-family: "Inter", "Segoe UI", "Cantarell", sans-serif;
			}
			notebook {
				background-color: @theme_bg_color;
				border: none;
			}
			notebook header tab {
				padding: 10px 20px;
				background-color: transparent;
				color: mix(@theme_fg_color, @theme_bg_color, 0.4);
				font-weight: 600;
				border: none;
				border-bottom: 2px solid transparent;
				transition: all 0.2s ease;
			}
			notebook header tab:hover {
				color: mix(@theme_fg_color, @theme_bg_color, 0.7);
			}
			notebook header tab:checked {
				background-color: transparent;
				color: @theme_selected_bg_color;
				border-bottom: 2px solid mix(@theme_selected_bg_color, @theme_bg_color, 0.65);
			}
			eventbox {
				background-color: transparent !important;
			}
			list {
				background-color: transparent;
				border: none;
			}
			row {
				padding: 10px 14px;
				margin-bottom: 8px;
				border: 1px solid mix(@theme_fg_color, @theme_bg_color, 0.08) !important;
				border-radius: 8px !important;
				background-color: @theme_base_color !important;
				transition: all 0.2s ease;
			}
			row:last-child {
				margin-bottom: 0;
			}
			row:hover {
				background-color: mix(@theme_selected_bg_color, @theme_base_color, 0.06) !important;
				border-color: mix(@theme_selected_bg_color, @theme_base_color, 0.25) !important;
			}
			row:selected {
				background-color: mix(@theme_selected_bg_color, @theme_base_color, 0.12) !important;
				color: @theme_fg_color;
			}
			row:selected:hover {
				background-color: mix(@theme_selected_bg_color, @theme_base_color, 0.18) !important;
				color: @theme_fg_color;
			}
			button {
				border-radius: 6px;
				padding: 8px 16px;
				font-weight: 500;
			}
			.btn-connect {
				background-image: linear-gradient(to bottom, #2ecc71, #27ae60) !important;
				color: #ffffff !important;
				border: 1px solid #219653 !important;
				box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.25), 0 1px 2px rgba(0, 0, 0, 0.2) !important;
				border-radius: 9999px !important;
				padding: 0 !important;
				min-width: 38px !important;
				min-height: 38px !important;
				width: 38px !important;
				height: 38px !important;
				transition: all 0.2s ease;
			}
			.btn-connect:hover {
				background-image: linear-gradient(to bottom, #2ebd70, #219653) !important;
				border-color: #1e8449 !important;
			}
			.btn-connect:active {
				background-image: linear-gradient(to bottom, #219653, #27ae60) !important;
				box-shadow: inset 0 2px 4px rgba(0, 0, 0, 0.25) !important;
			}
			textview text {
				background-color: @theme_base_color;
				color: @theme_text_color;
				font-family: "JetBrains Mono", "Fira Code", "Liberation Mono", monospace;
				font-size: 10pt;
				padding: 12px;
			}
			entry {
				border-radius: 6px;
				padding: 8px;
			}
			.server-title {
				font-weight: bold;
				font-size: 11pt;
			}
			.server-subtitle {
				font-size: 9pt;
				opacity: 0.65;
			}
			.badge {
				background-color: mix(@theme_fg_color, @theme_bg_color, 0.06);
				color: mix(@theme_fg_color, @theme_bg_color, 0.65);
				border-radius: 6px;
				padding: 2px 8px;
				font-size: 8.5pt;
				font-weight: 600;
			}
		`
		cssProvider.LoadFromData(css)
		screen, err := gdk.ScreenGetDefault()
		if err == nil {
			gtk.AddProviderForScreen(screen, cssProvider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)
		}
	}
}

// isSystemDark checks if the system theme preference is dark.
func isSystemDark() bool {
	// First, check the standard Freedesktop appearance portal
	cmd := exec.Command("gdbus", "call", "--session",
		"--dest", "org.freedesktop.portal.Desktop",
		"--object-path", "/org/freedesktop/portal/desktop",
		"--method", "org.freedesktop.portal.Settings.Read",
		"org.freedesktop.appearance", "color-scheme")
	out, err := cmd.Output()
	if err == nil {
		if strings.Contains(string(out), "uint32 1") {
			return true
		}
		if strings.Contains(string(out), "uint32 2") {
			return false
		}
	}

	// Fallback to gsettings color-scheme
	cmd = exec.Command("gsettings", "get", "org.gnome.desktop.interface", "color-scheme")
	out, err = cmd.Output()
	if err == nil {
		if strings.Contains(strings.ToLower(string(out)), "dark") {
			return true
		}
	}

	// Fallback to checking the current theme name
	cmd = exec.Command("gsettings", "get", "org.gnome.desktop.interface", "gtk-theme")
	out, err = cmd.Output()
	if err == nil {
		themeName := strings.ToLower(string(out))
		if strings.Contains(themeName, "dark") {
			return true
		}
	}

	return false
}

// applySystemTheme applies the detected dark mode preference and updates style classes.
func applySystemTheme(app *AppUI) {
	settings, err := gtk.SettingsGetDefault()
	if err != nil {
		return
	}
	isDark := isSystemDark()
	settings.Set("gtk-application-prefer-dark-theme", isDark)

	if app != nil && app.Window != nil {
		styleCtx, err := app.Window.GetStyleContext()
		if err == nil {
			if isDark {
				styleCtx.AddClass("dark-mode")
				styleCtx.RemoveClass("light-mode")
			} else {
				styleCtx.AddClass("light-mode")
				styleCtx.RemoveClass("dark-mode")
			}
		}
	}
}

// StartApp starts the GTK application
func StartApp(cfg *config.AppConfig, initialRDPFile string) error {
	gtk.Init(nil)
	applySystemTheme(nil)
	SetupCSS()

	app := &AppUI{
		Config: cfg,
	}

	err := app.BuildMainWindow()
	if err != nil {
		return err
	}

	applySystemTheme(app)

	// Periodic check for system theme changes (every 2 seconds)
	go func() {
		lastDark := isSystemDark()
		for {
			time.Sleep(2 * time.Second)
			currentDark := isSystemDark()
			if currentDark != lastDark {
				lastDark = currentDark
				glib.IdleAdd(func() {
					applySystemTheme(app)
				})
			}
		}
	}()

	app.Window.ShowAll()

	err = app.StartIPCServer()
	if err != nil {
		fmt.Printf("[IPC] Failed to start IPC server: %v\n", err)
	}

	if initialRDPFile != "" {
		glib.IdleAdd(func() {
			var dummyServer config.ServerConfig
			dummyServer.HostIP = initialRDPFile
			dummyServer.Name = filepath.Base(initialRDPFile)

			if app.LogTextBuffer != nil {
				endIter := app.LogTextBuffer.GetEndIter()
				msg := fmt.Sprintf("[CLI Mode] Opening RDP connection from file: %s\n", initialRDPFile)
				app.LogTextBuffer.Insert(endIter, msg)
			}

			app.RunConnectionWithAuthFallback(dummyServer)
		})
	}

	gtk.Main()
	return nil
}
