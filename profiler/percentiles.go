package profiler

import (
	"github.com/montanaflynn/stats"
	"time"
	"fmt"
	"strings"
	"strconv"
	"log"
	"sort"
	"math"
)

type M map[string]interface{}

type ProfileRecord struct {
	//Id           bson.ObjectId `bson:"_id"`
	Version      string    `json:"version" bson:"version"`
	Timestamp    int64     `json:"timestamp" bson:"timestamp"`
	ShortMessage string    `json:"short_message" bson:"short_message"`
	Room         string    `json:"room" bson:"room"`
	Process      int       `json:"process" bson:"process"`
	Level        string    `json:"level" bson:"level"`
	Host         string    `json:"host" bson:"host"`
	FullMessage  string    `json:"full_message" bson:"full_message"`
	Environment  string    `json:"environment" bson:"environment"`
	Elapsed      float64   `json:"elapsed" bson:"elapsed"`
	AppName      string    `json:"app_name" bson:"app_name"`
	Time         time.Time `json:"time" bson:"time"`
}

type ProfileRecordView struct {
	Version      string    `json:"version" bson:"version"`
	ShortMessage string    `json:"short_message" bson:"short_message"`
	Room         string    `json:"room" bson:"room"`
	Process      int       `json:"process" bson:"process"`
	FullMessage  string    `json:"full_message" bson:"full_message"`
	AppName      string    `json:"app_name" bson:"app_name"`
	Time         time.Time `json:"time" bson:"time"`

	FullTime      float64 `json:"full_time"`
	Count         int     `json:"count"`
	Percentile2   float64 `json:"2_percentile"`
	Percentile25  float64 `json:"25_percentile"`
	Percentile50  float64 `json:"50_percentile"`
	Percentile75  float64 `json:"75_percentile"`
	Percentile90  float64 `json:"90_percentile"`
	Percentile98  float64 `json:"98_percentile"`
	Percentile100 float64 `json:"100_percentile"`
}

func HashRecord(fields []string, record ProfileRecord) string {
	var hash []string
	for _, key := range fields {
		switch key {
		case "app_name":
			hash = append(hash, record.AppName)
		case "process":
			hash = append(hash, strconv.Itoa(record.Process))
		case "version":
			hash = append(hash, record.Version)
		case "room":
			hash = append(hash, record.Room)
		case "short_message":
			hash = append(hash, record.ShortMessage)
		case "full_message":
			hash = append(hash, record.FullMessage)
		case "time":
			hash = append(hash, record.Time.Format("2006-01-02T15:04:05Z"))
		}
	}
	return strings.Join(hash, "*")

}

func copyslice(input stats.Float64Data) stats.Float64Data {
	s := make(stats.Float64Data, input.Len())
	copy(s, input)
	return s
}

func sortedCopy(input stats.Float64Data) (copy stats.Float64Data) {
	copy = copyslice(input)
	sort.Float64s(copy)
	return
}

func perc(input stats.Float64Data, percent float64) (percentile float64, err error) {
	// Find the length of items in the slice
	il := input.Len()

	// Return an error for empty slices
	if il == 0 {
		return math.NaN(), stats.EmptyInput
	}

	// Return error for less than 0 or greater than 100 percentages
	if percent < 0 || percent > 100 {
		return math.NaN(), stats.BoundsErr
	}

	// Return the last item
	if percent == 100.0 {
		return input[il-1], nil
	}

	// Find ordinal ranking
	or := int(math.Ceil(float64(il) * percent / 100))

	// Return the item that is in the place of the ordinal rank
	if or == 0 {
		return input[0], nil
	}
	return input[or-1], nil
}

func GeneratePercentiles2(item ProfileRecord, seq []float64) ProfileRecordView {
	view := ProfileRecordView{AppName: item.AppName, Process: item.Process, Version: item.Version, Room: item.Room, ShortMessage: item.ShortMessage, FullMessage: item.FullMessage, Time: item.Time}

	var fullTime float64
	var count int

	for _, val := range seq {
		fullTime += val
		count++
	}
	view.FullTime = fullTime
	view.Count = count

	sort.Float64s(seq)

	p2, err := perc(seq, 2)
	p25, err := perc(seq, 25)
	p50, err := perc(seq, 50)
	p75, err := perc(seq, 75)
	p90, err := perc(seq, 90)
	p98, err := perc(seq, 98)
	p100, err := perc(seq, 100)

	//p2, err := stats.Percentile(seq, 2)
	//p25, err := stats.Percentile(seq, 25)
	//p50, err := stats.Percentile(seq, 50)
	//p75, err := stats.Percentile(seq, 75)
	//p90, err := stats.Percentile(seq, 90)
	//p98, err := stats.Percentile(seq, 98)
	//p100, err := stats.Percentile(seq, 100)
	if err != nil {
		log.Fatal(err)
	}
	view.Percentile2 = p2
	view.Percentile25 = p25
	view.Percentile50 = p50
	view.Percentile75 = p75
	view.Percentile90 = p90
	view.Percentile98 = p98
	view.Percentile100 = p100

	return view
}

