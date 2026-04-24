package runner

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// filterWriter writes all output to log, and conditionally to out based on verbose/summary logic.
type filterWriter struct {
	log     io.Writer
	out     io.Writer
	verbose bool
	buf     bytes.Buffer
}

func (f *filterWriter) Write(p []byte) (int, error) {
	f.buf.Write(p)
	for {
		line, rest, found := strings.Cut(f.buf.String(), "\n")
		if !found {
			break
		}
		f.buf.Reset()
		f.buf.WriteString(rest)
		fmt.Fprintln(f.log, line)
		if f.verbose || isSummaryLine(line) {
			fmt.Fprintln(f.out, line)
		}
	}
	return len(p), nil
}

func (f *filterWriter) flush() {
	if f.buf.Len() > 0 {
		line := f.buf.String()
		fmt.Fprint(f.log, line)
		if f.verbose || isSummaryLine(line) {
			fmt.Fprint(f.out, line)
		}
		f.buf.Reset()
	}
}

// isSummaryLine returns true for lines that are useful signal in summary mode.
func isSummaryLine(line string) bool {
	return strings.Contains(line, ": error:") ||
		strings.Contains(line, ": warning:") ||
		strings.HasPrefix(line, "** BUILD") ||
		strings.HasPrefix(line, "** TEST") ||
		(strings.Contains(line, "Executed ") && strings.Contains(line, "test")) ||
		(strings.HasPrefix(line, "Test Suite '") && (strings.Contains(line, "passed") || strings.Contains(line, "failed")))
}

// Run executes a command in the given directory, streaming output live.
// In summary mode (verbose=false), only errors, warnings, and status lines are printed.
// Full output is always written to a temp log file.
func Run(dir string, verbose bool, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("no command provided")
	}

	logName := fmt.Sprintf("jig-%s-%s.log", filepath.Base(args[0]), time.Now().Format("20060102-150405"))
	logPath := filepath.Join(os.TempDir(), logName)
	logFile, err := os.Create(logPath)
	if err != nil {
		return fmt.Errorf("failed to create log file: %w", err)
	}
	defer logFile.Close()

	fw := &filterWriter{
		log:     logFile,
		out:     os.Stdout,
		verbose: verbose,
	}

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = dir
	cmd.Stdin = os.Stdin
	cmd.Stdout = fw
	cmd.Stderr = fw

	start := time.Now()
	cmdErr := cmd.Run()
	fw.flush()
	elapsed := time.Since(start)

	fmt.Printf("Elapsed: %.1fs\n", elapsed.Seconds())
	fmt.Printf("Log: %s\n", logPath)

	return cmdErr
}
