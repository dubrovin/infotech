package persister

import (
	"../storage"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

func getTimeStr(t time.Time) string {
	return fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())
}

//dump all nodes to file
func Dump(nodes *map[string]storage.Node) {
	t := time.Now()
	f, _ := os.Create(strings.Join([]string{"persister/data/dump_", getTimeStr(t), ".log"}, ""))
	defer f.Close()
	for k, v := range *nodes {
		f.WriteString("key: " + k + " " + v.String() + "\n")
	}

}

type Persister struct {
	interval time.Duration
	stop     chan bool
	wg       sync.WaitGroup
	storage  *storage.Storage
}

func Persist(nodes *map[string]storage.Node) {
	f, err := os.OpenFile("persister/data/current.log", os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		f, _ = os.Create("persister/data/current.log")
	}
	defer f.Close()
	for k, v := range *nodes {
		if v.IsDumped() {
			return
		} else {
			f.WriteString("key: " + k + " " + v.String() + "\n")
			v.MakeDumped()
		}

	}
}

func (persister *Persister) PersistStorage() {
	persister.stop = make(chan bool)
	ticker := time.NewTicker(persister.interval)
	for {
		select {
		case <-ticker.C:
			persister.storage.LockMutex()
			Persist(persister.storage.GetNodes())
			persister.storage.UnlockMutex()
		case <-persister.stop:
			ticker.Stop()
			return
		}
	}
}

func RunPersister(storage *storage.Storage, interval time.Duration) {
	persi := &Persister{
		interval: interval,
		storage:  storage,
	}
	go persi.PersistStorage()
}
