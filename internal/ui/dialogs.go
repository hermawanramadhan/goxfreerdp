package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	
	"goxfreerdp/internal/config"
)

func showInfoDialog(parent *gtk.Window, message string) {
	dialog := gtk.MessageDialogNew(
		parent,
		gtk.DIALOG_MODAL,
		gtk.MESSAGE_INFO,
		gtk.BUTTONS_OK,
		"%s",
		message,
	)
	defer dialog.Destroy()
	dialog.Run()
}

func showErrorDialog(parent *gtk.Window, message string) {
	dialog := gtk.MessageDialogNew(
		parent,
		gtk.DIALOG_MODAL,
		gtk.MESSAGE_ERROR,
		gtk.BUTTONS_OK,
		"%s",
		message,
	)
	defer dialog.Destroy()
	dialog.Run()
}

func confirmDeleteDialog(parentWindow *gtk.Window, serverName string) bool {
	dialog := gtk.MessageDialogNew(
		parentWindow,
		gtk.DIALOG_MODAL,
		gtk.MESSAGE_QUESTION,
		gtk.BUTTONS_YES_NO,
		"Are you sure you want to delete server '%s'?",
		serverName,
	)
	defer dialog.Destroy()

	response := dialog.Run()
	return response == gtk.RESPONSE_YES
}

func createHeaderLabel(text string) (*gtk.Label, error) {
	lbl, err := gtk.LabelNew("")
	if err != nil {
		return nil, err
	}
	lbl.SetMarkup(fmt.Sprintf("<b><span size='large'>%s</span></b>", text))
	lbl.SetHAlign(gtk.ALIGN_START)
	lbl.SetMarginTop(12)
	lbl.SetMarginBottom(6)
	return lbl, nil
}

func createThreeWayDropdown(defaultValue string) (*gtk.ComboBoxText, error) {
	combo, err := gtk.ComboBoxTextNew()
	if err != nil {
		return nil, err
	}
	combo.Append("default", "Default")
	combo.Append("yes", "Yes")
	combo.Append("no", "No")

	val := strings.ToLower(strings.TrimSpace(defaultValue))
	if val == "yes" || val == "true" || val == "1" || val == "enable" || val == "enabled" {
		combo.SetActive(1)
	} else if val == "no" || val == "false" || val == "0" || val == "disable" || val == "disabled" {
		combo.SetActive(2)
	} else {
		combo.SetActive(0)
	}
	return combo, nil
}

func getThreeWayDropdownValue(combo *gtk.ComboBoxText) string {
	active := combo.GetActive()
	if active == 1 {
		return "yes"
	}
	if active == 2 {
		return "no"
	}
	return "default"
}

func createEngineDropdown(defaultValue string) (*gtk.ComboBoxText, error) {
	combo, err := gtk.ComboBoxTextNew()
	if err != nil {
		return nil, err
	}
	combo.Append("default", "Default Settings")
	combo.Append("xfreerdp", "xfreerdp (v2)")
	combo.Append("xfreerdp3", "xfreerdp3 (v3)")

	val := strings.TrimSpace(defaultValue)
	if val == "xfreerdp" {
		combo.SetActive(1)
	} else if val == "xfreerdp3" {
		combo.SetActive(2)
	} else {
		combo.SetActive(0)
	}
	return combo, nil
}

func getEngineDropdownValue(combo *gtk.ComboBoxText) string {
	active := combo.GetActive()
	if active == 1 {
		return "xfreerdp"
	}
	if active == 2 {
		return "xfreerdp3"
	}
	return "default"
}

func createTLSDropdown(defaultValue string) (*gtk.ComboBoxText, error) {
	combo, err := gtk.ComboBoxTextNew()
	if err != nil {
		return nil, err
	}
	levels := []string{"default", "0", "1", "2", "3", "4", "5"}
	levelLabels := []string{
		"Default Settings",
		"0 (Minimum security / SSL CTX level 0)",
		"1 (Default OpenSSL level 1)",
		"2 (OpenSSL level 2 - 112 bit security)",
		"3 (OpenSSL level 3 - 128 bit security)",
		"4 (OpenSSL level 4 - 192 bit security)",
		"5 (OpenSSL level 5 - 256 bit security)",
	}
	for i, lbl := range levelLabels {
		combo.Append(levels[i], lbl)
	}

	val := strings.TrimSpace(defaultValue)
	activeIdx := 0
	for i, v := range levels {
		if v == val {
			activeIdx = i
			break
		}
	}
	combo.SetActive(activeIdx)
	return combo, nil
}

