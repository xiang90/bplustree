package bplustree

func search(n node, key int) (*kv, *leafNode) {
	i, ok := n.find(key)
	if !ok {
		return nil, n.(*leafNode)
	}

	switch t := n.(type) {
	case *leafNode:
		return &t.kvs[i], t
	case *interiorNode:
		return search(t.kcs[i].child, key)
	default:
		panic("wrong type")
	}
}
