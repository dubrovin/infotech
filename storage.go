package main

import (
	"fmt"
	"sync"
	"time"
)

type AlreadyExistError struct {
	Key string
}

func (e AlreadyExistError) Error() string {
	return fmt.Sprintf("Node by key=%v already exist", e.Key)
}

type NotFoundError struct {
	Key string
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf("Node by key=%v does not exist", e.Key)
}

type ConversionError struct {
	Key string
}

func (e ConversionError) Error() string {
	return fmt.Sprintf("Cant convert to %v", e.Key)
}

type Node struct {
	val interface{}
	//time to live
	ttl int64
}

func (n *Node) get(i interface{}) (interface{}, error) {
	switch i.(type) {
	case int:
		s, ok := n.val.([]int)
		if ok {
			return s[i.(int)], nil
		}
		return nil, &ConversionError{"[]int"}

	case string:
		s, ok := n.val.(map[string]int)
		if ok {
			return s[i.(string)], nil
		}
		return nil, &ConversionError{"map[string]int"}

	}
	return nil, nil
}

func GetTimeDuration(t time.Duration) int64 {
	return time.Now().Add(t).UnixNano()
}

type Storage struct {
	nodes map[string]Node
	mu    sync.RWMutex
	// checker *Cherker
}

func (storage *Storage) Set(
	k string,
	v interface{},
	ttl time.Duration,
) (interface{}, error) {
	storage.mu.Lock()
	_, ok := storage.nodes[k]
	if ok {
		storage.mu.Unlock()
		return storage, &AlreadyExistError{k}
	}
	storage.nodes[k] = Node{val: v, ttl: GetTimeDuration(ttl)}
	storage.mu.Unlock()
	return storage.nodes[k], nil
}

func (storage *Storage) Get(k string) (interface{}, error) {
	storage.mu.RLock()
	_, ok := storage.nodes[k]
	if ok {
		storage.mu.RUnlock()
		return storage.nodes[k], nil
	}
	storage.mu.RUnlock()
	return nil, &NotFoundError{k}
}

func (storage *Storage) Update(
	k string,
	v interface{},
	ttl time.Duration,
) (interface{}, error) {
	storage.mu.Lock()
	_, ok := storage.nodes[k]
	if ok {
		storage.nodes[k] = Node{val: v, ttl: GetTimeDuration(ttl)}
		storage.mu.Unlock()
		return storage.nodes[k], nil
	}
	storage.mu.Unlock()
	return nil, &NotFoundError{k}
}

func (storage *Storage) Delete(k string) (bool, error) {
	storage.mu.Lock()
	_, ok := storage.nodes[k]
	if ok {
		delete(storage.nodes, k)
		storage.mu.Unlock()
		return true, nil
	}
	storage.mu.Unlock()
	return false, &NotFoundError{k}
}

func (storage *Storage) DeleteExpiredNodes() {
	currentTime := time.Now().UnixNano()
	storage.mu.Lock()
	for k, v := range storage.nodes {
		if v.ttl > 0 && v.ttl < currentTime {
			delete(storage.nodes, k)
		}
	}
	storage.mu.Unlock()
}

func main() {
	storage := new(Storage)
	fmt.Println(storage)
	storage.nodes = make(map[string]Node)
	fmt.Println(storage)

	m := make(map[string]int)
	m["test"] = 11
	s := make([]int, 3)
	s = append(s, 1)

	if r, e := storage.Set("stringkey", "stringvalue", 10); e != nil {
		fmt.Println(e)
	} else {
		fmt.Println(r)
	}
	fmt.Println(storage)
	if r, e := storage.Update("stringkey", "stringasdvalue", 10); e != nil {
		fmt.Println(e)
	} else {
		fmt.Println(r)
	}

	storage.Set("listkey", s, time.Second*2)
	storage.Set("dictkey", m, 10)
	fmt.Println(*storage)
	storage.DeleteExpiredNodes()
	fmt.Println(*storage)

}
