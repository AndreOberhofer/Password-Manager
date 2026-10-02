package credentialStore

import (
	"github.com/google/uuid"
)

type jsonBranch struct {
	Name     string            `json:"name"`
	Parent   string            `json:"parent"`
	Contents map[string]string `json:"contents"`
}

type jsonTree struct {
	Branches map[string]jsonBranch `json:"branches"`
}

func newJsonTree() jsonTree {
	return jsonTree{
		Branches: make(map[string]jsonBranch),
	}
}

type branch struct {
	name     string
	parent   string
	contents map[string]string
}

type tree struct {
	branches map[string]branch
}

func newTree() *tree {
	return &tree{
		branches: map[string]branch{
			"": branch{
				name:     "",
				parent:   "",
				contents: make(map[string]string),
			},
		},
	}
}

// IsTreeElementABranch checks if the ID is a branch
func (tree tree) IsTreeElementABranch(id string) bool {
	_, ok := tree.branches[id]
	return ok
}

func (tree *tree) addNewTreeBranch(newBranchName string, parentId string) { // TODO: make check if ID exists
	// Add new branch
	id := uuid.NewString()
	tree.branches[id] = branch{
		name:     newBranchName,
		parent:   parentId,
		contents: make(map[string]string),
	}

	// Add new branch to parent content
	tree.branches[parentId].contents[id] = id
}

func (tree *tree) removeTreeBranch(branchId string) { // TODO: make check if ID exists
	// Remove in parent branch
	parentId := tree.branches[branchId].parent
	delete(tree.branches[parentId].contents, branchId)

	// Remove branch
	delete(tree.branches, branchId)
}
