package validator

import (
	"fmt"
	"io/fs"
	"regexp"

	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
)

func SetValidatorEntryText(entry *widget.Entry, errorText *canvas.Text) {
	entry.SetOnValidationChanged(func(err error) {
		if err != nil {
			errorText.Text = err.Error()
			errorText.Show()
		} else {
			errorText.Hide()
		}
	})
}

func EmptyStringValidator(s string) error {
	if s == "" {
		return fmt.Errorf("can't be empty")
	}
	return nil
}

func SaveLocationValidator(s string) error {
	if s == "" {
		return fmt.Errorf("can't be empty")
	} else if !fs.ValidPath(s) {
		return fmt.Errorf("not a valid directory path")
	}
	return nil
}

func BetterSaveLocationValidator(s string) error {
	//regex := `^[a-zA-Z0-9\/_\-\.]*.\.json$|^[A-Z]\:{1}\/{1}[a-zA-Z0-9\/_\-\.]*.\.json$`
	regex := `^[a-zA-Z0-9\/_\-\.]*.$|^[A-Z]\:{1}\/{1}[a-zA-Z0-9\/_\-\.]*.$`

	match, err := regexp.MatchString(regex, s)
	if err != nil {
		return err
	}

	if s == "" {
		return fmt.Errorf("can't be empty")
	} else if !match {
		return fmt.Errorf("isn't a valid path")
	}
	return nil
}
