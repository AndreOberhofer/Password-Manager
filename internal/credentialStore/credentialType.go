package credentialStore

import (
	"password_manager/internal/assets"

	"fyne.io/fyne/v2"
)

type CredentialType int

func (credentialType CredentialType) String() string {
	// Return nothing if the credential type is not in scope
	if int(credentialType) < 0 || int(credentialType) >= len(CredentialTypeNames) {
		return ""
	}
	return CredentialTypeNames[credentialType]
}

func (credentialType CredentialType) GetIcon() fyne.Resource {
	// Load standard icon if credential type is not in scope
	if int(credentialType) < 0 || int(credentialType) >= len(CredentialTypeIcons) {
		return assets.LoadIcon(assets.IconCredential)
	}
	return CredentialTypeIcons[credentialType]
}

const (
	Password CredentialType = iota
	Email
	User
)

var CredentialTypeNames = []string{
	"Password",
	"Email",
	"User",
}

var CredentialTypeIcons = []fyne.Resource{
	assets.LoadIcon(assets.IconPasswordCredential),
	assets.LoadIcon(assets.IconEmailCredential),
	assets.LoadIcon(assets.IconUserCredential),
}
