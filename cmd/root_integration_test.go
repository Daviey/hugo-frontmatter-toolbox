// Package cmd_test contains integration tests for the command line interface.
package cmd

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"
)

// captureOutput captures the output of a function that writes to os.Stdout and os.Stderr.
func captureOutput(f func()) string {
	r, w, _ := os.Pipe()
	stdout := os.Stdout
	stderr := os.Stderr
	os.Stdout = w
	os.Stderr = w

	outC := make(chan string)
	go func() {
		var buf bytes.Buffer
		_, _ = buf.ReadFrom(r)
		outC <- buf.String()
	}()

	f()
	_ = w.Close()
	os.Stdout = stdout
	os.Stderr = stderr
	return <-outC
}

// TestExecute_Help tests the execution of the command with the help flag.
func TestExecute_Help(t *testing.T) {
	exitFunc = func(code int) { panic(fmt.Sprintf("exit %d", code)) }
	defer func() { exitFunc = os.Exit }()

	output := captureOutput(func() {
		defer func() {
			if r := recover(); r != nil {
				// Check if the panic is from our exit function
				panicStr := fmt.Sprint(r)
				if !strings.HasPrefix(panicStr, "exit") {
					// If it's not our expected panic, propagate it
					panic(r)
				}
				// Otherwise, ignore it as it's expected
			}
		}()
		os.Args = []string{"hugo-frontmatter-toolbox"}
		Execute()
	})

	if !strings.Contains(output, "Batch edit Hugo frontmatter") {
		t.Errorf("Expected help text, got: %s", output)
	}
}

// TestExecute_Version tests the execution of the command with the version flag.
func TestExecute_Version(t *testing.T) {
	exitFunc = func(code int) { panic(fmt.Sprintf("exit %d", code)) }
	defer func() { exitFunc = os.Exit }()

	output := captureOutput(func() {
		defer func() {
			if r := recover(); r != nil {
				// Check if the panic is from our exit function
				panicStr := fmt.Sprint(r)
				if !strings.HasPrefix(panicStr, "exit") {
					// If it's not our expected panic, propagate it
					panic(r)
				}
				// Otherwise, ignore it as it's expected
			}
		}()
		os.Args = []string{"hugo-frontmatter-toolbox", "--version"}
		Execute()
	})

	if !strings.Contains(output, "hugo-frontmatter-toolbox") {
		t.Errorf("Expected version output, got: %s", output)
	}
}
