package jman

import (
	"fmt"
)

type differences []difference

func (d differences) report() string {
	var report string
	for _, diff := range d {
		if diff.diff == "" {
			report += diff.subDiffs.report()
		} else {
			report += fmt.Sprintf("%s\n", diff.String())
		}
	}
	return report
}

func (d differences) hasPath(path string) bool {
	for _, d := range d {
		if d.path == path {
			return true
		}
	}

	return false
}

type difference struct {
	diff     string
	path     string
	subDiffs differences
}

func (d difference) String() string {
	if d.path[:2] != "$." {
		return fmt.Sprintf("%s.%s %s", base, d.path, d.diff)
	}
	return fmt.Sprintf("%s %s", d.path, d.diff)
}
