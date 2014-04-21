package bplustree

import (
	"fmt"
	"testing"
	"time"
)

func TestInsert(t *testing.T) {
	testCount := 1000000
	bt := newBTree()

	start := time.Now()
	for i := 1; i <= testCount; i++ {
		bt.insert(i, "")
	}
	fmt.Println(time.Now().Sub(start))

	verifyTree(bt, testCount, t)
}

func verifyTree(b *BTree, count int, t *testing.T) {
	verifyRoot(b, t)

	for i := 0; i < b.root.count; i++ {
		verifyNode(b.root.kcs[i].child, b.root, t)
	}

	leftMost := findLeftMost(b.root)
	verifyLeaf(leftMost, count, t)
}

// min child: 1
// max child: MaxKC
func verifyRoot(b *BTree, t *testing.T) {
	if b.root.parent() != nil {
		t.Errorf("root.parent: want = nil, got = %p", b.root.parent())
	}

	if b.root.count < 1 {
		t.Errorf("root.min.child: want >=1, got = %d", b.root.count)
	}

	if b.root.count > MaxKC {
		t.Errorf("root.max.child: want <= %d, got = %d", MaxKC, b.root.count)
	}
}

func verifyNode(n node, parent *interiorNode, t *testing.T) {
	switch nn := n.(type) {
	case *interiorNode:
		if nn.count < MaxKC/2+1 {
			t.Errorf("interior.min.child: want >= %d, got = %d", MaxKC/2+1, nn.count)
		}

		if nn.count > MaxKC {
			t.Errorf("interior.max.child: want <= %d, got = %d", MaxKC, nn.count)
		}

		if nn.parent() != parent {
			t.Errorf("interior.parent: want = %p, got = %p", parent, nn.parent())
		}

		var last int
		for i := 0; i < nn.count; i++ {
			key := nn.kcs[i].key
			if key != 0 && key < last {
				t.Errorf("interior.sort.key: want > %d, got = %d", last, key)
			}
			last = key

			if i == nn.count-1 && key != 0 {
				t.Errorf("interior.last.key: want = 0, got = %d", key)
			}

			verifyNode(nn.kcs[i].child, nn, t)
		}

	case *leafNode:
		if nn.parent() != parent {
			t.Errorf("leaf.parent: want = %p, got = %p", parent, nn.parent())
		}

		if nn.count < MaxKV/2+1 {
			t.Errorf("leaf.min.child: want >= %d, got = %d", MaxKV/2+1, nn.count)
		}

		if nn.count > MaxKV {
			t.Errorf("leaf.max.child: want <= %d, got = %d", MaxKV, nn.count)
		}
	}
}

func verifyLeaf(leftMost *leafNode, count int, t *testing.T) {
	curr := leftMost
	last := 0
	c := 0

	for curr != nil {
		for i := 0; i < curr.count; i++ {
			key := curr.kvs[i].key

			if key <= last {
				t.Errorf("leaf.sort.key: want > %d, got = %d", last, key)
			}
			last = key
			c++
		}
		curr = curr.next
	}

	if c != count {
		t.Errorf("leaf.count: want = %d, got = %d", count, c)
	}
}

func findLeftMost(n node) *leafNode {
	switch nn := n.(type) {
	case *interiorNode:
		return findLeftMost(nn.kcs[0].child)
	case *leafNode:
		return nn
	default:
		panic("")
	}
}
