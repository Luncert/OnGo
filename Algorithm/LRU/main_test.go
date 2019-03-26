package main

import (
	"fmt"
	"testing"
)

func TestKeyNode(t *testing.T) {
	n1 := &KeyNode{"a", nil}
	n1.Append("b")
	fmt.Println(n1)
}

func TestKeyList(t *testing.T) {
	list := NewKeyList()
	list.Push("a")
	list.Push("b")
	list.Push("c")
	list.MoveToHead("b")
	list.MoveToHead("c")
	list.MoveToHead("a")
	fmt.Println(list)
}

func TestLRUCache(t *testing.T) {
	l := NewLRUCache(4)
	l.Put("a", 0)
	l.Put("b", 0)
	l.Put("c", 0)
	l.Put("d", 0)
	l.PrintKeys()

	l.Get("c")
	l.PrintKeys()

	l.Get("d")
	l.PrintKeys()

	l.Put("e", 0)
	l.PrintKeys()

	l.Put("g", 0)
	l.PrintKeys()
}
