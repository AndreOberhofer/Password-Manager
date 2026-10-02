package main

import (
	"password_manager/internal/window"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func main() {
	app := app.NewWithID("com.password-manager")
	mainWindow := app.NewWindow("Main")

	content := window.VaultView(app, mainWindow)

	mainWindow.SetContent(content)
	mainWindow.SetMaster()
	mainWindow.Show()
	mainWindow.Resize(fyne.NewSize(500, 300))
	app.Run()
}
