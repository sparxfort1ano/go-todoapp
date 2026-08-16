package main

import (
	"fmt"
	"io"
	"time"
)

func fprintHeader(w io.Writer) {
	fmt.Fprintln(w, "╔═══════════════════════════════════════════════════════════════════════╗")
	fmt.Fprintln(w, "║    Load test /tasks — comparison: without cache vs with cache         ║")
	fmt.Fprintln(w, "╚═══════════════════════════════════════════════════════════════════════╝")
	fmt.Fprintln(w)
	fmt.Fprintf(w, "  Date:                  %s\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Fprintf(w, "  Base URL:              %s\n", *baseURL)
	fmt.Fprintf(w, "  Users:                 %d\n", *numUsers)
	fmt.Fprintf(w, "  Tasks per user seed:   %d\n", *tasksPerUser)
	fmt.Fprintf(w, "  Concurrency:           %d\n", *concurrency)
	fmt.Fprintf(w, "  Phase duration:        %s\n", *phaseDuration)
	fmt.Fprintf(w, "  Read burst:            %d\n", *readBurstCount)
	fmt.Fprintf(w, "  Mixed r/w:             %d reads + %d writes\n", *mixedReadCount, *mixedWriteCount)
	fmt.Fprintln(w)
}

func fprintSnapTable(w io.Writer, title string, rows []Snap) {
	fmt.Fprintf(w, "\n  ┌─ %s\n", title)
	fmt.Fprintf(w, "  │ %-32s │ %7s │ %6s │ %6s │ %9s │ %9s │ %9s │ %9s\n",
		"Endpoint", "OK", "FAIL", "Err%", "Avg", "p50", "p95", "p99")
	fmt.Fprintln(w, "  │ "+
		"─────────────────────────────────┼─────────┼────────┼────────┼───────────┼───────────┼───────────┼───────────")
	for _, r := range rows {
		if r.total() == 0 {
			fmt.Fprintf(w, "  │ %-32s │     —   │   —    │   —    │      —      │      —      │      —      │      —\n", r.name)
			continue
		}

		fmt.Fprintf(w, "  │ %-32s │ %7d │ %6d │ %6s │ %9s │ %9s │ %9s │ %9s\n",
			r.name, r.ok, r.fail, fmt.Sprintf("%.1f%%", r.errPct()),
			fmtDur(r.avg), fmtDur(r.p50), fmtDur(r.p95), fmtDur(r.p99))
	}
	fmt.Fprintln(w, "  └──")
}

type compareSection struct {
	title string
	a, b  []Snap
}

func fprintCrossCompare(w io.Writer, labelA, labelB string, sections []compareSection) {
	fmt.Fprintln(w)
	fmt.Fprintln(w, "  ╔═════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════╗")
	fmt.Fprintf(w, "  ║ Comparison: A = %-20s  vs  B = %-20s                                                                  ║\n", labelA, labelB)
	fmt.Fprintln(w, "  ╠═════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════╣")

	for _, sec := range sections {
		fmt.Fprintf(w, "\n  ┌─ %s\n", sec.title)

		fmt.Fprintf(w, "  │ %-32s │ %12s │ %14s │ %10s │ %12s │ %14s │ %10s │ %8s\n",
			"Endpoint", labelA+" p50", labelB+" p50", "Δ p50",
			labelA+" p95", labelB+" p95", "Δ p95", "Speedup")

		fmt.Fprintln(w, "  │ "+
			"─────────────────────────────────┼──────────────┼────────────────┼────────────┼──────────────┼────────────────┼────────────┼─────────")

		n := min(len(sec.a), len(sec.b))
		for i := range n {
			a, b := sec.a[i], sec.b[i]
			dp50 := b.p50 - a.p50
			dp95 := b.p95 - a.p95

			sp := float64(a.p50) / float64(b.p50)

			fmt.Fprintf(w, "  │ %-32s │ %12s │ %14s │ %10s │ %12s │ %14s │ %10s │ %7.2fx\n",
				a.name,
				fmtDur(a.p50), fmtDur(b.p50), fmtDurSigned(dp50),
				fmtDur(a.p95), fmtDur(b.p95), fmtDurSigned(dp95), sp)
		}
		fmt.Fprintln(w, "  └──")
	}

	fmt.Fprintln(w)
	fmt.Fprintln(w, "  ╠═════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════╣")
	fmt.Fprintln(w, "  ║ Speedup > 1.0x → B is faster than A.    Δ < 0 → B is faster.                                                                        ║")
	fmt.Fprintln(w, "  ╚═════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════╝")
}

func fmtDur(d time.Duration) string {
	us := d.Microseconds()
	if us < 1000 {
		return fmt.Sprintf("%dµs", us)
	}
	return fmt.Sprintf("%.2fms", float64(us)/1000)
}

func fmtDurSigned(d time.Duration) string {
	if d > 0 {
		return "+" + fmtDur(d)
	} else if d < 0 {
		return "-" + fmtDur(-d)
	}
	return fmtDur(0)
}
