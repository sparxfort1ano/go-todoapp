package main

import (
	"flag"
	"time"
)

var (
	baseURL         = flag.String("base-url", "http://127.0.0.1:5050/api/v1", "Base URL API")
	numUsers        = flag.Int("users", 5, "Number of test users")
	tasksPerUser    = flag.Int("tasks-per-user", 20, "Number of tasks per user")
	phaseDuration   = flag.Duration("phase-duration", 20*time.Second, "Duration of each phase")
	concurrency     = flag.Int("concurrency", 50, "Number of goroutines per phase")
	readBurstCount  = flag.Int("read-burst", 10, "Number of consecutive GETs in a read-heavy iteration")
	mixedReadCount  = flag.Int("mixed-reads", 3, "Number of consecutive GETs in a read-write iteration")
	mixedWriteCount = flag.Int("mixed-writes", 3, "Number of consecutive create-update-delete's in a read-write iteration")
	reportFile      = flag.String("report", "load_test_report.txt", "Path for load test report file")
)
