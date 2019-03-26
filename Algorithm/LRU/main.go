package main

import (
	"bytes"
	"fmt"
)

func main() {}

// LRUCache ...
type LRUCache struct {
	size int
	data map[string]interface{}
	keys *KeyList
}

// NewLRUCache ...
func NewLRUCache(size int) *LRUCache {
	return &LRUCache{
		size: size,
		data: make(map[string]interface{}, 0),
		keys: NewKeyList(),
	}
}

// Put ...
func (l *LRUCache) Put(key string, value interface{}) {
	if l.keys.size == l.size {
		if key, ok := l.keys.Pop(); ok {
			delete(l.data, key)
		}
	}
	l.data[key] = value
	l.keys.Push(key)
}

// Get ...
func (l *LRUCache) Get(key string) (value interface{}, ok bool) {
	value, ok = l.data[key]
	ok = l.keys.MoveToHead(key)
	return
}

// PrintKeys ...
func (l *LRUCache) PrintKeys() {
	fmt.Println(l.keys)
}

// KeyNode ...
type KeyNode struct {
	key  string
	next *KeyNode
}

// Append ...
func (n *KeyNode) Append(key string) {
	n.next = &KeyNode{key, nil}
}

func (n *KeyNode) String() string {
	var buf bytes.Buffer
	for tmp := n; tmp != nil; tmp = tmp.next {
		buf.Write([]byte(fmt.Sprintf("%s -> ", tmp.key)))
	}
	buf.Write([]byte("nil"))
	return buf.String()
}

// KeyList ...
type KeyList struct {
	head, tail *KeyNode
	size       int
}

// NewKeyList ...
func NewKeyList() *KeyList {
	root := &KeyNode{}
	return &KeyList{head: root, tail: root, size: 0}
}

// Push ...
func (l *KeyList) Push(key string) {
	l.tail.Append(key)
	l.tail = l.tail.next
	l.size++
}

// Pop ...
func (l *KeyList) Pop() (key string, ok bool) {
	if l.size > 0 {
		last := l.tail
		tmp := l.head.next
		for ; tmp != nil && tmp.next != last; tmp = tmp.next {
		}
		if tmp != nil && tmp.next == last {
			l.tail = tmp
			l.size--
			key, ok = last.key, true
		}
	}
	return
}

// MoveToHead ...
func (l *KeyList) MoveToHead(key string) bool {
	for pre, node := l.head, l.head.next; node != nil; pre, node = node, node.next {
		if node.key == key {
			if node == l.tail {
				l.tail = pre
			}
			pre.next = node.next
			l.head.next, node.next = node, l.head.next
			return true
		}
	}
	return false
}

func (l *KeyList) String() (s string) {
	if l.size > 0 {
		s = l.head.next.String()
	}
	return
}
