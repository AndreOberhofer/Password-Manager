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
	"github.com/google/uuid"
)

func AddCredentialDialog(title string, credentialManager *credentialStore.CredentialManager, branchId string, window fyne.Window) *dialog.CustomDialog {
	var newDialog *dialog.CustomDialog
	var addButton *widget.Button

	// Create entry field
	titleEntry := widget.NewEntry()
	usernameEntry := widget.NewEntry()
	emailEntry := widget.NewEntry()
	passwordEntry := widget.NewPasswordEntry()

	// Set place holder text
	titleEntry.PlaceHolder = "Enter title..."
	usernameEntry.PlaceHolder = "Enter username..."
	emailEntry.PlaceHolder = "Enter email..."
	passwordEntry.PlaceHolder = "Enter password..."

	// Handle entry validator
	titleEntry.Validator = validator.EmptyStringValidator

	// Create validator text
	titleEntryValidatorText := utils.CreateEntryErrorDisplayText()

	// Handle validator text
	validator.SetValidatorEntryText(titleEntry, titleEntryValidatorText)

	// Create combined entry container
	titleEntryContainer := container.New(
		layout.NewBorderLayout(
			nil,
			titleEntryValidatorText,
			nil,
			nil,
		),
		titleEntry,
		titleEntryValidatorText,
	)

	// Hide entries
	hideAllEntries := func() {
		titleEntryContainer.Hide()
		usernameEntry.Hide()
		emailEntry.Hide()
		passwordEntry.Hide()
	}
	hideAllEntries()

	// Create dropdown select
	credentialTypeSelect := widget.NewSelect(
		[]string{
			credentialStore.Password.String(),
			credentialStore.Email.String(),
			credentialStore.User.String(),
		},
		func(s string) {
			// Show base fields
			hideAllEntries()
			titleEntryContainer.Show()

			// Show special fields
			switch s {
			case credentialStore.Password.String():
				passwordEntry.Show()
			case credentialStore.Email.String():
				emailEntry.Show()
				passwordEntry.Show()
			case credentialStore.User.String():
				usernameEntry.Show()
				passwordEntry.Show()
			}
		},
	)

	// Handle on change events to enable/disable button
	disableButtonOnInvalidEntries := func(s string) {
		if titleEntry.Validate() != nil {
			addButton.Disable()
			return
		}
		addButton.Enable()
	}
	titleEntry.OnChanged = disableButtonOnInvalidEntries

	// Create button
	addButton = widget.NewButton("Add", func() {
		err := credentialManager.AddCredential(credentialStore.NewCredential(uuid.NewString(), branchId, titleEntry.Text, credentialStore.CredentialType(credentialTypeSelect.SelectedIndex()), emailEntry.Text, usernameEntry.Text, passwordEntry.Text))
		if err != nil {
			dialog.NewError(err, window).Show()
			return
		}
		window.Content().Refresh()
	})
	cancelButton := widget.NewButton("Cancel", func() {
		newDialog.Dismiss()
	})

	// Disable add button
	addButton.Disable()

	// Create container
	entryFieldsContainer := container.NewVBox(
		titleEntryContainer,
		usernameEntry,
		emailEntry,
		passwordEntry,
	)
	entryFieldsContainer = container.New(
		layout.NewCustomPaddedLayout(
			10,
			10,
			0,
			0,
		),
		entryFieldsContainer,
	)
	content := container.NewVBox(
		credentialTypeSelect,
		entryFieldsContainer,
		container.NewHBox(addButton, cancelButton),
	)

	newDialog = dialog.NewCustomWithoutButtons(title, content, window)
	newDialog.Resize(fyne.NewSize(300, 0)) // TODO: make dynamic
	return newDialog
}
