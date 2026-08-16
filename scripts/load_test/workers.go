package main

import (
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

func readHeavyWorker(
	stop <-chan struct{}, wg *sync.WaitGroup,
	base string, myUID int, allUIDs []int, tasksByUser map[int][]int,
	sListSelf, sListOther, sGetByID *Stats, iters *atomic.Int64,
) {
	defer wg.Done()

	myTasks := tasksByUser[myUID]
	for {
		otherUID := pickOther(myUID, allUIDs)
		otherTasks := tasksByUser[otherUID]

		for range *readBurstCount {
			select {
			case <-stop:
				return
			default:
			}
			{
				url := fmt.Sprintf("%s/tasks?user_id=%d&limit=10&offset=%d",
					base, myUID, rand.IntN(max(1, len(myTasks)-10)))
				start := time.Now()
				resp, _, err := doJSON(http.MethodGet, url, nil)
				sListSelf.Record(time.Since(start), err == nil && resp.StatusCode == http.StatusOK)
			}
			{
				url := fmt.Sprintf("%s/tasks?user_id=%d&limit=10&offset=%d",
					base, otherUID, rand.IntN(max(1, len(otherTasks)-10)))
				start := time.Now()
				resp, _, err := doJSON(http.MethodGet, url, nil)
				sListOther.Record(time.Since(start), err == nil && resp.StatusCode == http.StatusOK)
			}
			if len(myTasks) > 0 {
				tid := myTasks[rand.IntN(len(myTasks))]
				start := time.Now()
				resp, _, err := doJSON(http.MethodGet, fmt.Sprintf("%s/tasks/%d", base, tid), nil)
				sGetByID.Record(time.Since(start), err == nil && resp.StatusCode == http.StatusOK)
			}
		}
		iters.Add(1)
	}
}

func mixedWorker(
	stop <-chan struct{}, wg *sync.WaitGroup,
	base string, myUID int, allUIDs []int, tasksByUser map[int][]int,
	sListSelf, sListOther, sGetByID, sPost, sPatch, sDelete *Stats,
	iters *atomic.Int64,
) {
	defer wg.Done()

	myTasks := make([]int, len(tasksByUser[myUID]))
	copy(myTasks, tasksByUser[myUID])
	for {
		otherUID := pickOther(myUID, allUIDs)
		otherTasks := tasksByUser[otherUID]

		for range *mixedReadCount {
			select {
			case <-stop:
				return
			default:
			}
			{
				url := fmt.Sprintf("%s/tasks?user_id=%d&limit=10&offset=%d",
					base, myUID, rand.IntN(max(1, len(myTasks)-10)))
				start := time.Now()
				resp, _, err := doJSON(http.MethodGet, url, nil)
				sListSelf.Record(time.Since(start), err == nil && resp.StatusCode == http.StatusOK)
			}
			{
				url := fmt.Sprintf("%s/tasks?user_id=%d&limit=10&offset=%d",
					base, otherUID, rand.IntN(max(1, len(otherTasks)-10)))
				start := time.Now()
				resp, _, err := doJSON(http.MethodGet, url, nil)
				sListOther.Record(time.Since(start), err == nil && resp.StatusCode == http.StatusOK)
			}
			if len(myTasks) > 0 {
				tid := myTasks[rand.IntN(len(myTasks))]
				start := time.Now()
				resp, _, err := doJSON(http.MethodGet, fmt.Sprintf("%s/tasks/%d", base, tid), nil)
				sGetByID.Record(time.Since(start), err == nil && resp.StatusCode == http.StatusOK)
			}
		}

		for range *mixedWriteCount {
			select {
			case <-stop:
				return
			default:
			}
			{
				start := time.Now()
				resp, body, err := doJSON(http.MethodPost, base+"/tasks", CreateTaskRequest{
					Title:        fmt.Sprintf("Mixed %s [u:%d]", randomString(6), myUID),
					Description:  fmt.Sprintf("Mixed load task user %d", myUID),
					AuthorUserID: myUID,
				})
				sPost.Record(time.Since(start), err == nil && resp.StatusCode == http.StatusCreated)
				if err == nil && resp.StatusCode == http.StatusCreated {
					var t TaskResponse
					json.Unmarshal(body, &t)
					myTasks = append(myTasks, t.ID)
				}
			}
			if len(myTasks) > 0 {
				tid := myTasks[rand.IntN(len(myTasks))]
				start := time.Now()
				resp, _, err := doJSON(http.MethodPatch, fmt.Sprintf("%s/tasks/%d", base, tid), PatchTaskRequest{
					Title:     fmt.Sprintf("Upd %s", randomString(4)),
					Completed: rand.IntN(2) == 1,
				})
				sPatch.Record(time.Since(start), err == nil && resp.StatusCode == http.StatusOK)
			}
			if len(myTasks) > 1 {
				idx := rand.IntN(len(myTasks))
				tid := myTasks[idx]
				myTasks[idx] = myTasks[len(myTasks)-1]
				myTasks = myTasks[:len(myTasks)-1]

				start := time.Now()
				resp, _, err := doJSON(http.MethodDelete, fmt.Sprintf("%s/tasks/%d", base, tid), nil)
				sDelete.Record(time.Since(start), err == nil && resp.StatusCode == http.StatusNoContent)
			}
		}
		iters.Add(1)
	}
}

func pickOther(my int, all []int) int {
	if len(all) <= 1 {
		return my
	}
	o := my
	for o == my {
		o = all[rand.IntN(len(all))]
	}
	return o
}
