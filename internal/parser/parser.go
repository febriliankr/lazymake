package parser

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Target struct {
	Name        string
	Description string
	File        string
}

var (
	targetRe      = regexp.MustCompile(`^([a-zA-Z0-9][a-zA-Z0-9._-]*)\s*:`)
	inlineComment = regexp.MustCompile(`##\s*(.+)$`)
	commentRe     = regexp.MustCompile(`^##?\s*(.+)$`)
)

func Parse(path string) ([]Target, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var targets []Target
	var pendingComment string

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()

		// Check for comment line (# or ##)
		if m := commentRe.FindStringSubmatch(line); m != nil {
			pendingComment = strings.TrimSpace(m[1])
			continue
		}

		// Check for target line (skip variable assignments like FOO := bar)
		if m := targetRe.FindStringSubmatch(line); m != nil {
			rest := line[len(m[0]):]
			if strings.HasPrefix(rest, "=") || strings.HasPrefix(rest, ":") {
				pendingComment = ""
				continue
			}
			name := m[1]

			// Skip special targets
			if strings.HasPrefix(name, ".") {
				pendingComment = ""
				continue
			}

			// Skip pattern rules
			if strings.Contains(name, "%") {
				pendingComment = ""
				continue
			}

			desc := ""
			// Check for inline ## comment first
			if im := inlineComment.FindStringSubmatch(line); im != nil {
				desc = strings.TrimSpace(im[1])
			} else if pendingComment != "" {
				desc = pendingComment
			}

			targets = append(targets, Target{
				Name:        name,
				Description: desc,
				File:        path,
			})
			pendingComment = ""
			continue
		}

		// Non-comment, non-target line resets pending comment
		if strings.TrimSpace(line) != "" {
			pendingComment = ""
		}
	}

	return targets, scanner.Err()
}

func FindMakefiles(dir string, recursive bool) ([]string, error) {
	names := []string{"Makefile", "makefile", "GNUmakefile"}

	if !recursive {
		var found []string
		for _, name := range names {
			p := filepath.Join(dir, name)
			if _, err := os.Stat(p); err == nil {
				found = append(found, p)
			}
		}
		return found, nil
	}

	skipDirs := map[string]bool{
		".git": true, "node_modules": true, "vendor": true, "dist": true,
	}

	nameSet := make(map[string]bool)
	for _, n := range names {
		nameSet[n] = true
	}

	var found []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() && skipDirs[info.Name()] {
			return filepath.SkipDir
		}
		if !info.IsDir() && nameSet[info.Name()] {
			found = append(found, path)
		}
		return nil
	})
	return found, err
}
