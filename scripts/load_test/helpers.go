package main

import (
	"bufio"
	"fmt"
	"math/rand/v2"
	"os"
	"strings"
)

func randomString(n int) string {
	const c = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = c[rand.IntN(len(c))]
	}
	return string(b)
}

func randomPhone() string {
	return fmt.Sprintf("+7%010d", rand.Int64N(1e10))
}

var stdinScanner = bufio.NewScanner(os.Stdin)

func readyForTest(prompt string) bool {
	fmt.Println()
	fmt.Printf("  ⏸  %s\n", prompt)
	fmt.Print("  ➜  Press Enter when you’re ready (or ‘q’ to cancel): ")
	stdinScanner.Scan()
	text := strings.TrimSpace(strings.ToLower(stdinScanner.Text()))
	return text != "q" && text != "quit" && text != "n" && text != "no"
}
