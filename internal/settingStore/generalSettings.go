package settingStore

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"password_manager/internal/utils"
)

// settings file location
const settingsFileLocation = "settings/settings.json"

type GeneralSettings struct {
	CopyTime int `json:"copyTime"`
}

// creates a standard instance of the object
func NewSettings() *GeneralSettings {
	return &GeneralSettings{
		CopyTime: 10,
	}
}

func (generalSettings *GeneralSettings) ReadFromFile() error {
	// check if file exists
	if _, err := os.Stat(settingsFileLocation); errors.Is(err, os.ErrNotExist) {
		// write settings file if it doesn't exist
		err := generalSettings.WriteToFile()
		if err != nil {
			return fmt.Errorf("couldn't create missing settings save file: %v", err)
		}
	} else if err != nil {
		// Return other errors
		return err
	}

	// open json file
	jsonFile, err := os.Open(settingsFileLocation)
	if err != nil {
		return err
	}

	defer jsonFile.Close()

	// read opened file as a byte array
	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return err
	}

	return json.Unmarshal(byteValue, &generalSettings)
}

func (generalSettings GeneralSettings) WriteToFile() error {
	err := utils.CheckAndCreateFile(settingsFileLocation)
	if err != nil {
		return err
	}

	byteData, err := json.Marshal(generalSettings)
	if err != nil {
		return err
	}

	return os.WriteFile(settingsFileLocation, byteData, 0644)
}
