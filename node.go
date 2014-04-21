package bplustree

const (
	MaxKV = 255
	MaxKC = 511
)

type node interface {
	find(key int) (int, bool)
	parent() *interiorNode
	setParent(*interiorNode)
	full() bool
}
