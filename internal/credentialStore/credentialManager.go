package credentialStore

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"password_manager/internal/encryption"
	"password_manager/internal/utils"
	"sort"
)

type jsonCredentialManager struct {
	EncryptionHash string                    `json:"encryptionHash"`
	Credentials    map[string]jsonCredential `json:"credentials"`
	Tree           jsonTree                  `json:"tree"`
}

func newJsonCredentialManager() jsonCredentialManager {
	return jsonCredentialManager{
		EncryptionHash: "",
		Credentials:    make(map[string]jsonCredential),
		Tree:           newJsonTree(),
	}
}

type CredentialManager struct {
	encryptionHash []byte
	encryptionKey  []byte
	credentials    map[string]*Credential
	Tree           *tree
	unsavedChanges bool
}

func NewCredentialManager() *CredentialManager {
	return &CredentialManager{
		encryptionHash: []byte{},
		encryptionKey:  []byte{},
		credentials:    make(map[string]*Credential),
		Tree:           newTree(),
		unsavedChanges: false,
	}
}

func (credentialManager *CredentialManager) ReadEncryptionHashFromFile(fileLocation string) error {
	// open json file
	jsonFile, err := os.Open(fileLocation)
	if err != nil {
		return fmt.Errorf("error when trying to open file:\n%v", err)
	}
	defer jsonFile.Close()

	// read opened file as a byte array
	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return fmt.Errorf("error when trying to read from file:\n%v", err)
	}

	// unmarshal json file to structs
	var jsonEncryptionHash struct {
		EncryptionHash string `json:"encryptionHash"`
	}
	err = json.Unmarshal(byteValue, &jsonEncryptionHash)
	if err != nil {
		return fmt.Errorf("error when trying to unmarshal from json format:\n%v", err)
	}

	// Map encryption hash
	credentialManager.encryptionHash = []byte(jsonEncryptionHash.EncryptionHash)
	return nil
}

// read from json file and fill into slice
func (credentialManager *CredentialManager) ReadFromFile(fileLocation string) error {
	// check if file exists
	if _, err := os.Stat(fileLocation); errors.Is(err, os.ErrNotExist) {
		// Create file if it doesn't exist yet
		err := credentialManager.WriteToFile(fileLocation)
		if err != nil {
			return fmt.Errorf("couldn't create missing credentials save file: %v", err)
		}
	} else if err != nil {
		// Return other errors
		return err
	}

	// open json file
	jsonFile, err := os.Open(fileLocation)
	if err != nil {
		return fmt.Errorf("error when trying to open file:\n%v", err)
	}
	defer jsonFile.Close()

	// read opened file as a byte array
	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return fmt.Errorf("error when trying to read from file:\n%v", err)
	}

	// unmarshal json file to structs
	jsonCredentialManager := newJsonCredentialManager()
	err = json.Unmarshal(byteValue, &jsonCredentialManager)
	if err != nil {
		return fmt.Errorf("error when trying to unmarshal from json format:\n%v", err)
	}

	// Map tree
	for key, jsonBranch := range jsonCredentialManager.Tree.Branches {
		credentialManager.Tree.branches[key] = branch{
			name:     jsonBranch.Name,
			parent:   jsonBranch.Parent,
			contents: jsonBranch.Contents,
		}
	}

	// Map credentials
	for key, jsonCredential := range jsonCredentialManager.Credentials {
		credentialManager.credentials[key] = &Credential{
			id:             jsonCredential.Id,
			branchId:       jsonCredential.BranchId,
			title:          jsonCredential.Title,
			credentialType: jsonCredential.CredentialType,
			email:          jsonCredential.Email,
			username:       jsonCredential.Username,
			password:       jsonCredential.Password,
			unsavedChanges: false,
		}
	}

	// Map encryption hash
	credentialManager.encryptionHash = []byte(jsonCredentialManager.EncryptionHash)

	credentialManager.ResetUnsavedChanges()
	return nil
}

