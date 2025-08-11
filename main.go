package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/go-vgo/robotgo"
)

func main() {
	filePath := flag.String("file", "", "path to file whose contents will be typed")
	delayStr := flag.String("delay", "0s", "delay before typing (e.g. 3s, 500ms)")
	flag.Parse()

	var input string
	var err error

	if *filePath != "" {
		b, e := os.ReadFile(*filePath)
		if e != nil {
			fmt.Fprintf(os.Stderr, "error reading file: %v\n", e)
			os.Exit(1)
		}
		input = string(b)
	} else {
		// Read from stdin if data is being piped
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			b, e := io.ReadAll(os.Stdin)
			if e != nil {
				fmt.Fprintf(os.Stderr, "error reading stdin: %v\n", e)
				os.Exit(1)
			}
			// Trim common trailing newlines from pipes like echo
			input = strings.TrimRight(string(b), "\r\n")
		} else {
			fmt.Fprintln(os.Stderr, "no input: provide data via stdin or --file")
			os.Exit(2)
		}
	}

	fmt.Printf("input: %q\n", input)

	d, err := time.ParseDuration(*delayStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid delay: %v\n", err)
		os.Exit(1)
	}

	if d > 0 {
		time.Sleep(d)
	}

	// Split by newlines and type each line separately
	lines := strings.Split(input, "\n")
	for i, line := range lines {
		// Trim trailing whitespace from each line
		line = strings.TrimRight(line, " \t\r")

		// Press enter before typing (except for the first line)
		if i > 0 {
			time.Sleep(timeout * time.Millisecond)
			robotgo.KeyTap(robotgo.Enter)
			time.Sleep(timeout * time.Millisecond)
		}

		if len(line) == 0 {
			continue
		}

		cols := strings.Split(line, "\t")
		for j, col := range cols {
			if j > 0 {
				time.Sleep(timeout * time.Millisecond)
				robotgo.KeyTap(robotgo.Tab)
				time.Sleep(timeout * time.Millisecond)
			}
			if len(col) > 0 {
				fmt.Printf("line %d col %d: %q\n", i+1, j+1, col)
				robotgo.TypeStr(col)
			}
		}
	}
}

const timeout = 10
