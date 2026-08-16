package main

import (
	"math"
	"sort"
	"sync"
	"time"
)

type Stats struct {
	mu        sync.Mutex
	name      string
	latencies []time.Duration
	ok, fail  int64
}

func NewStats(name string) *Stats {
	return &Stats{name: name}
}

func (s *Stats) Record(d time.Duration, success bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.latencies = append(s.latencies, d)
	if success {
		s.ok++
	} else {
		s.fail++
	}
}

type Snap struct {
	name                    string
	ok, fail                int64
	avg, p50, p95, p99, max time.Duration
}

func (s *Stats) Snapshot() Snap {
	sn := Snap{name: s.name, ok: s.ok, fail: s.fail}
	if len(s.latencies) == 0 {
		return sn
	}

	sortedLatencies := make([]time.Duration, len(s.latencies))
	copy(sortedLatencies, s.latencies)
	sort.Slice(sortedLatencies, func(i, j int) bool {
		return sortedLatencies[i] < sortedLatencies[j]
	})

	var sum time.Duration
	for _, l := range sortedLatencies {
		sum += l
	}

	sn.avg = sum / time.Duration(len(sortedLatencies))
	sn.p50 = pctl(sortedLatencies, 50)
	sn.p95 = pctl(sortedLatencies, 95)
	sn.p99 = pctl(sortedLatencies, 99)
	sn.max = sortedLatencies[len(sortedLatencies)-1]

	return sn
}

func (s Snap) total() int64 {
	return s.ok + s.fail
}

func (s Snap) errPct() float64 {
	total := s.total()
	if total == 0 {
		return 0
	}

	return float64(s.fail) / float64(total) * 100.
}

func pctl(sorted []time.Duration, p float64) time.Duration {
	idx := int(math.Ceil(p/100*float64(len(sorted)))) - 1
	return sorted[idx]
}