// write to json file
func (credentialManager *CredentialManager) WriteToFile(fileLocation string) error {
	err := utils.CheckAndCreateFile(fileLocation)
	if err != nil {
		return err
	}

	// map unexported internal structs to exported json structs
	jsonCredentialManager := newJsonCredentialManager()

	// Map credentials
	for key, credential := range credentialManager.credentials {
		jsonCredentialManager.Credentials[key] = jsonCredential{
			Id:             credential.id,
			BranchId:       credential.branchId,
			Title:          credential.title,
			CredentialType: credential.credentialType,
			Email:          credential.email,
			Username:       credential.username,
			Password:       credential.password,
		}
	}

	// Map encryption hash
	jsonCredentialManager.EncryptionHash = string(credentialManager.encryptionHash)

	// Map tree
	for key, branch := range credentialManager.Tree.branches {
		jsonCredentialManager.Tree.Branches[key] = jsonBranch{
			Name:     branch.name,
			Parent:   branch.parent,
			Contents: branch.contents,
		}
	}

	// marshal structs to json format
	byteData, err := json.Marshal(jsonCredentialManager)
	if err != nil {
		return fmt.Errorf("error when trying to marshal to json format: %v", err)
	}

	// write json data to file
	err = os.WriteFile(fileLocation, byteData, 0644)
	if err != nil {
		return fmt.Errorf("error when writing to file: %v", err)
	}

	credentialManager.ResetUnsavedChanges()
	return nil
}

// get encryption hash
func (credentialManager CredentialManager) GetEncryptionHash() []byte {
	return credentialManager.encryptionHash
}

// get encryption key
func (credentialManager CredentialManager) GetEncryptionKey() string {
	return string(credentialManager.encryptionKey)
}

// set encryption key
func (credentialManager *CredentialManager) SetEncryptionKey(encryptionKey string) error {
	// Set encryption key
	credentialManager.encryptionKey = []byte(encryptionKey)

	// Set encryption hash
	encryptionHash, err := encryption.Hash([]byte(encryptionKey))
	if err != nil {
		return err
	}
	credentialManager.encryptionHash = encryptionHash
	credentialManager.setUnsavedChanges()
	return nil
}

// get if there are unsaved changes
func (credentialManager CredentialManager) AreUnsavedChanges() bool {
	for _, credential := range credentialManager.credentials {
		if credential.hasUnsavedChanges() {
			return true
		}
	}
	return credentialManager.unsavedChanges
}

// set unsaved changes
func (credentialManager *CredentialManager) setUnsavedChanges() {
	credentialManager.unsavedChanges = true
}

// reset unsaved changes
func (credentialManager *CredentialManager) ResetUnsavedChanges() {
	for _, credential := range credentialManager.credentials {
		credential.resetUnsavedChanges()
	}
	credentialManager.unsavedChanges = false
}

// AddCredential adds new credential
func (credentialManager *CredentialManager) AddCredential(credential *Credential) error {
	// Add credential
	_, exists := credentialManager.credentials[credential.id]
	if exists {
		return fmt.Errorf("credential with this ID already exists")
	}
	credentialManager.credentials[credential.id] = credential

	// Add credential to tree branches
	_, exists = credentialManager.Tree.branches[credential.branchId].contents[credential.id]
	if exists {
		return fmt.Errorf("credential with this ID already exists in the tree")
	}
	credentialManager.Tree.branches[credential.branchId].contents[credential.id] = credential.id

	credentialManager.setUnsavedChanges()
	return nil
}

// AddCredential adds new credential
func (credentialManager *CredentialManager) RemoveCredential(credentialId string) error {
	// Remove credential from tree branch
	_, exists := credentialManager.Tree.branches[credentialManager.credentials[credentialId].branchId].contents[credentialId]
	if !exists {
		return fmt.Errorf("credential doesn't exist in tree branch")
	}
	delete(credentialManager.Tree.branches[credentialManager.credentials[credentialId].branchId].contents, credentialId)

	// Remove credential
	_, exists = credentialManager.credentials[credentialId]
	if !exists {
		return fmt.Errorf("credential doesn't exist")
	}
	delete(credentialManager.credentials, credentialId)

	credentialManager.setUnsavedChanges()
	return nil
}

