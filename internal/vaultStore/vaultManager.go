package vaultStore

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"password_manager/internal/utils"
)

// vault file location
const vaultFileLocation = "vaults/vaults.json"

type VaultManager struct {
	Vaults []Vault
}

func NewVaultManager() *VaultManager {
	return &VaultManager{
		Vaults: []Vault{},
	}
}

func (vaultManager *VaultManager) ReadFromFile() error {
	// check if file exists
	if _, err := os.Stat(vaultFileLocation); errors.Is(err, os.ErrNotExist) {
		// Create file if it doesn't exist yet
		err := vaultManager.WriteToFile()
		if err != nil {
			return fmt.Errorf("couldn't create missing vaults save file: %v", err)
		}
	} else if err != nil {
		// Return other errors
		return err
	}

	// open json file
	jsonFile, err := os.Open(vaultFileLocation)
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	// read opened file as a byte array
	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return err
	}

	return json.Unmarshal(byteValue, &vaultManager.Vaults)
}

func (vaultManager *VaultManager) WriteToFile() error {
	err := utils.CheckAndCreateFile(vaultFileLocation)
	if err != nil {
		return err
	}

	byteData, err := json.Marshal(vaultManager.Vaults)
	if err != nil {
		return err
	}

	return os.WriteFile(vaultFileLocation, byteData, 0644)
}

func (vaultManager VaultManager) GetVaultNames() *[]string {
	var names []string
	for _, vault := range vaultManager.Vaults {
		names = append(names, vault.Name)
	}
	return &names
}
