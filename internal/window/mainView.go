package window

import (
	"context"
	"image/color"
	"password_manager/internal/assets"
	"password_manager/internal/credentialStore"
	"password_manager/internal/customDialog"
	"password_manager/internal/customLayout"
	"password_manager/internal/settingStore"
	"password_manager/internal/utils"
	"password_manager/internal/validator"
	"password_manager/internal/vaultStore"
	"password_manager/internal/widgetExtension"
	"sync"
	"sync/atomic"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var copyTime binding.Float = binding.NewFloat()

func mainView(app fyne.App, window fyne.Window, settings *settingStore.GeneralSettings, credentialManager *credentialStore.CredentialManager, vaultManager *vaultStore.VaultManager, vaultId int) *fyne.Container {
	var ctx context.Context
	var cancel context.CancelFunc
	var running atomic.Bool
	var waitGroup sync.WaitGroup
	var tree *widget.Tree
	var selectedCredentialId string
	var applyButton *widget.Button

	// Create editor background image
	lockImage := canvas.NewImageFromResource(assets.LoadBackground(assets.BackgroundLock))
	lockImage.SetMinSize(fyne.NewSize(200, 200))
	lockImage.Translucency = 0.9

	// Create source and destination rectangle for drag indicator
	sourceRectangle := canvas.NewRectangle(color.RGBA{0, 0, 0, 0})
	sourceRectangle.CornerRadius = 8
	destinationRectangle := canvas.NewRectangle(color.RGBA{39, 174, 245, 0})
	destinationRectangle.CornerRadius = 8

	// Create icon and title for credential editing container
	selectedCredentialIcon := canvas.NewImageFromResource(assets.LoadIcon(assets.IconCredential))
	selectedCredentialIcon.SetMinSize(fyne.NewSize(50, 50))
	selectedCredentialTitle := canvas.NewText("", color.White)
	selectedCredentialTitle.TextSize = 16

	// Create label
	vaultName := widget.NewLabel(vaultManager.Vaults[vaultId].Name)
	vaultName.TextStyle.Bold = true

	// Create copy bar
	copyTimeBar := widget.NewProgressBar()
	copyTimeBar.TextFormatter = func() string {
		return ""
	}
	copyTimeBar.Bind(copyTime)
	copyTimeBar.Min = float64(0)
	copyTimeBar.Max = float64(settings.CopyTime * 10)

	// Create entry field
	titleEntry := widget.NewEntry()
	usernameEntry := widget.NewEntry()
	emailEntry := widget.NewEntry()
	passwordEntry := widget.NewPasswordEntry()

	// Set place holder text
	titleEntry.PlaceHolder = "Enter title..."
	usernameEntry.PlaceHolder = "Enter username..."
	emailEntry.PlaceHolder = "Enter email..."
	passwordEntry.PlaceHolder = "Enter password..."

	// Handle entry validator
	titleEntry.Validator = validator.EmptyStringValidator

	// Create validator text
	titleEntryValidatorText := utils.CreateEntryErrorDisplayText()

	// Handle validator text
	validator.SetValidatorEntryText(titleEntry, titleEntryValidatorText)

	// Set place holder text to that the validator text is not immediately displayed
	titleEntry.SetText("PLACEHOLDER")

	// Handle on change events to enable/disable button
	disableButtonOnInvalidEntries := func(s string) {
		if titleEntry.Validate() != nil {
			applyButton.Disable()
			return
		}
		applyButton.Enable()
	}
	titleEntry.OnChanged = disableButtonOnInvalidEntries

	// Create tool bar
	treeToolBar := widget.NewToolbar(
		widget.NewToolbarAction(theme.ViewRefreshIcon(), func() {
			if credentialManager.AreUnsavedChanges() {
				customDialog.ConfirmDialog(
					"Reload data",
					"Are you sure you want to reload the data from the file?\nUnsaved changes will be lost.",
					func() {
						err := credentialManager.ReadFromFile(vaultManager.Vaults[vaultId].SaveLocation)
						if err != nil {
							dialog.NewError(err, window).Show()
							return
						}
						window.Content().Refresh()
					},
					window).Show()
			} else {
				err := credentialManager.ReadFromFile(vaultManager.Vaults[vaultId].SaveLocation)
				if err != nil {
					dialog.NewError(err, window).Show()
					return
				}
				window.Content().Refresh()
			}
		}),
		widget.NewToolbarAction(theme.FolderNewIcon(), func() {
			customDialog.AddFolderDialog("Add branch", "", credentialManager, window).Show()
		}),
		widget.NewToolbarAction(assets.LoadIcon(assets.IconAddCredential), func() {
			customDialog.AddCredentialDialog("Add credential", credentialManager, "", window).Show()
		}),
		widget.NewToolbarAction(theme.DocumentSaveIcon(), func() {
			err := credentialManager.WriteToFile(vaultManager.Vaults[vaultId].SaveLocation)
			if err != nil {
				dialog.NewError(err, window).Show()
				return
			}
		}),
		widget.NewToolbarAction(assets.LoadIcon(assets.IconExpandTree), func() {
			tree.OpenAllBranches()
		}),
		widget.NewToolbarAction(assets.LoadIcon(assets.IconCollapseTree), func() {
			tree.CloseAllBranches()
		}),
	)
	vaultToolBar := widget.NewToolbar(
		widget.NewToolbarAction(theme.LogoutIcon(), func() {
			if credentialManager.AreUnsavedChanges() {
				customDialog.ConfirmCancelDialog(
					"Unsaved changes!",
					"You have unsaved changes.\nDo you want to save before exiting?",
					func() {
						err := credentialManager.WriteToFile(vaultManager.Vaults[vaultId].SaveLocation)
						if err != nil {
							dialog.NewError(err, window).Show()
							return
						}
						window.SetContent(VaultView(app, window))
					},
					func() {
						credentialManager.ResetUnsavedChanges()
						window.SetContent(VaultView(app, window))
					},
					window,
				).Show()
			} else {
				credentialManager.ResetUnsavedChanges()
				window.SetContent(VaultView(app, window))
			}
		}),
		widget.NewToolbarAction(theme.SettingsIcon(), func() {
			window.SetContent(settingsView(app, window, settings, credentialManager, vaultManager, vaultId))
		}),
	)

	// Create tree
	tree = widget.NewTree(
		func(uid widget.TreeNodeID) []string {
			return credentialManager.GetTreeBranchContent(uid)
		},
		func(uid widget.TreeNodeID) bool {
			return credentialManager.Tree.IsTreeElementABranch(uid)
		},
		func(branch bool) fyne.CanvasObject {
			icon := widget.NewIcon(nil)
			treeNode := widgetExtension.NewTreeNode(
				branch,
				credentialManager,
				tree,
				destinationRectangle,
				sourceRectangle,
				func(uid string) {
					tree.Select(uid)

					if tree.IsBranch(uid) {
						tree.ToggleBranch(uid)
					}
				},
			)
			return container.New(
				layout.NewBorderLayout(
					nil,
					nil,
					icon,
					nil,
				),
				icon,
				treeNode,
			)
		},
		func(uid string, branch bool, obj fyne.CanvasObject) {
			container := obj.(*fyne.Container)

			// Get icon
			var resource fyne.Resource
			if branch {
				resource = theme.FolderIcon()
			} else {
				credential, err := credentialManager.GetCredential(uid)
				if err != nil {
					dialog.NewError(err, window).Show()
					return
				}
				resource = credential.GetType().GetIcon()
			}

			container.Objects[0].(*widget.Icon).SetResource(resource)
			container.Objects[1].(*widgetExtension.TreeNode).SetUID(uid)
			container.Objects[1].(*widgetExtension.TreeNode).SetText(credentialManager.GetBranchOrCredentialTitle(uid))
		},
	)

	// Create button
	copyButton := widget.NewButton("Copy", func() {
		app.Clipboard().SetContent(passwordEntry.Text)

		if running.Load() {
			cancel() // stop old goroutine
		} else {
			running.Store(true)
		}

		waitGroup.Wait()
		ctx, cancel = context.WithCancel(context.Background())

		// will empty the clipboard after 10 seconds
		waitGroup.Add(1)
		go func() {
			defer waitGroup.Done()

			for i := (settings.CopyTime * 10); i >= 0; i-- {
				select {
				case <-ctx.Done():
					return
				default:
					err := copyTime.Set(float64(i))
					if err != nil {
						dialog.NewError(err, window).Show()
					}
					time.Sleep(time.Duration(100 * time.Millisecond))
				}
			}
			app.Clipboard().SetContent("")
		}()
	})
	applyButton = widget.NewButton("Apply", func() {
		credential, err := credentialManager.GetCredential(selectedCredentialId)
		if err != nil {
			dialog.NewError(err, window).Show()
			return
		}

		credential.SetTitle(titleEntry.Text)
		credential.SetPassword(passwordEntry.Text)
		credential.SetEmail(emailEntry.Text)
		credential.SetUsername(usernameEntry.Text)
	})

	// Hide entries
	hideAndEmptyAllEntries := func() {
		// Jump out if the credential title is still hidden (everything else will be hidden too then)
		if selectedCredentialTitle.Hidden {
			return
		}

		// Hide icon and title
		selectedCredentialIcon.Hide()
		selectedCredentialTitle.Hide()

		// Hide entries
		titleEntry.Hide()
		titleEntryValidatorText.Hide()
		usernameEntry.Hide()
		emailEntry.Hide()
		passwordEntry.Hide()
		copyButton.Hide()

		// Hide buttons
		applyButton.Hide()

		// I don't use the function when settings the text so the "OnChanged" event isn't called
		// Because that would write to the data
		titleEntry.Text = ""
		usernameEntry.Text = ""
		emailEntry.Text = ""
		passwordEntry.Text = ""

		// Display lock image
		lockImage.Show()
	}
	hideAndEmptyAllEntries()

	tree.OnSelected = func(uid widget.TreeNodeID) {
		// Hide and empty all fields
		hideAndEmptyAllEntries()

		// If it is a branch jump out
		if tree.IsBranch(uid) {
			return
		}

		// Get credential
		selectedCredentialId = uid
		credential, err := credentialManager.GetCredential(selectedCredentialId)
		if err != nil {
			dialog.NewError(err, window).Show()
			return
		}

		// Hide lock image
		lockImage.Hide()

		// Show and fill base fields
		selectedCredentialIcon.Show()
		selectedCredentialIcon.Resource = credential.GetType().GetIcon()
		selectedCredentialIcon.Refresh()
		selectedCredentialTitle.Show()
		selectedCredentialTitle.Text = credential.GetTitle()
		selectedCredentialTitle.Refresh()
		titleEntry.Show()
		titleEntry.SetText(credential.GetTitle())
		titleEntry.Refresh()

		// Show and fill special fields
		switch credential.GetType() {
		case credentialStore.Password:
			passwordEntry.Show()
			passwordEntry.Text = credential.GetPassword()
			passwordEntry.Refresh()
			copyButton.Show()
		case credentialStore.Email:
			emailEntry.Show()
			emailEntry.Text = credential.GetEmail()
			emailEntry.Refresh()
			passwordEntry.Show()
			passwordEntry.Text = credential.GetPassword()
			passwordEntry.Refresh()
			copyButton.Show()
		case credentialStore.User:
			usernameEntry.Show()
			usernameEntry.Text = credential.GetUsername()
			usernameEntry.Refresh()
			passwordEntry.Show()
			passwordEntry.Text = credential.GetPassword()
			passwordEntry.Refresh()
			copyButton.Show()
		}

		// Show buttons
		applyButton.Show()
	}

	// Handle window event
	window.SetCloseIntercept(func() {
		if credentialManager.AreUnsavedChanges() {
			customDialog.ConfirmCancelDialog(
				"Unsaved changes!",
				"You have unsaved changes.\nDo you want to save before exiting?",
				func() {
					err := credentialManager.WriteToFile(vaultManager.Vaults[vaultId].SaveLocation)
					if err != nil {
						dialog.NewError(err, window).Show()
						return
					}
					app.Quit()
				},
				func() {
					app.Quit()
				},
				window,
			).Show()
		} else {
			app.Quit()
		}
	})

	// Create container
	vaultToolbarContainer := container.NewHBox(
		vaultName,
		vaultToolBar,
	)
	vaulToolBarContainer := container.New(
		layout.NewBorderLayout(
			nil,
			nil,
			nil,
			vaultToolbarContainer,
		),
		vaultToolbarContainer,
	)
	copyTimeBarContainer := container.New(
		layout.NewCustomPaddedLayout(
			0,
			10,
			10,
			10,
		),
		copyTimeBar,
	)
	titleEntryContainer := container.New(
		layout.NewBorderLayout(
			nil,
			titleEntryValidatorText,
			nil,
			nil,
		),
		titleEntry,
		titleEntryValidatorText,
	)
	passwordEntryContainer := container.New(
		layout.NewBorderLayout(nil, nil, nil, copyButton),
		passwordEntry,
		copyButton,
	)
	credentialEditorButtonContainer := container.New(
		layout.NewBorderLayout(
			nil,
			nil,
			applyButton,
			nil,
		),
		applyButton,
	)
	credentialEditorContainer := container.NewStack(
		container.New(
			customLayout.NewResponsivePaddingLayout(0.1),
			container.NewVBox(
				container.NewCenter(
					container.NewHBox(
						selectedCredentialIcon,
						selectedCredentialTitle,
					),
				),
				titleEntryContainer,
				usernameEntry,
				emailEntry,
				passwordEntryContainer,
				credentialEditorButtonContainer,
			),
		),
		container.NewCenter(
			lockImage,
		),
	)
	credentialEditorContainerWithBar := container.New(
		layout.NewBorderLayout(
			nil,
			copyTimeBarContainer,
			nil,
			nil,
		),
		copyTimeBarContainer,
		credentialEditorContainer,
	)
	listWithToolbarContainer := container.New(
		layout.NewBorderLayout(
			treeToolBar,
			nil,
			nil,
			nil,
		),
		treeToolBar,
		tree,
	)
	listAndCredentialEditorContainer := container.NewHSplit(
		listWithToolbarContainer,
		credentialEditorContainerWithBar,
	)
	listAndCredentialEditorContainer.SetOffset(0.3)

	return container.NewStack(
		container.New(
			layout.NewBorderLayout(
				vaulToolBarContainer,
				nil,
				nil,
				nil,
			),
			vaulToolBarContainer,
			listAndCredentialEditorContainer,
		),
		destinationRectangle,
		sourceRectangle,
	)
}
