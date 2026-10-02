package window

import (
	"image/color"
	"password_manager/internal/credentialStore"
	"password_manager/internal/customLayout"
	"password_manager/internal/settingStore"
	"password_manager/internal/utils"
	"password_manager/internal/vaultStore"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"golang.org/x/crypto/bcrypt"
)

func loginView(app fyne.App, window fyne.Window, vaultManager *vaultStore.VaultManager, vaultIndex int, settings *settingStore.GeneralSettings) *fyne.Container {
	var isCriticalError bool = false

	// create and fill credential manager
	credentialManager := credentialStore.NewCredentialManager()
	err := credentialManager.ReadEncryptionHashFromFile(vaultManager.Vaults[vaultIndex].SaveLocation)
	if err != nil {
		dialog.NewError(err, window).Show()
		isCriticalError = true
	}

	// create texts
	loginTitleText := utils.CreateText("Login to", color.White, 30, fyne.TextAlignCenter, fyne.TextStyle{Bold: true}, false)
	vaultNameTitleText := utils.CreateText("\""+vaultManager.Vaults[vaultIndex].Name+"\"", color.White, 30, fyne.TextAlignCenter, fyne.TextStyle{Bold: true}, false)
	loginFailedMessage := utils.CreateText("", color.RGBA{R: 255, G: 0, B: 0, A: 255}, 15, fyne.TextAlignCenter, fyne.TextStyle{Bold: true}, false)

	// create entries
	passwordEntry := widget.NewPasswordEntry()

	onLogin := func() {
		err := bcrypt.CompareHashAndPassword(credentialManager.GetEncryptionHash(), []byte(passwordEntry.Text))
		if err != nil {
			loginFailedMessage.Text = "Login failed!"
			return
		}

		err = credentialManager.SetEncryptionKey(passwordEntry.Text)
		if err != nil {
			dialog.NewError(err, window).Show()
			return
		}
		err = credentialManager.ReadFromFile(vaultManager.Vaults[vaultIndex].SaveLocation)
		if err != nil {
			dialog.NewError(err, window).Show()
			isCriticalError = true
			return
		}
		window.SetContent(mainView(app, window, settings, credentialManager, vaultManager, vaultIndex))
	}

	// create buttons
	loginButton := widget.NewButton("Login", func() {
		onLogin()
	})
	returnButton := widget.NewButtonWithIcon("Go back", theme.MailReplyIcon(), func() {
		window.SetContent(VaultView(app, window))
	})

	// Set button importance
	loginButton.Importance = widget.HighImportance

	passwordEntry.OnSubmitted = func(s string) {
		onLogin()
	}

	// disable components if save file doesnt exist
	if isCriticalError {
		passwordEntry.Disable()
		loginButton.Disable()
	}

	// container
	loginContainer := container.NewVBox(
		loginTitleText,
		vaultNameTitleText,
		passwordEntry,
		loginButton,
		returnButton,
		loginFailedMessage,
	)

	return container.New(
		customLayout.NewResponsivePaddingLayout(0.2),
		container.NewVBox(
			layout.NewSpacer(),
			loginContainer,
			layout.NewSpacer(),
		),
	)

}
