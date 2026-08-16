package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func main() {
	flag.Parse()

	if dir := filepath.Dir(*reportFile); dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Printf("directory making error %s: %v\n", dir, err)
			return
		}
	}
	f, err := os.Create(*reportFile)
	if err != nil {
		fmt.Printf("file creating error %s: %v\n", *reportFile, err)
		return
	}
	defer f.Close()
	bw := bufio.NewWriter(f)
	defer bw.Flush()

	out := io.MultiWriter(os.Stdout, bw)

	fprintHeader(out)
	fmt.Fprintf(out, "  Report:           %s\n\n", *reportFile)

	if !readyForTest("Start the server WITHOUT the cache") {
		fmt.Println("	canceled")
		return
	}

	if err := checkServer(*baseURL); err != nil {
		return
	}

	fmt.Fprintln(out, "\n -- SEED (RUN A) --")
	fmt.Fprintf(out, "  Creating %d users... ", *numUsers)
	userIDs, err := seedUsers(*baseURL, *numUsers)
	if err != nil {
		fmt.Fprintf(out, "ERROR: %v\n", err)
		return
	}
	fmt.Fprintf(out, "OK (IDs: %v)\n", userIDs)

	fmt.Fprintf(out, "  Creating %d tasks per user... ", *tasksPerUser)
	tasksByUser, err := seedTasks(*baseURL, userIDs, *tasksPerUser)
	if err != nil {
		fmt.Fprintf(out, "ERROR: %v\n", err)
		return
	}
	var total int
	for _, t := range tasksByUser {
		total += len(t)
	}
	fmt.Fprintf(out, "OK (%d tasks)\n", total)

	fmt.Fprintln(out, "\n  ════════════════════════════════════════════")
	fmt.Fprintln(out, "    RUN A: NO CACHE")
	fmt.Fprintln(out, "  ════════════════════════════════════════════")

	resultA := runFullTest(*baseURL, userIDs, tasksByUser)

	fprintSnapTable(out, "A / READ-HEAVY", resultA.ReadSnaps)
	fprintSnapTable(out, "A / MIXED — read", resultA.MixedReadSnaps)
	fprintSnapTable(out, "A / MIXED — mutations", resultA.MixedWriteSnaps)

	fmt.Fprintln(out)
	fmt.Fprintln(out, "  ────────────────────────────────────────────────────────────")
	fmt.Fprintln(out, "  Run A completed. Please RESET the database now (e.g., restart containers)!")
	fmt.Fprintln(out, "  ────────────────────────────────────────────────────────────")

	if !readyForTest("Start the server WITH the cache (make sure DB is empty and Redis is up)") {
		fmt.Println("   canceled")
		return
	}

	if err := checkServer(*baseURL); err != nil {
		return
	}

	fmt.Fprintln(out, "\n -- SEED (RUN B) --")
	fmt.Fprintf(out, "  Creating %d users... ", *numUsers)
	userIDs, err = seedUsers(*baseURL, *numUsers)
	if err != nil {
		fmt.Fprintf(out, "ERROR: %v\n", err)
		return
	}
	fmt.Fprintf(out, "OK (IDs: %v)\n", userIDs)

	fmt.Fprintf(out, "  Creating %d tasks per user... ", *tasksPerUser)
	tasksByUser, err = seedTasks(*baseURL, userIDs, *tasksPerUser)
	if err != nil {
		fmt.Fprintf(out, "ERROR: %v\n", err)
		return
	}
	var totalB int
	for _, t := range tasksByUser {
		totalB += len(t)
	}
	fmt.Fprintf(out, "OK (%d tasks)\n", totalB)

	fmt.Fprintln(out, "\n  ════════════════════════════════════════════")
	fmt.Fprintln(out, "    RUN B: WITH CACHE")
	fmt.Fprintln(out, "  ════════════════════════════════════════════")

	resultB := runFullTest(*baseURL, userIDs, tasksByUser)

	fprintSnapTable(out, "B / READ-HEAVY", resultB.ReadSnaps)
	fprintSnapTable(out, "B / MIXED — read", resultB.MixedReadSnaps)
	fprintSnapTable(out, "B / MIXED — mutations", resultB.MixedWriteSnaps)

	fmt.Fprintln(out)
	fmt.Fprintln(out, "  ────────────────────────────────────────────────────────────")
	fmt.Fprintln(out, "  Run B completed.")
	fmt.Fprintln(out, "  ────────────────────────────────────────────────────────────")

	fprintCrossCompare(out, "no-cache", "with-cache", []compareSection{
		{"READ-HEAVY (read-only)", resultA.ReadSnaps, resultB.ReadSnaps},
		{"MIXED — чтение (with cache invalidation)", resultA.MixedReadSnaps, resultB.MixedReadSnaps},
		{"MIXED — mutations", resultA.MixedWriteSnaps, resultB.MixedWriteSnaps},
	})

	fmt.Printf("\n  ✔ Report saved successfully: %s\n", *reportFile)
}
