package main

import (
	"./persister"
	"./storage"
	"fmt"
	"time"
)

func main() {
	stor := new(storage.Storage)
	stor.New()
	fmt.Println(stor)
	m := make(map[string]int)
	m["test"] = 11
	s := make([]int, 3)
	s = append(s, 1)

	if r, e := stor.Set("stringkey", "stringvalue", 10); e != nil {
		fmt.Println(e)
	} else {
		fmt.Println(r)
	}
	if r, e := stor.Update("stringkey", "stringasdvalue", 10); e != nil {
		fmt.Println(e)
	} else {
		fmt.Println(r)
	}

	stor.Set("listkey", s, time.Second*5)
	stor.Set("dictkey", m, -1)
	fmt.Println("Storage before running checker")
	fmt.Println(stor)
	// persister.Dump(stor.GetNodes())
	persister.RunPersister(stor, time.Second*3)
	storage.RunChecker(stor, time.Second*3)
	time.Sleep(time.Second * 9)
	storage.StopChecker(stor)
	fmt.Println("Storage after running checker")
	fmt.Println(stor)
	fmt.Println("End")

}
