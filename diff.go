package main

import (
	"fmt"
	"strings"

	"github.com/sergi/go-diff/diffmatchpatch"
)

func renderDiff(diffs []diffmatchpatch.Diff) string {
	var lines []string

	for _, diff := range diffs {
		switch diff.Type {
		case diffmatchpatch.DiffInsert:
			// Lines added are prefixed with "+"
			lines = append(lines, fmt.Sprintf("+ %s", diff.Text))
		case diffmatchpatch.DiffDelete:
			// Lines removed are prefixed with "-"
			lines = append(lines, fmt.Sprintf("- %s", diff.Text))
		}
	}

	return strings.Join(lines, "\n")
}
