package utils

import (
	"errors"
	"image/color"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
)

// TODO: change
func CreateText(text string, textColor color.Color, textSize float32, textAlignment fyne.TextAlign, textStyle fyne.TextStyle, isHidden bool) *canvas.Text {
	title := canvas.NewText(text, textColor)
	title.TextSize = textSize
	title.Alignment = textAlignment
	title.TextStyle = textStyle
	title.Hidden = isHidden
	return title
}

// TODO: remove if this stays unused
func VerticalSpacer(px float32) fyne.CanvasObject {
	rectangle := canvas.NewRectangle(color.Transparent)
	rectangle.Resize(fyne.NewSize(1, px))
	return rectangle
}

// TODO: remove if this stays unused
func showToast(canvas fyne.Canvas, message string) {
	// Create a label with the message you want to display
	toast := widget.NewLabel(message)

	// Create a Popup to display the toast
	toastPopup := widget.NewPopUp(toast, canvas)
	toastPopup.ShowAtPosition(fyne.NewPos(canvas.Size().Width-toastPopup.MinSize().Width-50, canvas.Size().Height-toastPopup.MinSize().Height-50)) // Position the popup

	// Set the toast to automatically disappear after 3 seconds
	go func() {
		fyne.Do(func() {
			time.Sleep(3 * time.Second)
			toastPopup.Hide() // Hide the toast after the delay
		})
	}()

	// Show the toast popup
	toastPopup.Show()
}

func GetLengthAsIndicesSlice(length int, combine string) []string {
	indices := make([]string, length)
	for i := range length {
		indices[i] = combine + "." + strconv.Itoa(i)
	}
	return indices
}

func CheckAndCreateFolder(directoryPath string) error {
	// Check if folder exists
	if _, err := os.Stat(directoryPath); errors.Is(err, os.ErrNotExist) {
		// Create parent directories if needed
		if err := os.MkdirAll(directoryPath, 0755); err != nil {
			return err
		}
	} else if err != nil {
		// Return other errors
		return err
	}
	return nil
}

func CheckAndCreateFile(filePath string) error {
	// check if file exists
	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		// Create parent directories if needed
		directoryPath := filepath.Dir(filePath)
		if err := os.MkdirAll(directoryPath, 0755); err != nil {
			return err
		}

		// Create file if it doesn't exist yet
		file, err := os.Create(filePath)
		if err != nil {
			return err
		}
		defer file.Close()
	} else if err != nil {
		// Return other errors
		return err
	}
	return nil
}

func CreateEntryErrorDisplayText() *canvas.Text {
	nameEntryValidatorText := canvas.NewText("", color.RGBA{R: 255, G: 0, B: 0, A: 255})
	nameEntryValidatorText.TextSize = 10
	return nameEntryValidatorText
}
