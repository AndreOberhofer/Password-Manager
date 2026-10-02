package window

import (
	"errors"
	"image/color"
	"os"
	"password_manager/internal/credentialStore"
	"password_manager/internal/customDialog"
	"password_manager/internal/customLayout"
	"password_manager/internal/settingStore"
	"password_manager/internal/vaultStore"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func VaultView(app fyne.App, window fyne.Window) *fyne.Container {
	var isCriticalError bool = false

	// Create vault manager
	vaultManager := vaultStore.NewVaultManager()
	err := vaultManager.ReadFromFile()
	if err != nil {
		dialog.NewError(err, window).Show()
		isCriticalError = true
	}

	// Create settings manager
	settings := settingStore.NewSettings()
	err = settings.ReadFromFile()
	if err != nil {
		dialog.NewError(err, window).Show()
		isCriticalError = true
	}

	// Create text
	titleText := canvas.NewText("Vaults", color.White)
	titleText.TextSize = 30
	titleText.Alignment = fyne.TextAlignCenter
	titleText.TextStyle.Bold = true

	// Create toolbar actions
	refreshToolbarAction := widget.NewToolbarAction(theme.ViewRefreshIcon(), func() {
		if err := vaultManager.ReadFromFile(); err != nil {
			dialog.NewError(err, window).Show()
			return
		}
		window.Content().Refresh()
	})
	addToolbarAction := widget.NewToolbarAction(theme.ContentAddIcon(), func() {
		customDialog.AddVaultDialog("Add vault", vaultManager, window).Show()
	})
	addFileToolbarAction := widget.NewToolbarAction(theme.DownloadIcon(), func() {
		addFileDialog := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.NewError(err, window).Show()
				return
			}
			if reader != nil {
				credentialManager := credentialStore.NewCredentialManager()
				err = credentialManager.ReadFromFile(reader.URI().Path())
				if err != nil {
					dialog.NewError(err, window).Show()
					return
				}
				customDialog.AddVaultFromFileDialog("Add vault from file", vaultManager, credentialManager, window).Show()
			}
		}, window)
		addFileDialog.Show()
	})

	// Create toolbar
	toolbar := widget.NewToolbar(
		refreshToolbarAction,
		addToolbarAction,
		addFileToolbarAction,
	)

	// Disable components on critical error
	if isCriticalError {
		refreshToolbarAction.Disable()
		addToolbarAction.Disable()
		addFileToolbarAction.Disable()
	}

	// Create list
	vaultList := widget.NewList(
		func() int {
			return len(vaultManager.Vaults)
		},
		func() fyne.CanvasObject {
			label := canvas.NewText("template", color.White)

			openButton := widget.NewButtonWithIcon("", theme.FolderOpenIcon(), func() {})
			deleteButton := widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {})

			buttonContainer := container.NewHBox(
				openButton,
				deleteButton,
			)

			return container.New(
				layout.NewBorderLayout(
					nil,
					nil,
					nil,
					buttonContainer,
				),
				label,
				buttonContainer,
			)
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			// retrieve different components from canvas object
			container := o.(*fyne.Container)
			label := container.Objects[0].(*canvas.Text)
			buttonContainer := container.Objects[1].(*fyne.Container)
			openButton := buttonContainer.Objects[0].(*widget.Button)
			deleteButton := buttonContainer.Objects[1].(*widget.Button)

			// set label
			label.Text = vaultManager.Vaults[i].Name
			label.TextStyle.Bold = true
			// turn label red if save file no longer exists
			if _, err := os.Stat(vaultManager.Vaults[i].SaveLocation); errors.Is(err, os.ErrNotExist) {
				label.Color = color.RGBA{R: 255, G: 0, B: 0, A: 255}
			}

			// handle button events
			openButton.OnTapped = func() {
				window.SetContent(loginView(app, window, vaultManager, i, settings))
			}
			deleteButton.OnTapped = func() {
				if i >= 0 && i < len(vaultManager.Vaults) {
					customDialog.ConfirmDialog(
						"Delete item",
						"Are you sure you want to delete the following item:\n"+vaultManager.Vaults[i].Name,
						func() {
							// remove credential file if it still exists
							_, err := os.Stat(vaultManager.Vaults[i].SaveLocation)
							if !errors.Is(err, os.ErrNotExist) {
								err = os.Remove(vaultManager.Vaults[i].SaveLocation)
								if err != nil {
									dialog.NewError(err, window).Show()
									return
								}
							}

							// remove vault from slice and update vault file
							vaultManager.Vaults = append(vaultManager.Vaults[:i], vaultManager.Vaults[i+1:]...) // TODO: could be made into a function
							if err := vaultManager.WriteToFile(); err != nil {
								dialog.NewError(err, window).Show()
								return
							}

							window.Content().Refresh()
						},
						window).Show()
				}
			}
		},
	)

	// Create container
	titleAndToolbar := container.NewVBox(
		titleText,
		toolbar,
	)

	return container.New(
		customLayout.NewResponsivePaddingLayout(0.15),
		container.New(
			layout.NewBorderLayout(
				titleAndToolbar,
				nil,
				nil,
				nil,
			),
			titleAndToolbar,
			vaultList,
		),
	)
}
