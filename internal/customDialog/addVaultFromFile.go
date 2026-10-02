package customDialog

import (
	"password_manager/internal/credentialStore"
	"password_manager/internal/utils"
	"password_manager/internal/validator"
	"password_manager/internal/vaultStore"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func AddVaultFromFileDialog(title string, vaultManager *vaultStore.VaultManager, credentialManager *credentialStore.CredentialManager, window fyne.Window) *dialog.CustomDialog {
	var newDialog *dialog.CustomDialog
	var addButton *widget.Button

	// Create entry field
	nameEntry := widget.NewEntry()
	saveLocationEntry := widget.NewEntry()

	// Set place holder text
	nameEntry.PlaceHolder = "Enter name..."
	saveLocationEntry.PlaceHolder = "Enter save location..."

	// Handle entry validator
	nameEntry.Validator = validator.EmptyStringValidator
	saveLocationEntry.Validator = validator.EmptyStringValidator

	// Create validator text
	nameEntryValidatorText := utils.CreateEntryErrorDisplayText()
	saveLocationEntryEntryValidatorText := utils.CreateEntryErrorDisplayText()

	// Handle validator text
	validator.SetValidatorEntryText(nameEntry, nameEntryValidatorText)
	validator.SetValidatorEntryText(saveLocationEntry, saveLocationEntryEntryValidatorText)

	// Handle on change events to enable/disable button
	disableButtonOnInvalidEntries := func(s string) {
		if nameEntry.Validate() != nil || saveLocationEntry.Validate() != nil {
			addButton.Disable()
			return
		}
		addButton.Enable()
	}
	nameEntry.OnChanged = disableButtonOnInvalidEntries
	saveLocationEntry.OnChanged = disableButtonOnInvalidEntries

	// Create button
	addButton = widget.NewButton("Add", func() {
		newVault := vaultStore.Vault{
			Name:         nameEntry.Text,
			SaveLocation: saveLocationEntry.Text,
		}

		// Add vault manager and save file
		vaultManager.Vaults = append(vaultManager.Vaults, newVault)
		if err := vaultManager.WriteToFile(); err != nil {
			dialog.NewError(err, window).Show()
			return
		}

		window.Content().Refresh()
		newDialog.Dismiss()
	})
	addButton.Disable()
	cancelButton := widget.NewButton("Cancel", func() {
		newDialog.Dismiss()
	})

	// Create container
	nameEntryContainer := container.New(
		layout.NewBorderLayout(
			nil,
			nameEntryValidatorText,
			nil,
			nil,
		),
		nameEntry,
		nameEntryValidatorText,
	)
	saveLocationEntryContainer := container.New(
		layout.NewBorderLayout(
			nil,
			saveLocationEntryEntryValidatorText,
			nil,
			nil,
		),
		saveLocationEntry,
		saveLocationEntryEntryValidatorText,
	)
	content := container.NewVBox(
		nameEntryContainer,
		saveLocationEntryContainer,
		container.NewHBox(addButton, cancelButton),
	)

	newDialog = dialog.NewCustomWithoutButtons(title, content, window)
	newDialog.Resize(fyne.NewSize(300, 0)) // TODO: make dynamic
	return newDialog
}
