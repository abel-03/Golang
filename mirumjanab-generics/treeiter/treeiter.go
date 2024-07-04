//go:build !solution

package treeiter

func DoInOrder[T interface { Left() *T;	Right() *T }] (root *T, callback func(node *T)) {
	if root == nil {
		return
	}
	DoInOrder((*root).Left(), callback)
	callback(root)
	DoInOrder((*root).Right(), callback)
}
