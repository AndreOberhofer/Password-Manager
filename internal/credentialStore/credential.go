package credentialStore

type jsonCredential struct {
	Id             string
	BranchId       string
	Title          string
	CredentialType CredentialType
	Email          string
	Username       string
	Password       string
}

type Credential struct {
	id             string
	branchId       string
	title          string
	credentialType CredentialType
	email          string
	username       string
	password       string
	unsavedChanges bool
}

func NewCredential(id string, branchId string, title string, credentialType CredentialType, email string, username string, password string) *Credential {
	return &Credential{
		id:             id,
		branchId:       branchId,
		title:          title,
		credentialType: credentialType,
		email:          email,
		username:       username,
		password:       password,
		unsavedChanges: false,
	}
}

func (credential Credential) GetId() string {
	return credential.id
}

func (credential Credential) GetBranchId() string {
	return credential.branchId
}

func (credential Credential) GetCredentialType() CredentialType {
	return credential.credentialType
}

func (credential *Credential) SetTitle(title string) {
	credential.title = title
	credential.setUnsavedChanges()
}

func (credential Credential) GetTitle() string {
	return credential.title
}

func (credential Credential) GetType() CredentialType {
	return credential.credentialType
}

func (credential *Credential) setUnsavedChanges() {
	credential.unsavedChanges = true
}

func (credential *Credential) resetUnsavedChanges() {
	credential.unsavedChanges = false
}

func (credential Credential) hasUnsavedChanges() bool {
	return credential.unsavedChanges
}

func (credential *Credential) SetPassword(password string) {
	credential.password = password
	credential.setUnsavedChanges()
}

func (credential Credential) GetPassword() string {
	return credential.password
}

func (credential *Credential) SetEmail(email string) {
	credential.email = email
	credential.setUnsavedChanges()
}

func (credential Credential) GetEmail() string {
	return credential.email
}

func (credential *Credential) SetUsername(username string) {
	credential.username = username
	credential.setUnsavedChanges()
}

func (credential Credential) GetUsername() string {
	return credential.username
}
