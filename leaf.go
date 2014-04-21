package bplustree

import (
	//"log"
	"sort"
)

type kv struct {
	key   int
	value string
}

type kvs [MaxKV]kv

func (a *kvs) Len() int           { return len(a) }
func (a *kvs) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a *kvs) Less(i, j int) bool { return a[i].key < a[j].key }

type leafNode struct {
	kvs   kvs
	count int
	next  *leafNode
	p     *interiorNode
}

func newLeafNode(p *interiorNode) *leafNode {
	return &leafNode{
		p: p,
	}
}

// find finds the index of a key in the leaf node.
// If the key exists in the node, it returns the index and true.
// If the key does not exist in the node, it returns index to
// insert the key (the index of the smallest key in the node that larger
// than the given key) and false.
func (l *leafNode) find(key int) (int, bool) {
	c := func(i int) bool {
		return l.kvs[i].key >= key
	}

	i := sort.Search(l.count, c)

	if i < l.count && l.kvs[i].key == key {
		return i, true
	}

	return i, false
}

// insert
func (l *leafNode) insert(key int, value string) (int, bool) {
	i, ok := l.find(key)

	if ok {
		//log.Println("insert.replace", i)
		l.kvs[i].value = value
		return 0, false
	}

	if !l.full() {
		copy(l.kvs[i+1:], l.kvs[i:l.count])
		l.kvs[i].key = key
		l.kvs[i].value = value
		l.count++
		return 0, false
	}

	next := l.split()

	if key < next.kvs[0].key {
		l.insert(key, value)
	} else {
		next.insert(key, value)
	}

	return next.kvs[0].key, true
}

func (l *leafNode) split() *leafNode {
	next := newLeafNode(nil)

	copy(next.kvs[0:], l.kvs[l.count/2+1:])

	next.count = MaxKV - l.count/2 - 1
	next.next = l.next

	l.count = l.count/2 + 1
	l.next = next

	return next
}

func (l *leafNode) full() bool { return l.count == MaxKV }

func (l *leafNode) parent() *interiorNode { return l.p }

func (l *leafNode) setParent(p *interiorNode) { l.p = p }