func getTLSDropdownValue(combo *gtk.ComboBoxText) string {
	levels := []string{"default", "0", "1", "2", "3", "4", "5"}
	active := combo.GetActive()
	if active >= 0 && active < len(levels) {
		return levels[active]
	}
	return "default"
}

type serverFormWidgets struct {
	EntryName       *gtk.Entry
	EntryHost       *gtk.Entry
	EntryPort       *gtk.Entry
	EntryUser       *gtk.Entry
	EntryPass       *gtk.Entry
	ComboEngine     *gtk.ComboBoxText
	ComboIgnoreCert *gtk.ComboBoxText
	ComboTLS        *gtk.ComboBoxText
	ComboClipboard  *gtk.ComboBoxText
	ComboSecNLA     *gtk.ComboBoxText
	ComboFullscreen *gtk.ComboBoxText
	ComboDynamicRes *gtk.ComboBoxText
	ComboMultimon   *gtk.ComboBoxText
	ComboSound      *gtk.ComboBoxText
	ComboShareHome  *gtk.ComboBoxText
	ComboFonts      *gtk.ComboBoxText
	ComboWallpaper  *gtk.ComboBoxText
	ComboThemes     *gtk.ComboBoxText
	EntryParams     *gtk.Entry
}

func (w *serverFormWidgets) ApplyToServer(server *config.ServerConfig) {
	server.Name, _ = w.EntryName.GetText()
	server.HostIP, _ = w.EntryHost.GetText()
	server.Port, _ = w.EntryPort.GetText()
	server.Username, _ = w.EntryUser.GetText()
	server.Password, _ = w.EntryPass.GetText()

	server.Engine = getEngineDropdownValue(w.ComboEngine)
	server.IgnoreCertificate = getThreeWayDropdownValue(w.ComboIgnoreCert)
	server.TLSSecLevel = getTLSDropdownValue(w.ComboTLS)
	server.Clipboard = getThreeWayDropdownValue(w.ComboClipboard)
	server.SecNLA = getThreeWayDropdownValue(w.ComboSecNLA)
	server.Fullscreen = getThreeWayDropdownValue(w.ComboFullscreen)
	server.DynamicRes = getThreeWayDropdownValue(w.ComboDynamicRes)
	server.Multimon = getThreeWayDropdownValue(w.ComboMultimon)
	server.Sound = getThreeWayDropdownValue(w.ComboSound)
	server.ShareHome = getThreeWayDropdownValue(w.ComboShareHome)
	server.FontSmoothing = getThreeWayDropdownValue(w.ComboFonts)
	server.Wallpaper = getThreeWayDropdownValue(w.ComboWallpaper)
	server.Themes = getThreeWayDropdownValue(w.ComboThemes)
	server.CustomParams, _ = w.EntryParams.GetText()
}

