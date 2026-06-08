package main

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type node struct {
	key   int
	value int
	left  *node
	right *node
}

func deleteNode(n *node, key int) (*node, bool) {
	if n == nil {
		return nil, false
	}

	var deleted bool
	if key < n.key {
		n.left, deleted = deleteNode(n.left, key)
	} else if key > n.key {
		n.right, deleted = deleteNode(n.right, key)
	} else {
		deleted = true
		if n.left == nil {
			return n.right, true
		} else if n.right == nil {
			return n.left, true
		}

		inheritor := n.right
		for inheritor.left != nil {
			inheritor = inheritor.left
		}

		n.key = inheritor.key
		n.value = inheritor.value

		n.right, _ = deleteNode(n.right, inheritor.key)
	}
	return n, deleted
}

func (m *OrderedMap) insertNode(n *node, key, value int) (*node, bool) {
	if n == nil {
		return &node{key: key, value: value}, true
	}

	if key < n.key {
		newLeft, inserted := m.insertNode(n.left, key, value)
		n.left = newLeft
		return n, inserted
	} else if key > n.key {
		newRight, inserted := m.insertNode(n.right, key, value)
		n.right = newRight
		return n, inserted
	}

	n.value = value
	return n, false
}

type OrderedMap struct {
	root *node
	size int
}

func NewOrderedMap() OrderedMap {
	return OrderedMap{}
}

func (m *OrderedMap) Insert(key, value int) {
	var inserted bool
	m.root, inserted = m.insertNode(m.root, key, value)
	if inserted {
		m.size++
	}
}

func (m *OrderedMap) Erase(key int) {
	var deleted bool
	m.root, deleted = deleteNode(m.root, key)
	if deleted {
		m.size--
	}
}

func (m *OrderedMap) Contains(key int) bool {
	curr := m.root
	for curr != nil {
		if key == curr.key {
			return true
		} else if key < curr.key {
			curr = curr.left
		} else {
			curr = curr.right
		}
	}
	return false
}

func (m *OrderedMap) Size() int {
	return m.size
}

func (m *OrderedMap) ForEach(action func(int, int)) {
	var traverse func(n *node)
	traverse = func(n *node) {
		if n == nil {
			return
		}
		traverse(n.left)
		action(n.key, n.value)
		traverse(n.right)
	}
	traverse(m.root)
}

func TestOrderedMap(t *testing.T) {
	data := NewOrderedMap()
	assert.Zero(t, data.Size())

	data.Insert(10, 10)
	data.Insert(5, 5)
	data.Insert(15, 15)
	data.Insert(2, 2)
	data.Insert(4, 4)
	data.Insert(12, 12)
	data.Insert(14, 14)

	assert.Equal(t, 7, data.Size())
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(3))
	assert.False(t, data.Contains(13))

	var keys []int
	expectedKeys := []int{2, 4, 5, 10, 12, 14, 15}
	data.ForEach(func(key, _ int) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))

	data.Erase(15)
	data.Erase(14)
	data.Erase(2)

	assert.Equal(t, 4, data.Size())
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(2))
	assert.False(t, data.Contains(14))

	keys = nil
	expectedKeys = []int{4, 5, 10, 12}
	data.ForEach(func(key, _ int) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))
}
