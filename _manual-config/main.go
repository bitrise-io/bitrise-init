package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/bitrise-io/bitrise-init/output"
	"github.com/bitrise-io/bitrise-init/scanner"
)

const exportDir = "generated"

func main() {
	os.Exit(run())
}

func run() int {
	fmt.Println("Generating manual config")

	scanResult, err := scanner.ManualConfig()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "scanner failed:", err)
		return 1
	}

	if err := os.MkdirAll(exportDir, 0700); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, fmt.Errorf("failed to create (%s): %s", exportDir, err))
		return 2
	}

	outputPth, err := output.WriteToFile(scanResult, output.YAMLFormat, filepath.Join(exportDir, "result"))
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "failed to save config:", err)
		return 3
	}

	fmt.Println("Config saved to", outputPth)

	return 0
}
