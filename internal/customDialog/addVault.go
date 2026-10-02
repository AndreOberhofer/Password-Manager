package customDialog

import (
	"os"
	"password_manager/internal/credentialStore"
	"password_manager/internal/utils"
	"password_manager/internal/validator"
	"password_manager/internal/vaultStore"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func AddVaultDialog(title string, vaultManager *vaultStore.VaultManager, window fyne.Window) *dialog.CustomDialog {
	var newDialog *dialog.CustomDialog
	var addButton *widget.Button

	// Create entry field
	nameEntry := widget.NewEntry()
	encryptionKeyEntry := widget.NewPasswordEntry()
	saveLocationEntry := widget.NewEntry()

	// Set place holder text
	nameEntry.PlaceHolder = "Enter name..."
	encryptionKeyEntry.PlaceHolder = "Enter encryption key..."
	saveLocationEntry.PlaceHolder = "Enter save location..."

	// Handle entry validator
	nameEntry.Validator = validator.EmptyStringValidator
	encryptionKeyEntry.Validator = validator.EmptyStringValidator
	saveLocationEntry.Validator = validator.EmptyStringValidator

	// Create validator text
	nameEntryValidatorText := utils.CreateEntryErrorDisplayText()
	encryptionKeyEntryValidatorText := utils.CreateEntryErrorDisplayText()
	saveLocationEntryEntryValidatorText := utils.CreateEntryErrorDisplayText()

	// Handle validator text
	validator.SetValidatorEntryText(nameEntry, nameEntryValidatorText)
	validator.SetValidatorEntryText(encryptionKeyEntry, encryptionKeyEntryValidatorText)
	validator.SetValidatorEntryText(saveLocationEntry, saveLocationEntryEntryValidatorText)

	// Handle on change events to enable/disable button
	disableButtonOnInvalidEntries := func(s string) {
		if nameEntry.Validate() != nil || encryptionKeyEntry.Validate() != nil || saveLocationEntry.Validate() != nil {
			addButton.Disable()
			return
		}
		addButton.Enable()
	}
	nameEntry.OnChanged = disableButtonOnInvalidEntries
	encryptionKeyEntry.OnChanged = disableButtonOnInvalidEntries
	saveLocationEntry.OnChanged = disableButtonOnInvalidEntries

	// Create button
	chooseFileButton := widget.NewButtonWithIcon("", theme.FolderOpenIcon(), func() {
		chooseFileDialog := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
			if writer != nil {
				// Close reader and delete file
				writer.Close()
				os.Remove(writer.URI().Path())

				// Put path into entry field
				saveLocationEntry.SetText(writer.URI().Path())
			}
		}, window)

		fileName := "data.json"
		if nameEntry.Text != "" {
			fileName = strings.Join([]string{nameEntry.Text, "data.json"}, "-")
		}

		chooseFileDialog.SetFileName(fileName)
		chooseFileDialog.Show()
	})
	addButton = widget.NewButton("Add", func() {
		saveLocation := saveLocationEntry.Text

		// Create and save credential manager
		newCredentialManager := credentialStore.NewCredentialManager()
		err := newCredentialManager.SetEncryptionKey(encryptionKeyEntry.Text)
		if err != nil {
			dialog.NewError(err, window).Show()
			return
		}
		err = newCredentialManager.WriteToFile(saveLocation)
		if err != nil {
			dialog.NewError(err, window).Show()
			return
		}

		// Create new vault instance
		newVault := vaultStore.Vault{ // TODO: could create a constructor for this
			Name:         nameEntry.Text,
			SaveLocation: saveLocation,
		}

		// Add vault to manager and save file
		vaultManager.Vaults = append(vaultManager.Vaults, newVault) // TODO: could do this with a function
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
	encryptionKeyContainer := container.New(
		layout.NewBorderLayout(
			nil,
			encryptionKeyEntryValidatorText,
			nil,
			nil,
		),
		encryptionKeyEntry,
		encryptionKeyEntryValidatorText,
	)
	saveLocationEntryContainer := container.New(
		layout.NewBorderLayout(
			nil,
			saveLocationEntryEntryValidatorText,
			nil,
			chooseFileButton,
		),
		saveLocationEntry,
		saveLocationEntryEntryValidatorText,
		chooseFileButton,
	)
	content := container.NewVBox(
		nameEntryContainer,
		encryptionKeyContainer,
		saveLocationEntryContainer,
		container.NewHBox(addButton, cancelButton),
	)

	newDialog = dialog.NewCustomWithoutButtons(title, content, window)
	newDialog.Resize(fyne.NewSize(300, 0)) // TODO: make dynamic
	return newDialog
}
