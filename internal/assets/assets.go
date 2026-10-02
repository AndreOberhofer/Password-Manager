package assets

import (
	"embed"

	"fyne.io/fyne/v2"
)

//go:embed icons/*.png
var iconFS embed.FS

//go:embed background/*.png
var backgroundFS embed.FS

const (
	IconCredential         string = "icons/credential.png"
	IconAddCredential      string = "icons/addCredential.png"
	IconUserCredential     string = "icons/userCredential.png"
	IconEmailCredential    string = "icons/emailCredential.png"
	IconPasswordCredential string = "icons/passwordCredential.png"
	IconExpandTree         string = "icons/expandTree.png"
	IconCollapseTree       string = "icons/collapseTree.png"
	BackgroundLock         string = "background/lock.png"
)

func LoadIcon(path string) *fyne.StaticResource {
	data, err := iconFS.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return fyne.NewStaticResource(path, data)
}

func LoadBackground(path string) *fyne.StaticResource {
	data, err := backgroundFS.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return fyne.NewStaticResource(path, data)
}
