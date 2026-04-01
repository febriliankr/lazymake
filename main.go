package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/febriliankr/lazymake/internal/parser"
	"github.com/febriliankr/lazymake/internal/runner"
	"github.com/febriliankr/lazymake/internal/tui"
)

var version = "dev"

func main() {
	showVersion := flag.Bool("v", false, "print version")
	file := flag.String("f", "", "path to Makefile")
	recursive := flag.Bool("r", false, "search subdirectories for Makefiles")
	dryRun := flag.Bool("n", false, "print the make command instead of running it")
	flag.Parse()

	if *showVersion {
		fmt.Println("lazymake " + version)
		os.Exit(0)
	}

	extraArgs := flag.Args()

	// Find Makefiles
	var files []string
	var err error

	if *file != "" {
		files = []string{*file}
	} else {
		dir, _ := os.Getwd()
		files, err = parser.FindMakefiles(dir, *recursive)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	}

	if len(files) == 0 {
		fmt.Fprintln(os.Stderr, "no Makefile found")
		os.Exit(1)
	}

	// Parse all targets
	var targets []parser.Target
	for _, f := range files {
		t, err := parser.Parse(f)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error parsing %s: %v\n", f, err)
			continue
		}
		targets = append(targets, t...)
	}

	if len(targets) == 0 {
		fmt.Fprintln(os.Stderr, "no targets found")
		os.Exit(1)
	}

	// Launch TUI
	selected, err := tui.Run(targets)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	if selected == nil {
		os.Exit(0)
	}

	// Run the selected target
	fmt.Printf("\n  Running: make %s\n\n", selected.Name)
	if err := runner.Run(*selected, extraArgs, *dryRun); err != nil {
		os.Exit(1)
	}
}