func GeneratePercentiles(item ProfileRecord, seq []float64) ProfileRecordView {
	view := ProfileRecordView{AppName: item.AppName, Process: item.Process, Version: item.Version, Room: item.Room, ShortMessage: item.ShortMessage, FullMessage: item.FullMessage, Time: item.Time}

	for _, val := range seq {
		view.FullTime += val
		view.Count++
	}
	p2, err := stats.PercentileNearestRank(seq, 2)
	p25, err := stats.PercentileNearestRank(seq, 25)
	p50, err := stats.PercentileNearestRank(seq, 50)
	p75, err := stats.PercentileNearestRank(seq, 75)
	p90, err := stats.PercentileNearestRank(seq, 90)
	p98, err := stats.PercentileNearestRank(seq, 98)
	p100, err := stats.PercentileNearestRank(seq, 100)
	if err != nil {
		log.Fatal(err)
	}
	view.Percentile2 = p2
	view.Percentile25 = p25
	view.Percentile50 = p50
	view.Percentile75 = p75
	view.Percentile90 = p90
	view.Percentile98 = p98
	view.Percentile100 = p100

	return view
}

func DrainMapper(fields []string, dbHelper *DBHelper) (map[string]ProfileRecord, map[string][]float64) {
	items := make(map[string]ProfileRecord)
	values := make(map[string][]float64)

	tmpRecord := ProfileRecord{}
	count := 0
	s1 := time.Now()

	iter := dbHelper.FindAll()

	for iter.Next(&tmpRecord) {
		count += 1
		//if count == 100000 {
		//	break
		//}
		key := HashRecord(fields, tmpRecord)
		items[key] = tmpRecord

		if _, ok := values[key]; ok {
			values[key] = append(values[key], tmpRecord.Elapsed)
		} else {
			values[key] = []float64{tmpRecord.Elapsed}
		}
	}
	s2 := time.Now()
	log.Println("drained", len(items), len(values), s2.Sub(s1))
	return items, values
}

func ReduceRecordsByFields(fields []string, timeFilter TimeChecker, items map[string]ProfileRecord, values map[string][]float64) (map[string]ProfileRecord, map[string][]float64) {
	newItems := make(map[string]ProfileRecord)
	newValues := make(map[string][]float64)

	for k, v := range items {
		if timeFilter != nil && !timeFilter.CheckTime(v.Time) {
			continue
		}
		seq := (values)[k]
		key := HashRecord(fields, v)
		//newItems[key] = v

		if _, ok := newValues[key]; ok {
			//{"app_name", "process", "version", "room", "short_message", "full_message", "time"}
			item := newItems[key]

			if v.AppName != newItems[key].AppName {
				item.AppName = "mixed"
			}
			if v.Process != newItems[key].Process {
				item.Process = -1
			}
			if v.Version != newItems[key].Version {
				item.Version = "mixed"
			}
			if v.Room != newItems[key].Room {
				item.Room = "mixed"
			}
			if v.ShortMessage != newItems[key].ShortMessage {
				item.ShortMessage = "mixed"
			}
			if v.FullMessage != newItems[key].FullMessage {
				item.FullMessage = "mixed"
			}
			if v.Time != newItems[key].Time {
				item.Time = time.Unix(0, 0).UTC()
			}
			newItems[key] = item

			newValues[key] = append(newValues[key], seq...)
		} else {
			newItems[key] = v
			newValues[key] = seq
		}

	}

	return newItems, newValues
}

func DrainByDate(start time.Time, end time.Time, fields []string, dbHelper *DBHelper) (map[string]ProfileRecord, map[string][]float64) {
	items := make(map[string]ProfileRecord)
	values := make(map[string][]float64)

	tmpRecord := ProfileRecord{}
	count := 0
	s1 := time.Now()

	iter := dbHelper.FindByDate(start, end)

	for iter.Next(&tmpRecord) {
		count += 1
		key := HashRecord(fields, tmpRecord)
		items[key] = tmpRecord

		if _, ok := values[key]; ok {
			values[key] = append(values[key], tmpRecord.Elapsed)
		} else {
			values[key] = []float64{tmpRecord.Elapsed}
		}
	}
	s2 := time.Now()
	fmt.Println(len(items), len(values), s2.Sub(s1))
	return items, values
}

func Process(items map[string]ProfileRecord, values map[string][]float64) *[]ProfileRecordView {
	var views []ProfileRecordView
	var a []int
	for _, v := range values {
		a = append(a, len(v))
	}
	c := 0
	for k, item := range items {
		seq := values[k]
		view := GeneratePercentiles2(item, seq)
		views = append(views, view)
		c++
	}
	return &views
}
