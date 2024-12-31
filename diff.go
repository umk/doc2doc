package main

import (
	"strings"

	"github.com/sergi/go-diff/diffmatchpatch"
)

const (
	diffPrefixEq   = " "
	diffPrefixPrev = "-"
	diffPrefixCur  = "+"
)

func renderDiff(diffs []diffmatchpatch.Diff) string {
	var lines []string

	var chunksPrev, chunksCur []string

	for _, diff := range diffs {
		diffLines := strings.Split(diff.Text, "\n")
		for i, chunk := range diffLines {
			if i > 0 {
				switch diff.Type {
				case diffmatchpatch.DiffDelete:
					lines = append(lines, diffPrefixPrev+strings.Join(chunksPrev, ""))

					chunksPrev = nil
				case diffmatchpatch.DiffInsert:
					lines = append(lines, diffPrefixCur+strings.Join(chunksCur, ""))

					chunksCur = nil
				case diffmatchpatch.DiffEqual:
					linePrev := strings.Join(chunksPrev, "")
					lineCur := strings.Join(chunksCur, "")

					if linePrev == lineCur {
						lines = append(lines, diffPrefixEq+linePrev)
					} else {
						lines = append(lines, diffPrefixPrev+linePrev)
						lines = append(lines, diffPrefixCur+lineCur)
					}

					chunksPrev = nil
					chunksCur = nil
				}
			}

			switch diff.Type {
			case diffmatchpatch.DiffDelete:
				chunksPrev = append(chunksPrev, chunk)
			case diffmatchpatch.DiffInsert:
				chunksCur = append(chunksCur, chunk)
			case diffmatchpatch.DiffEqual:
				chunksPrev = append(chunksPrev, chunk)
				chunksCur = append(chunksCur, chunk)
			}
		}
	}

	linePrev := strings.Join(chunksPrev, "")
	lineCur := strings.Join(chunksCur, "")

	if linePrev == lineCur {
		if len(linePrev) > 0 {
			lines = append(lines, diffPrefixEq+linePrev)
		}
	} else {
		if len(linePrev) > 0 {
			lines = append(lines, diffPrefixPrev+linePrev)
		}
		if len(lineCur) > 0 {
			lines = append(lines, diffPrefixCur+lineCur)
		}
	}

	return strings.Join(lines, "\n")
}
