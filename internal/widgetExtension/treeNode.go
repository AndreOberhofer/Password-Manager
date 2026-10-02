package widgetExtension

import (
	"image/color"
	"password_manager/internal/assets"
	"password_manager/internal/credentialStore"
	"password_manager/internal/customDialog"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var uidBoxes = make(map[string]box)

type box struct {
	x      float32
	y      float32
	width  float32
	height float32
}

func (box *box) IsOverlapping(x float32, y float32) bool {
	if x >= box.x &&
		x <= box.x+box.width &&
		y >= box.y &&
		y <= box.y+box.height {
		return true
	}
	return false
}

type TreeNode struct {
	widget.Label
	uid                  string
	isBranch             bool
	credentialManager    *credentialStore.CredentialManager
	tree                 *widget.Tree
	destinationRectangle *canvas.Rectangle
	sourceRectangle      *canvas.Rectangle
	currentDragUid       string
	onSelect             func(string)
}

func NewTreeNode(isBranch bool, credentialManager *credentialStore.CredentialManager, tree *widget.Tree, destinationRectangle *canvas.Rectangle, sourceRectangle *canvas.Rectangle, onSelect func(string)) *TreeNode {
	treeNode := &TreeNode{
		isBranch:             isBranch,
		credentialManager:    credentialManager,
		tree:                 tree,
		destinationRectangle: destinationRectangle,
		sourceRectangle:      sourceRectangle,
		onSelect:             onSelect,
	}
	treeNode.ExtendBaseWidget(treeNode)
	return treeNode
}

func (n *TreeNode) SetUID(uid string) {
	n.uid = uid
}

func (n *TreeNode) updateBox() {
	if n.isBranch {
		pos := fyne.CurrentApp().Driver().AbsolutePositionForObject(n)
		uidBoxes[n.uid] = box{
			x:      pos.X,
			y:      pos.Y,
			width:  n.Size().Width,
			height: n.Size().Height,
		}
	}
}

func (n *TreeNode) Move(pos fyne.Position) {
	n.Label.Move(pos)
	n.updateBox()
}

func (n *TreeNode) Resize(size fyne.Size) {
	n.Label.Resize(size)
	n.updateBox()
}

func (n *TreeNode) Tapped(_ *fyne.PointEvent) {
	if n.onSelect != nil {
		n.onSelect(n.uid)
	}
}

func (n *TreeNode) Dragged(ev *fyne.DragEvent) {
	if _, _, _, a := n.sourceRectangle.FillColor.RGBA(); a == 0 {
		pos := fyne.CurrentApp().Driver().AbsolutePositionForObject(n)
		n.sourceRectangle.Resize(fyne.NewSize(n.Size().Width, n.Size().Height+1))
		n.sourceRectangle.Move(fyne.NewPos(pos.X-5, pos.Y-5))
		n.sourceRectangle.FillColor = color.RGBA{0, 0, 0, 99}
	}

	for uid, box := range uidBoxes {
		if box.IsOverlapping(ev.AbsolutePosition.X, ev.AbsolutePosition.Y) {
			n.currentDragUid = uid
			n.destinationRectangle.Resize(fyne.NewSize(box.width, box.height+1))
			n.destinationRectangle.Move(fyne.NewPos(box.x-5, box.y-5))

			if _, _, _, a := n.destinationRectangle.FillColor.RGBA(); a == 0 {
				n.destinationRectangle.FillColor = color.RGBA{39, 174, 245, 99}
			}
		}
	}
}

func (n *TreeNode) DragEnd() {
	n.destinationRectangle.FillColor = color.RGBA{39, 174, 245, 0}
	n.sourceRectangle.FillColor = color.RGBA{0, 0, 0, 0}

	if n.credentialManager.IsCredential(n.uid) {
		window := fyne.CurrentApp().Driver().AllWindows()[0]
		credential, err := n.credentialManager.GetCredential(n.uid)
		if err != nil {
			dialog.NewError(err, window).Show()
			return
		}

		customDialog.ConfirmDialog(
			"Move credential",
			"Are you sure you want to move the following credential from '"+n.credentialManager.GetBranchOrCredentialTitle(credential.GetBranchId())+"' to '"+n.credentialManager.GetBranchOrCredentialTitle(n.currentDragUid)+"'",
			func() {
				n.credentialManager.MoveCredential(credential.GetBranchId(), n.currentDragUid, credential.GetId())
				n.tree.Refresh()

				// Closing and reopening the branch to fix the indenting
				n.tree.CloseBranch(credential.GetBranchId())
				n.tree.OpenBranch(credential.GetBranchId())
			},
			window,
		).Show()
	}
}

func (n *TreeNode) MouseDown(me *desktop.MouseEvent) {
	if me.Button == desktop.MouseButtonSecondary {
		// Get the canvas from the object
		canvas := fyne.CurrentApp().
			Driver().
			CanvasForObject(n)
		window := fyne.CurrentApp().Driver().AllWindows()[0]

		if n.isBranch {
			// Create the popup menu for branches
			addBranchMenuItem := fyne.NewMenuItemWithIcon("Add branch", theme.FolderNewIcon(), func() {
				customDialog.AddFolderDialog("Add folder", n.uid, n.credentialManager, window).Show()
			})
			addCredentialMenuItem := fyne.NewMenuItemWithIcon("Add credential", assets.LoadIcon(assets.IconAddCredential), func() {
				customDialog.AddCredentialDialog("Add credential", n.credentialManager, n.uid, window).Show()
			})
			removeBranchMenuItem := fyne.NewMenuItemWithIcon("Remove branch", theme.DeleteIcon(), func() {
				customDialog.ConfirmDialog(
					"Delete branch",
					"Are you sure you want to delete the following branch:\n"+n.credentialManager.GetBranchOrCredentialTitle(n.uid)+"\nAll of the data in the branch will be deleted too.",
					func() {
						n.credentialManager.RemoveTreeBranch(n.uid)
						window.Content().Refresh()
					},
					window).Show()
			})
			menu := fyne.NewMenu("Tree Item Menu",
				addBranchMenuItem,
				addCredentialMenuItem,
				removeBranchMenuItem,
			)
			popupMenu := widget.NewPopUpMenu(menu, canvas)
			popupMenu.Move(fyne.NewPos(me.AbsolutePosition.X, me.AbsolutePosition.Y))
			popupMenu.Show()
		} else {
			// Create the popup menu for credentials
			removeCredentialMenuItem := fyne.NewMenuItemWithIcon("Remove credential", theme.DeleteIcon(), func() {
				customDialog.ConfirmDialog(
					"Delete credential",
					"Are you sure you want to delete the following credential:\n"+n.credentialManager.GetBranchOrCredentialTitle(n.uid),
					func() {
						n.credentialManager.RemoveCredential(n.uid)
						window.Content().Refresh()
					},
					window).Show()
			})
			menu := fyne.NewMenu("Tree Item Menu",
				removeCredentialMenuItem,
			)
			popupMenu := widget.NewPopUpMenu(menu, canvas)
			popupMenu.Move(fyne.NewPos(me.AbsolutePosition.X, me.AbsolutePosition.Y))
			popupMenu.Show()
		}
	}
}

func (n *TreeNode) MouseUp(me *desktop.MouseEvent) {}
