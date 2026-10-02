package customDialog

import (
	"password_manager/internal/credentialStore"
	"password_manager/internal/utils"
	"password_manager/internal/validator"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func AddFolderDialog(title string, parentBranchId string, credentialManager *credentialStore.CredentialManager, window fyne.Window) *dialog.CustomDialog {
	var newDialog *dialog.CustomDialog
	var addButton *widget.Button

	// Create entry field
	folderNameEntry := widget.NewEntry()

	// Set place holder text
	folderNameEntry.PlaceHolder = "Enter folder name..."

	// Handle entry validator
	folderNameEntry.Validator = validator.EmptyStringValidator

	// Create validator text
	folderNameEntryValidatorText := utils.CreateEntryErrorDisplayText()

	// Handle validator text
	validator.SetValidatorEntryText(folderNameEntry, folderNameEntryValidatorText)

	// Handle on change events to enable/disable button
	disableButtonOnInvalidEntries := func(s string) {
		if folderNameEntry.Validate() != nil {
			addButton.Disable()
			return
		}
		addButton.Enable()
	}
	folderNameEntry.OnChanged = disableButtonOnInvalidEntries

	// Create button
	addButton = widget.NewButton("Add", func() {
		credentialManager.AddNewTreeBranch(folderNameEntry.Text, parentBranchId)
		window.Content().Refresh()
	})
	addButton.Disable()
	cancelButton := widget.NewButton("Cancel", func() {
		newDialog.Dismiss()
	})

	// Create container
	folderNameEntryContainer := container.New(
		layout.NewBorderLayout(
			nil,
			folderNameEntryValidatorText,
			nil,
			nil,
		),
		folderNameEntry,
		folderNameEntryValidatorText,
	)
	content := container.NewVBox(
		folderNameEntryContainer,
		container.NewHBox(addButton, cancelButton),
	)

	newDialog = dialog.NewCustomWithoutButtons(title, content, window)
	newDialog.Resize(fyne.NewSize(300, 0)) // TODO: make dynamic
	return newDialog
}
