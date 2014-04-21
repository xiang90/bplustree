package bplustree

type BTree struct {
	root     *interiorNode
	first    *leafNode
	leaf     int
	interior int
	height   int
}

func newBTree() *BTree {
	leaf := newLeafNode(nil)
	r := newInteriorNode(nil, leaf)
	leaf.p = r
	return &BTree{
		root:     r,
		first:    leaf,
		leaf:     1,
		interior: 1,
		height:   2,
	}
}

// first returns the first leafNode
func (bt *BTree) First() *leafNode {
	return bt.first
}

// insert inserts a (key, value) into the B+ tree
func (bt *BTree) Insert(key int, value string) {
	_, oldIndex, leaf := search(bt.root, key)
	p := leaf.parent()

	mid, bump := leaf.insert(key, value)
	if !bump {
		return
	}

	var midNode node
	midNode = leaf

	p.kcs[oldIndex].child = leaf.next
	leaf.next.setParent(p)

	interior, interiorP := p, p.parent()

	for {
		var oldIndex int
		var newNode *interiorNode

		isRoot := interiorP == nil

		if !isRoot {
			oldIndex, _ = interiorP.find(key)
		}

		mid, newNode, bump = interior.insert(mid, midNode)
		if !bump {
			return
		}

		if !isRoot {
			interiorP.kcs[oldIndex].child = newNode
			newNode.setParent(interiorP)

			midNode = interior
		} else {
			bt.root = newInteriorNode(nil, newNode)
			newNode.setParent(bt.root)

			bt.root.insert(mid, interior)
			return
		}

		interior, interiorP = interiorP, interior.parent()
	}
}

// Search searches the key in B+ tree
// If the key exists, it returns the value of key and true
// If the key does not exist, it returns an empty string and false
func (bt *BTree) Search(key int) (string, bool) {
	kv, _, _ := search(bt.root, key)
	if kv == nil {
		return "", false
	}
	return kv.value, true
}

func search(n node, key int) (*kv, int, *leafNode) {
	curr := n
	oldIndex := -1

	for {
		switch t := curr.(type) {
		case *leafNode:
			i, ok := t.find(key)
			if !ok {
				return nil, oldIndex, t
			}
			return &t.kvs[i], oldIndex, t
		case *interiorNode:
			i, _ := t.find(key)
			curr = t.kcs[i].child
			oldIndex = i
		default:
			panic("")
		}
	}
}
