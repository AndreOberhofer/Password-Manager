package customDialog

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// ConfirmDialog returns a dialog which can either be confirmed or denied
func ConfirmDialog(title string, message string, onConfirm func(), window fyne.Window) *dialog.CustomDialog {
	var newDialog *dialog.CustomDialog

	dialogMessage := widget.NewLabelWithStyle(message, fyne.TextAlignCenter, fyne.TextStyle{})
	dialogMessage.Wrapping = fyne.TextWrapWord

	confirmButton := widget.NewButtonWithIcon("Yes", theme.ConfirmIcon(), func() {
		onConfirm()
		newDialog.Dismiss()
	})
	confirmButton.Importance = widget.HighImportance
	denyButton := widget.NewButtonWithIcon("No", theme.CancelIcon(), func() {
		newDialog.Dismiss()
	})

	content := container.NewVBox(
		dialogMessage,
		container.New(
			layout.NewCenterLayout(),
			container.New(
				layout.NewCustomPaddedLayout(
					20,
					0,
					0,
					0,
				),
				container.NewGridWithColumns(
					2,
					confirmButton,
					denyButton,
				),
			),
		),
	)

	newDialog = dialog.NewCustomWithoutButtons(title, content, window)
	newDialog.Resize(fyne.NewSize(300, 0))
	return newDialog
}
