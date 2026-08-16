package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

func runPhase(
	name string, dur time.Duration, conc int, userIDs []int,
	launch func(stop <-chan struct{}, wg *sync.WaitGroup, uid int, iters *atomic.Int64),
) {
	fmt.Printf("\n  ▶ Phase: %s (%s, %d goroutines)\n", name, dur, conc)
	pstop := make(chan struct{})
	var iters atomic.Int64
	go showProgress(pstop, name, dur, &iters)

	stop := make(chan struct{})
	var wg sync.WaitGroup
	for i := range conc {
		uid := userIDs[i%len(userIDs)]
		wg.Add(1)
		go func(uid int) {
			launch(stop, &wg, uid, &iters)
		}(uid)
	}

	<-time.After(dur)
	close(stop)
	wg.Wait()
	close(pstop)
	fmt.Printf("    Completed: %d iterations\n", iters.Load())
}

// showProgress displays statistics on the number of iterations to stdout at every tick.
func showProgress(stop <-chan struct{}, phase string, dur time.Duration, iters *atomic.Int64) {
	start := time.Now()
	tick := time.NewTicker(2 * time.Second)
	defer tick.Stop()
	for {
		select {
		case <-stop:
			fmt.Println()
			return
		case <-tick.C:
			el := time.Since(start).Round(time.Second)
			pct := float64(el) / float64(dur) * 100
			if pct > 100 {
				pct = 100
			}
			fmt.Printf("\r    [%s] %s / %s (%4.0f%%)  iterations: %d  ",
				phase, el, dur, pct, iters.Load())
		}
	}
}

type TestResult struct {
	ReadSnaps       []Snap
	MixedReadSnaps  []Snap
	MixedWriteSnaps []Snap
}

func runFullTest(base string, userIDs []int, tasksByUser map[int][]int) TestResult {
	warmUpCache(base, userIDs)

	rSelf := NewStats("GET /tasks?user_id=self")
	rOther := NewStats("GET /tasks?user_id=other")
	rByID := NewStats("GET /tasks/{id}")
	runPhase("READ-ONLY", *phaseDuration, *concurrency, userIDs,
		func(stop <-chan struct{}, wg *sync.WaitGroup, uid int, iters *atomic.Int64) {
			readHeavyWorker(stop, wg, base, uid, userIDs, tasksByUser, rSelf, rOther, rByID, iters)
		})

	warmUpCache(base, userIDs)

	mSelf := NewStats("GET /tasks?user_id=self")
	mOther := NewStats("GET /tasks?user_id=other")
	mByID := NewStats("GET /tasks/{id}")
	mPost := NewStats("POST /tasks")
	mPatch := NewStats("PATCH /tasks/{id}")
	mDel := NewStats("DELETE /tasks/{id}")
	runPhase("MIXED", *phaseDuration, *concurrency, userIDs,
		func(stop <-chan struct{}, wg *sync.WaitGroup, uid int, iters *atomic.Int64) {
			mixedWorker(stop, wg, base, uid, userIDs, tasksByUser,
				mSelf, mOther, mByID, mPost, mPatch, mDel, iters)
		})

	return TestResult{
		ReadSnaps:       []Snap{rSelf.Snapshot(), rOther.Snapshot(), rByID.Snapshot()},
		MixedReadSnaps:  []Snap{mSelf.Snapshot(), mOther.Snapshot(), mByID.Snapshot()},
		MixedWriteSnaps: []Snap{mPost.Snapshot(), mPatch.Snapshot(), mDel.Snapshot()},
	}
}