func buildServerFormGrid(server *config.ServerConfig) (*gtk.Grid, *serverFormWidgets, error) {
	if server == nil {
		server = &config.ServerConfig{
			Engine:            "default",
			IgnoreCertificate: "default",
			TLSSecLevel:       "default",
			Clipboard:         "default",
			SecNLA:            "default",
			Fullscreen:        "default",
			DynamicRes:        "default",
			Multimon:          "default",
			Sound:             "default",
			ShareHome:         "default",
			FontSmoothing:     "default",
			Wallpaper:         "default",
			Themes:            "default",
		}
	}

	grid, err := gtk.GridNew()
	if err != nil {
		return nil, nil, err
	}
	grid.SetColumnSpacing(15)
	grid.SetRowSpacing(10)
	grid.SetMarginStart(15)
	grid.SetMarginEnd(15)
	grid.SetMarginTop(15)
	grid.SetMarginBottom(15)

	w := &serverFormWidgets{}
	r := 0

	// General Connection Details
	headerGen, _ := createHeaderLabel("General Info")
	grid.Attach(headerGen, 0, r, 2, 1)
	r++

	lblName, _ := gtk.LabelNew("Server Name:")
	lblName.SetHAlign(gtk.ALIGN_START)
	w.EntryName, _ = gtk.EntryNew()
	w.EntryName.SetPlaceholderText("e.g. Main Office")
	w.EntryName.SetHExpand(true)
	w.EntryName.SetText(server.Name)
	grid.Attach(lblName, 0, r, 1, 1)
	grid.Attach(w.EntryName, 1, r, 1, 1)
	r++

	lblHost, _ := gtk.LabelNew("Host / IP:")
	lblHost.SetHAlign(gtk.ALIGN_START)
	w.EntryHost, _ = gtk.EntryNew()
	w.EntryHost.SetPlaceholderText("e.g. 192.168.1.100")
	w.EntryHost.SetHExpand(true)
	w.EntryHost.SetText(server.HostIP)
	grid.Attach(lblHost, 0, r, 1, 1)
	grid.Attach(w.EntryHost, 1, r, 1, 1)
	r++

	lblPort, _ := gtk.LabelNew("Port:")
	lblPort.SetHAlign(gtk.ALIGN_START)
	w.EntryPort, _ = gtk.EntryNew()
	w.EntryPort.SetPlaceholderText("default setting (3389)")
	w.EntryPort.SetHExpand(true)
	w.EntryPort.SetText(server.Port)
	grid.Attach(lblPort, 0, r, 1, 1)
	grid.Attach(w.EntryPort, 1, r, 1, 1)
	r++

	lblUser, _ := gtk.LabelNew("Username:")
	lblUser.SetHAlign(gtk.ALIGN_START)
	w.EntryUser, _ = gtk.EntryNew()
	w.EntryUser.SetPlaceholderText("e.g. administrator")
	w.EntryUser.SetHExpand(true)
	w.EntryUser.SetText(server.Username)
	grid.Attach(lblUser, 0, r, 1, 1)
	grid.Attach(w.EntryUser, 1, r, 1, 1)
	r++

	lblPass, _ := gtk.LabelNew("Password:")
	lblPass.SetHAlign(gtk.ALIGN_START)
	w.EntryPass, _ = gtk.EntryNew()
	w.EntryPass.SetVisibility(false)
	w.EntryPass.SetPlaceholderText("Optional password")
	w.EntryPass.SetHExpand(true)
	w.EntryPass.SetText(server.Password)
	grid.Attach(lblPass, 0, r, 1, 1)
	grid.Attach(w.EntryPass, 1, r, 1, 1)
	r++

	chkShowPass, _ := gtk.CheckButtonNewWithLabel("Show Password")
	chkShowPass.SetHAlign(gtk.ALIGN_START)
	grid.Attach(chkShowPass, 1, r, 1, 1)
	r++

	chkShowPass.Connect("toggled", func() {
		w.EntryPass.SetVisibility(chkShowPass.GetActive())
	})

	// Advanced Parameter Overrides Grid
	overridesGrid, err := gtk.GridNew()
	if err != nil {
		return nil, nil, err
	}
	overridesGrid.SetColumnSpacing(15)
	overridesGrid.SetRowSpacing(10)

	or := 0

	lblEngine, _ := gtk.LabelNew("RDP Engine:")
	lblEngine.SetHAlign(gtk.ALIGN_START)
	w.ComboEngine, _ = createEngineDropdown(server.Engine)
	overridesGrid.Attach(lblEngine, 0, or, 1, 1)
	overridesGrid.Attach(w.ComboEngine, 1, or, 1, 1)
	or++

	lblIgnoreCert, _ := gtk.LabelNew("Ignore Certificate:")
	lblIgnoreCert.SetHAlign(gtk.ALIGN_START)
	w.ComboIgnoreCert, _ = createThreeWayDropdown(server.IgnoreCertificate)
	overridesGrid.Attach(lblIgnoreCert, 0, or, 1, 1)
	overridesGrid.Attach(w.ComboIgnoreCert, 1, or, 1, 1)
	or++

	lblTLS, _ := gtk.LabelNew("TLS Security Level:")
	lblTLS.SetHAlign(gtk.ALIGN_START)
	w.ComboTLS, _ = createTLSDropdown(server.TLSSecLevel)
	overridesGrid.Attach(lblTLS, 0, or, 1, 1)
	overridesGrid.Attach(w.ComboTLS, 1, or, 1, 1)
	or++

	lblClipboard, _ := gtk.LabelNew("Use Clipboard:")
	lblClipboard.SetHAlign(gtk.ALIGN_START)
	w.ComboClipboard, _ = createThreeWayDropdown(server.Clipboard)
	overridesGrid.Attach(lblClipboard, 0, or, 1, 1)
	overridesGrid.Attach(w.ComboClipboard, 1, or, 1, 1)
	or++

	lblSecNLA, _ := gtk.LabelNew("Network Level Auth (NLA):")
	lblSecNLA.SetHAlign(gtk.ALIGN_START)
	w.ComboSecNLA, _ = createThreeWayDropdown(server.SecNLA)
	overridesGrid.Attach(lblSecNLA, 0, or, 1, 1)
	overridesGrid.Attach(w.ComboSecNLA, 1, or, 1, 1)
	or++

	lblFullscreen, _ := gtk.LabelNew("Fullscreen:")
	lblFullscreen.SetHAlign(gtk.ALIGN_START)
	w.ComboFullscreen, _ = createThreeWayDropdown(server.Fullscreen)
	overridesGrid.Attach(lblFullscreen, 0, or, 1, 1)
	overridesGrid.Attach(w.ComboFullscreen, 1, or, 1, 1)
	or++

	lblDynamicRes, _ := gtk.LabelNew("Dynamic Resolution:")
	lblDynamicRes.SetHAlign(gtk.ALIGN_START)
	w.ComboDynamicRes, _ = createThreeWayDropdown(server.DynamicRes)
	overridesGrid.Attach(lblDynamicRes, 0, or, 1, 1)
	overridesGrid.Attach(w.ComboDynamicRes, 1, or, 1, 1)
	or++

	lblMultimon, _ := gtk.LabelNew("Multimonitor:")
	lblMultimon.SetHAlign(gtk.ALIGN_START)
	w.ComboMultimon, _ = createThreeWayDropdown(server.Multimon)
	overridesGrid.Attach(lblMultimon, 0, or, 1, 1)
	overridesGrid.Attach(w.ComboMultimon, 1, or, 1, 1)
	or++

	lblSound, _ := gtk.LabelNew("Redirect Audio:")
	lblSound.SetHAlign(gtk.ALIGN_START)
	w.ComboSound, _ = createThreeWayDropdown(server.Sound)
	overridesGrid.Attach(lblSound, 0, or, 1, 1)
	overridesGrid.Attach(w.ComboSound, 1, or, 1, 1)
	or++

	lblShareHome, _ := gtk.LabelNew("Share Home Drive:")
	lblShareHome.SetHAlign(gtk.ALIGN_START)
	w.ComboShareHome, _ = createThreeWayDropdown(server.ShareHome)
	overridesGrid.Attach(lblShareHome, 0, or, 1, 1)
	overridesGrid.Attach(w.ComboShareHome, 1, or, 1, 1)
	or++

	lblFonts, _ := gtk.LabelNew("Font Smoothing:")
	lblFonts.SetHAlign(gtk.ALIGN_START)
	w.ComboFonts, _ = createThreeWayDropdown(server.FontSmoothing)
	overridesGrid.Attach(lblFonts, 0, or, 1, 1)
	overridesGrid.Attach(w.ComboFonts, 1, or, 1, 1)
	or++

	lblWallpaper, _ := gtk.LabelNew("Wallpaper:")
	lblWallpaper.SetHAlign(gtk.ALIGN_START)
	w.ComboWallpaper, _ = createThreeWayDropdown(server.Wallpaper)
	overridesGrid.Attach(lblWallpaper, 0, or, 1, 1)
	overridesGrid.Attach(w.ComboWallpaper, 1, or, 1, 1)
	or++

	lblThemes, _ := gtk.LabelNew("Themes:")
	lblThemes.SetHAlign(gtk.ALIGN_START)
	w.ComboThemes, _ = createThreeWayDropdown(server.Themes)
	overridesGrid.Attach(lblThemes, 0, or, 1, 1)
	overridesGrid.Attach(w.ComboThemes, 1, or, 1, 1)
	or++

	lblParams, _ := gtk.LabelNew("Custom Params:")
	lblParams.SetHAlign(gtk.ALIGN_START)
	w.EntryParams, _ = gtk.EntryNew()
	w.EntryParams.SetPlaceholderText("default setting")
	w.EntryParams.SetHExpand(true)
	w.EntryParams.SetText(server.CustomParams)
	overridesGrid.Attach(lblParams, 0, or, 1, 1)
	overridesGrid.Attach(w.EntryParams, 1, or, 1, 1)
	or++

	expander, err := gtk.ExpanderNew("Advanced Parameter Overrides")
	if err != nil {
		return nil, nil, err
	}
	expander.SetMarginTop(10)
	expander.SetMarginBottom(10)
	expander.Add(overridesGrid)

	grid.Attach(expander, 0, r, 2, 1)
	r++

	return grid, w, nil
}

