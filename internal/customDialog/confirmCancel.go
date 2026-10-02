package customDialog

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// ConfirmCancelDialog returns a dialog which can be confirmed, denied or canceled
func ConfirmCancelDialog(title string, message string, onConfirm func(), onDeny func(), window fyne.Window) *dialog.CustomDialog {
	var newDialog *dialog.CustomDialog

	dialogMessage := widget.NewLabelWithStyle(message, fyne.TextAlignCenter, fyne.TextStyle{})
	dialogMessage.Wrapping = fyne.TextWrapWord

	confirmButton := widget.NewButtonWithIcon("Yes", theme.ConfirmIcon(), func() {
		onConfirm()
		newDialog.Dismiss()
	})
	confirmButton.Importance = widget.HighImportance
	denyButton := widget.NewButtonWithIcon("No", theme.CancelIcon(), func() {
		onDeny()
		newDialog.Dismiss()
	})
	cancelButton := widget.NewButtonWithIcon("Cancel", theme.ContentRemoveIcon(), func() {
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
					3,
					confirmButton,
					denyButton,
					cancelButton,
				),
			),
		),
	)

	newDialog = dialog.NewCustomWithoutButtons(title, content, window)
	return newDialog
}
