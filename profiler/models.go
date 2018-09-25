package profiler

import (
	"time"
	"sync"
)

type LocalRawStorage struct {
	sync.Mutex
	Items  map[string]ProfileRecord
	Values map[string][]float64
}

type LocalProcessStorage struct {
	sync.Mutex
	Data []ProfileRecordView
}

type TimeChecker interface {
	CheckTime(time.Time) bool
}

type HandlerTableParams struct {
	StartTs time.Time `json:"start_ts"`
	EndTs   time.Time `json:"end_ts"`
}

type HandlerReduceParams struct {
	StartTs time.Time `json:"start_ts"`
	EndTs   time.Time `json:"end_ts"`
	Fields  []string  `json:"fields"`
}

func (p *HandlerTableParams) CheckTime(t time.Time) bool {
	start := p.StartTs
	if !start.IsZero() {
		if t.Before(start) {
			return false
		}
	}
	end := p.EndTs
	if !end.IsZero() {
		if t.After(end) {
			return false
		}
	}
	return true
}

func (p *HandlerReduceParams) CheckTime(t time.Time) bool {
	start := p.StartTs
	if !start.IsZero() {
		if t.Before(start) {
			return false
		}
	}
	end := p.EndTs
	if !end.IsZero() {
		if t.After(end) {
			return false
		}
	}
	return true
}