func setupServerDialog(parentWindow *gtk.Window, title string) (*gtk.Dialog, *gtk.ScrolledWindow, error) {
	dialog, err := gtk.DialogNew()
	if err != nil {
		return nil, nil, err
	}
	dialog.SetTitle(title)
	dialog.SetTransientFor(parentWindow)
	dialog.SetModal(true)
	dialog.SetDefaultSize(550, 500)

	dialog.AddButton("Cancel", gtk.RESPONSE_CANCEL)
	dialog.AddButton("Save", gtk.RESPONSE_OK)

	contentArea, err := dialog.GetContentArea()
	if err != nil {
		dialog.Destroy()
		return nil, nil, err
	}

	scroll, err := gtk.ScrolledWindowNew(nil, nil)
	if err != nil {
		dialog.Destroy()
		return nil, nil, err
	}
	scroll.SetPolicy(gtk.POLICY_NEVER, gtk.POLICY_AUTOMATIC)
	
	contentArea.PackStart(scroll, true, true, 0)
	
	return dialog, scroll, nil
}

func showAddServerDialog(app *AppUI) {
	dialog, scroll, err := setupServerDialog(app.Window, "Add Server")
	if err != nil {
		fmt.Printf("Failed to create Add Server dialog: %v\n", err)
		return
	}
	defer dialog.Destroy()

	grid, widgets, err := buildServerFormGrid(nil)
	if err != nil {
		fmt.Printf("Failed to build server form: %v\n", err)
		return
	}
	scroll.Add(grid)
	dialog.ShowAll()

	response := dialog.Run()
	if response == gtk.RESPONSE_OK {
		var newServer config.ServerConfig
		widgets.ApplyToServer(&newServer)

		if strings.TrimSpace(newServer.Name) == "" || strings.TrimSpace(newServer.HostIP) == "" {
			fmt.Println("[Add Server] Failed: Server name and Host/IP cannot be empty.")
			return
		}

		newServer.ID = fmt.Sprintf("srv-%d", time.Now().UnixNano()/int64(time.Millisecond))
		app.Config.Servers = append(app.Config.Servers, newServer)

		if err := config.SaveConfig(*app.Config); err != nil {
			fmt.Printf("[Config] Failed to save new server: %v\n", err)
		} else {
			fmt.Printf("[Config] Server '%s' added successfully.\n", newServer.Name)
			if err := app.PopulateServerList(); err != nil {
				fmt.Printf("Failed to refresh server list: %v\n", err)
			}
		}
	}
}

