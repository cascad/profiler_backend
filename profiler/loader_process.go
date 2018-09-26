package profiler

import (
	"log"
	"time"
)

var fields = []string{"app_name", "process", "version", "room", "short_message", "full_message", "time"}

func ProccessLoad(raw *LocalRawStorage, process *LocalProcessStorage, dbh *DBHelper) {
	s1 := time.Now()
	items, values := DrainMapper(fields, dbh)
	raw.Lock()
	raw.Items = items
	raw.Values = values
	raw.Unlock()

	views := Process(items, values)
	process.Lock()
	process.Data = *views
	process.Unlock()
	s2 := time.Now()
	log.Println("processed", s2.Sub(s1), len(items), len(process.Data))
}