func (credentialManager *CredentialManager) GetCredential(credentialId string) (*Credential, error) {
	// Check if credential type exists
	credential, exists := credentialManager.credentials[credentialId]
	if !exists {
		return &Credential{}, fmt.Errorf("couldn't find credentials with ID %v", credentialId)
	}

	return credential, nil
}

// AddNewTreeBranch adds a new branch
func (credentialManager *CredentialManager) AddNewTreeBranch(newBranchName string, parentBranchId string) {
	credentialManager.Tree.addNewTreeBranch(newBranchName, parentBranchId)
	credentialManager.setUnsavedChanges()
}

// RemoveTreeBranch removes tree branch
func (credentialManager *CredentialManager) RemoveTreeBranch(branchId string) {
	// Remove all of the content in the branch
	for _, contentId := range credentialManager.Tree.branches[branchId].contents {
		_, isBranch := credentialManager.Tree.branches[contentId]
		if isBranch {
			credentialManager.RemoveTreeBranch(contentId)
		} else {
			credentialManager.RemoveCredential(contentId)
		}
	}

	// Remove the branch itself
	credentialManager.Tree.removeTreeBranch(branchId)
	credentialManager.setUnsavedChanges()
}

func (credentialManager CredentialManager) GetBranchOrCredentialTitle(id string) string {
	if _, isBranch := credentialManager.Tree.branches[id]; isBranch {
		return credentialManager.Tree.branches[id].name
	}

	if _, isCredential := credentialManager.credentials[id]; isCredential {
		return credentialManager.credentials[id].title
	}

	return ""
}

// GetTreeBranchContent returns the content of the branch
func (credentialManager CredentialManager) GetTreeBranchContent(branchId string) []string {
	// Check if branch exists
	branch, exists := credentialManager.Tree.branches[branchId]
	if !exists {
		return []string{}
	}

	// Get content
	content := make([]string, 0, len(branch.contents))
	for contentElement := range branch.contents {
		content = append(content, contentElement)
	}

	// Sort depending on title/name
	sort.Slice(content, func(i, j int) bool {
		nameFirst := credentialManager.GetBranchOrCredentialTitle(content[i])
		nameSecond := credentialManager.GetBranchOrCredentialTitle(content[j])

		if nameFirst == nameSecond {
			return content[i] < content[j]
		}

		return nameFirst < nameSecond
	})

	return content
}

// MoveCredential moves credential into other branch
func (credentialManager *CredentialManager) MoveCredential(oldBranchId string, newBranchId string, credentialId string) error {
	// Remove credential from old branch
	_, exists := credentialManager.Tree.branches[oldBranchId].contents[credentialId]
	if !exists {
		return fmt.Errorf("credential doesn't exist in tree branch")
	}
	delete(credentialManager.Tree.branches[oldBranchId].contents, credentialId)

	// Add credential to new tree branch
	_, exists = credentialManager.Tree.branches[newBranchId].contents[credentialId]
	if exists {
		return fmt.Errorf("credential with this ID already exists in the tree")
	}
	credentialManager.Tree.branches[newBranchId].contents[credentialId] = credentialId

	// Update credential branch ID
	credentialManager.credentials[credentialId].branchId = newBranchId

	credentialManager.setUnsavedChanges()
	return nil
}

// IsCredential checks if uid is in the credential map
func (credentialManager CredentialManager) IsCredential(uid string) bool {
	_, exists := credentialManager.credentials[uid]
	return exists
}

// IsBranch checks if uid is in the branch map
func (credentialManager CredentialManager) IsBranch(uid string) bool {
	_, exists := credentialManager.Tree.branches[uid]
	return exists
}