func showEditServerDialog(app *AppUI, serverID string) {
	var targetIdx = -1
	for i, s := range app.Config.Servers {
		if s.ID == serverID {
			targetIdx = i
			break
		}
	}
	if targetIdx == -1 {
		fmt.Printf("[Edit] Error: Server ID %s not found.\n", serverID)
		return
	}
	server := app.Config.Servers[targetIdx]

	dialog, scroll, err := setupServerDialog(app.Window, "Edit Server")
	if err != nil {
		fmt.Printf("Failed to create Edit Server dialog: %v\n", err)
		return
	}
	defer dialog.Destroy()

	grid, widgets, err := buildServerFormGrid(&server)
	if err != nil {
		fmt.Printf("Failed to build server form: %v\n", err)
		return
	}
	scroll.Add(grid)
	dialog.ShowAll()

	response := dialog.Run()
	if response == gtk.RESPONSE_OK {
		widgets.ApplyToServer(&server)

		if strings.TrimSpace(server.Name) == "" || strings.TrimSpace(server.HostIP) == "" {
			fmt.Println("[Edit Server] Failed: Server name and Host/IP cannot be empty.")
			return
		}

		app.Config.Servers[targetIdx] = server

		if err := config.SaveConfig(*app.Config); err != nil {
			fmt.Printf("[Config] Failed to save edited server: %v\n", err)
		} else {
			fmt.Printf("[Config] Changes to server '%s' saved successfully.\n", server.Name)
			if err := app.PopulateServerList(); err != nil {
				fmt.Printf("Failed to refresh server list: %v\n", err)
			}
		}
	}
}

