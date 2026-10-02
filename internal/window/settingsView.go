package window

import (
	"fmt"
	"image/color"
	"os"
	"password_manager/internal/credentialStore"
	"password_manager/internal/customLayout"
	"password_manager/internal/settingStore"
	"password_manager/internal/utils"
	"password_manager/internal/validator"
	"password_manager/internal/vaultStore"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// TODO: check validators in here to make them the same as in addVaultView
func settingsView(app fyne.App, window fyne.Window, settings *settingStore.GeneralSettings, credentialManager *credentialStore.CredentialManager, vaultManager *vaultStore.VaultManager, vaultId int) *fyne.Container {
	var applyButton *widget.Button

	// Settings title
	settingsTitle := canvas.NewText("Settings", color.White)
	settingsTitle.TextStyle.Bold = true
	settingsTitle.TextSize = 25

	// Settings icon
	settingsIcon := canvas.NewImageFromResource(theme.SettingsIcon())
	settingsIcon.SetMinSize(fyne.NewSize(50, 50))

	// title for general section
	generalTitle := canvas.NewText("General", color.White)
	generalTitle.TextStyle.Bold = true
	generalTitle.TextSize = 20

	// label and values for general section
	copyTimeLabel := widget.NewLabel("Copy time")
	copyTimeValue := widget.NewEntry()

	// fill entries with data for general section
	copyTimeValue.SetText(strconv.Itoa(settings.CopyTime))

	// handle validators for general section
	copyTimeValue.Validator = func(s string) error {
		if s == "" {
			applyButton.Disable()
			return fmt.Errorf("can't be empty")
		} else if _, err := strconv.Atoi(s); err != nil {
			applyButton.Disable()
			return fmt.Errorf("needs to be a number")
		}

		applyButton.Enable()
		return nil
	}

	// create containers for general section
	formContainerGeneralSettings := container.NewVBox(
		generalTitle,
		container.New(
			layout.NewFormLayout(),
			copyTimeLabel,
			copyTimeValue,
		),
	)
	formContainerGeneralSettings = container.New(
		layout.NewCustomPaddedLayout(
			0,
			20,
			0,
			0,
		),
		formContainerGeneralSettings,
	)

	// title for vault section
	vaultTitle := canvas.NewText("Vault", color.White)
	vaultTitle.TextStyle.Bold = true
	vaultTitle.TextSize = 20

	// label and values for vault section
	vaultNameLabel := widget.NewLabel("Vault name")
	encryptionKeyLabel := widget.NewLabel("Encryption key")
	dataSaveLocationLabel := widget.NewLabel("Data save location")
	encryptionKeyValue := widget.NewPasswordEntry()
	vaultNameValue := widget.NewEntry()
	dataSaveLocationValue := widget.NewEntry()

	// fill entries with data for vault section
	vaultNameValue.SetText(vaultManager.Vaults[vaultId].Name)
	encryptionKeyValue.SetText(credentialManager.GetEncryptionKey())
	dataSaveLocationValue.SetText(vaultManager.Vaults[vaultId].SaveLocation)

	// handle validators for vault section
	vaultNameValue.Validator = validator.EmptyStringValidator
	encryptionKeyValue.Validator = validator.EmptyStringValidator
	dataSaveLocationValue.Validator = validator.EmptyStringValidator

	// create containers for vault section
	formContainerVaultSettings := container.NewVBox(
		vaultTitle,
		container.New(
			layout.NewFormLayout(),
			vaultNameLabel,
			vaultNameValue,
			encryptionKeyLabel,
			encryptionKeyValue,
			dataSaveLocationLabel,
			dataSaveLocationValue,
		),
	)
	formContainerVaultSettings = container.New(
		layout.NewCustomPaddedLayout(
			0,
			20,
			0,
			0,
		),
		formContainerVaultSettings,
	)

	// handle on change events to disable/enable button
	disableButtonOnInvalidEntries := func(s string) {
		if copyTimeValue.Validate() != nil || vaultNameValue.Validate() != nil || encryptionKeyValue.Validate() != nil || dataSaveLocationValue.Validate() != nil {
			applyButton.Disable()
			return
		}
		applyButton.Enable()
	}
	copyTimeValue.OnChanged = disableButtonOnInvalidEntries
	vaultNameValue.OnChanged = disableButtonOnInvalidEntries
	encryptionKeyValue.OnChanged = disableButtonOnInvalidEntries
	dataSaveLocationValue.OnChanged = disableButtonOnInvalidEntries

	// Create button
	returnButton := widget.NewButtonWithIcon("Go back", theme.MailReplyIcon(), func() {
		window.SetContent(mainView(app, window, settings, credentialManager, vaultManager, vaultId))
	})
	applyButton = widget.NewButton("Apply", func() {
		generalChangesDetected := false
		vaultChangesDetected := false

		// check if copytime changed
		copyTimeValueConverted, err := strconv.Atoi(copyTimeValue.Text)
		if err != nil {
			dialog.NewError(err, window).Show()
		} else {
			if settings.CopyTime != copyTimeValueConverted {
				settings.CopyTime = copyTimeValueConverted

				if !generalChangesDetected {
					generalChangesDetected = true
				}
			}
		}

		// handle vault name change
		if vaultManager.Vaults[vaultId].Name != vaultNameValue.Text {
			vaultManager.Vaults[vaultId].Name = vaultNameValue.Text

			if !vaultChangesDetected {
				vaultChangesDetected = true
			}
		}

		// check if encryption key changed
		if credentialManager.GetEncryptionKey() != encryptionKeyValue.Text {
			err := credentialManager.SetEncryptionKey(encryptionKeyValue.Text)
			if err != nil {
				dialog.NewError(err, window).Show()
			} else {
				dialog.NewInformation("Changed encryption key", "Encryption key changed and all of the credentials will be encrypted with the new key on save", window).Show()
			}
		}

		// check if save file location changed
		// move save file if locations has changed
		if vaultManager.Vaults[vaultId].SaveLocation != dataSaveLocationValue.Text {
			originalFileContent, err := os.ReadFile(vaultManager.Vaults[vaultId].SaveLocation)
			if err != nil {
				dialog.NewError(err, window).Show()
			} else {
				// Check if file already exists
				err := utils.CheckAndCreateFile(dataSaveLocationValue.Text)
				if err != nil {
					dialog.NewError(err, window).Show()
				} else {
					// Write data from old file to new file
					err := os.WriteFile(dataSaveLocationValue.Text, originalFileContent, 0644)
					if err != nil {
						dialog.NewError(err, window).Show()
					} else {
						// Remove original file
						err := os.Remove(vaultManager.Vaults[vaultId].SaveLocation)
						if err != nil {
							dialog.NewError(err, window).Show()
						} else {
							// Update vault
							vaultManager.Vaults[vaultId].SaveLocation = dataSaveLocationValue.Text

							if !vaultChangesDetected {
								vaultChangesDetected = true
							}
						}
					}
				}
			}
		}

		if generalChangesDetected {
			if err := settings.WriteToFile(); err != nil {
				dialog.NewError(err, window).Show()
			}
		}

		if vaultChangesDetected {
			if err := vaultManager.WriteToFile(); err != nil {
				dialog.NewError(err, window).Show()
			}
		}
	})

	// Set button importance
	applyButton.Importance = widget.HighImportance

	// Create container
	returnButtonContainer := container.New(
		layout.NewBorderLayout(
			nil,
			nil,
			nil,
			returnButton,
		),
		returnButton,
	)
	settingsTitleWithIcon := container.NewCenter(
		container.NewHBox(
			settingsIcon,
			settingsTitle,
		),
	)
	settingsContainer := container.New(
		layout.NewBorderLayout(
			returnButtonContainer,
			applyButton,
			nil,
			nil,
		),
		returnButtonContainer,
		container.NewScroll(
			container.NewVBox(
				settingsTitleWithIcon,
				formContainerGeneralSettings,
				formContainerVaultSettings,
			),
		),
		applyButton,
	)

	return container.New(
		customLayout.NewResponsivePaddingLayout(0.05),
		settingsContainer,
	)
}
