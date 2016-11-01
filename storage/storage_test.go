package storage

import (
	"testing"
)

func TestEmptyStorage(t *testing.T) {
	stor := new(Storage)
	stor.New()
	got := stor.Len()
	if got != 0 {
		t.Errorf("got %d, want %d", got, 0)
	}
}

func TestStorageWithMap(t *testing.T) {
	stor := new(Storage)
	stor.New()
	m := make(map[string]int)
	m["test"] = 11
	stor.Set("map", m, 1)
	got := stor.Len()
	if got != 1 {
		t.Errorf("got %d, want %d", got, 1)
	}

	if got, ok := stor.Get("map"); ok != nil {
		t.Errorf("got %d, want %d", got, m)
	}
}

func TestStorageWithSlice(t *testing.T) {
	stor := new(Storage)
	stor.New()
	s := make([]int, 3)
	s = append(s, 1)
	stor.Set("slice", s, 1)
	got := stor.Len()
	if got != 1 {
		t.Errorf("got %d, want %d", got, 1)
	}

	if got, ok := stor.Get("slice"); ok != nil {
		t.Errorf("got %d, want %d", got, s)
	}
}

func TestStorageWithString(t *testing.T) {
	stor := new(Storage)
	stor.New()
	stor.Set("stringkey", "stringvalue", 10)
	got := stor.Len()
	if got != 1 {
		t.Errorf("got %d, want %d", got, 1)
	}

	if got, ok := stor.Get("stringkey"); ok != nil {
		t.Errorf("got %d, want %d", got, "stringvalue")
	}
}

func TestStorageUpdateString(t *testing.T) {
	stor := new(Storage)
	stor.New()
	stor.Set("stringkey", "stringvalue", 10)
	got := stor.Len()
	if got != 1 {
		t.Errorf("got %d, want %d", got, 1)
	}

	if got, ok := stor.Update("stringkey", "newvalue", 11); ok != nil {
		t.Errorf("got %d, want %d", got, "newvalue")
	}
	n, _ := stor.Get("stringkey")
	node := n.(Node)
	if got, _ := node.Get("newvalue"); got != "newvalue" {
		t.Errorf("got %d, want %d", got, "newvalue")
	}
}

func TestStorageUpdateSlice(t *testing.T) {
	stor := new(Storage)
	stor.New()
	s := make([]int, 3)
	s = append(s, 1)
	stor.Set("slice", s, 1)
	got := stor.Len()
	if got != 1 {
		t.Errorf("got %d, want %d", got, 1)
	}

	newSlice := make([]int, 1)
	newSlice[0] = 1
	if _, ok := stor.Update("slice", newSlice, 11); ok != nil {
		t.Errorf("Cant update %v", ok)
	}
	n, _ := stor.Get("slice")
	node := n.(Node)
	if got, _ := node.Get(0); got != 1 {
		t.Errorf("got %d, want %d", got, 1)
	}
}

func TestStorageUpdateMap(t *testing.T) {
	stor := new(Storage)
	stor.New()
	m := make(map[string]int)
	m["test"] = 11
	stor.Set("map", m, 1)
	got := stor.Len()
	if got != 1 {
		t.Errorf("got %d, want %d", got, 1)
	}
	newMap := make(map[string]int)
	newMap["newtest"] = 22
	if _, ok := stor.Update("map", newMap, 11); ok != nil {
		t.Errorf("Cant update %v", ok)
	}
	n, _ := stor.Get("map")
	node := n.(Node)
	if got, _ := node.Get("newtest"); got != 22 {
		t.Errorf("got %d, want %d", got, 22)
	}
}

func TestStorageDeleteMapNode(t *testing.T) {
	stor := new(Storage)
	stor.New()
	m := make(map[string]int)
	m["test"] = 11
	stor.Set("map", m, 1)
	if got := stor.Len(); got != 1 {
		t.Errorf("got %d, want %d", got, 1)
	}
	stor.Delete("map")
	if got := stor.Len(); got != 0 {
		t.Errorf("got %d, want %d", got, 0)
	}
}

func TestStorageDeleteStringNode(t *testing.T) {
	stor := new(Storage)
	stor.New()
	stor.Set("stringkey", "value", 1)
	if got := stor.Len(); got != 1 {
		t.Errorf("got %d, want %d", got, 1)
	}
	stor.Delete("stringkey")
	if got := stor.Len(); got != 0 {
		t.Errorf("got %d, want %d", got, 0)
	}
}

func TestStorageGetNodes(t *testing.T) {
	stor := new(Storage)
	stor.New()
	stor.Set("stringkey", "value", 1)
	if got := stor.Len(); got != 1 {
		t.Errorf("got %d, want %d", got, 1)
	}

	nodes := stor.GetNodes()
	if got := len(*nodes); got != 1 {
		t.Errorf("got %d, want %d", got, 1)
	}
}