func showContextMenu(app *AppUI, serverID string, event *gdk.EventButton) {
	menu, err := gtk.MenuNew()
	if err != nil {
		fmt.Printf("Failed to create menu: %v\n", err)
		return
	}

	var serverName string
	for _, s := range app.Config.Servers {
		if s.ID == serverID {
			serverName = s.Name
			break
		}
	}

	editItem, err := gtk.MenuItemNewWithLabel("Edit Server")
	if err == nil {
		editItem.Connect("activate", func() {
			showEditServerDialog(app, serverID)
		})
		menu.Append(editItem)
	}

	deleteItem, err := gtk.MenuItemNewWithLabel("Delete Server")
	if err == nil {
		deleteItem.Connect("activate", func() {
			if confirmDeleteDialog(app.Window, serverName) {
				for i, item := range app.Config.Servers {
					if item.ID == serverID {
						app.Config.Servers = append(app.Config.Servers[:i], app.Config.Servers[i+1:]...)
						break
					}
				}

				if err := config.SaveConfig(*app.Config); err != nil {
					fmt.Printf("[Config] Failed to save after deletion: %v\n", err)
				} else {
					fmt.Printf("[Config] Server '%s' deleted successfully.\n", serverName)
					app.PopulateServerList()
				}
			}
		})
		menu.Append(deleteItem)
	}

	menu.ShowAll()
	menu.PopupAtPointer(event.Event)
}

func promptPasswordDialog(parentWindow *gtk.Window, connectionName string) (string, bool) {
	dialog, err := gtk.DialogNew()
	if err != nil {
		fmt.Printf("Failed to create prompt password dialog: %v\n", err)
		return "", false
	}
	defer dialog.Destroy()

	dialog.SetTitle("Password Required")
	dialog.SetTransientFor(parentWindow)
	dialog.SetModal(true)
	dialog.SetDefaultSize(350, 150)

	dialog.AddButton("Cancel", gtk.RESPONSE_CANCEL)
	dialog.AddButton("Connect", gtk.RESPONSE_OK)

	contentArea, err := dialog.GetContentArea()
	if err != nil {
		return "", false
	}

	box, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 10)
	if err != nil {
		return "", false
	}
	box.SetMarginStart(15)
	box.SetMarginEnd(15)
	box.SetMarginTop(15)
	box.SetMarginBottom(15)

	lblText, err := gtk.LabelNew(fmt.Sprintf("Enter password for connection to '%s':", connectionName))
	if err != nil {
		return "", false
	}
	lblText.SetHAlign(gtk.ALIGN_START)

	entryPass, err := gtk.EntryNew()
	if err != nil {
		return "", false
	}
	entryPass.SetVisibility(false)
	entryPass.SetActivatesDefault(true) // Pressing enter triggers default action (Connect)

	chkShowPass, err := gtk.CheckButtonNewWithLabel("Show Password")
	if err != nil {
		return "", false
	}
	chkShowPass.Connect("toggled", func() {
		entryPass.SetVisibility(chkShowPass.GetActive())
	})

	box.PackStart(lblText, false, false, 0)
	box.PackStart(entryPass, false, false, 0)
	box.PackStart(chkShowPass, false, false, 0)

	contentArea.Add(box)
	dialog.ShowAll()

	dialog.SetDefaultResponse(gtk.RESPONSE_OK)

	response := dialog.Run()
	if response == gtk.RESPONSE_OK {
		pass, _ := entryPass.GetText()
		return pass, true
	}

	return "", false
}
