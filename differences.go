package jman

import (
	"fmt"
)

type Differences []Difference

func (d Differences) Report() string {
	var report string
	for _, diff := range d {
		if diff.diff == "" {
			report += diff.subDiffs.Report()
		} else {
			report += fmt.Sprintf("%s\n", diff.String())
		}
	}
	return report
}

func (d Differences) HasKey(key string) bool {
	for _, d := range d {
		if d.path == key {
			return true
		}
	}

	return false
}

type Prefix string

type Difference struct {
	diff     string
	path     string
	subDiffs Differences
}

func (d Difference) String() string {
	if d.path[:2] != "$." {
		return fmt.Sprintf("%s.%s %s", base, d.path, d.diff)
	}
	return fmt.Sprintf("%s %s", d.path, d.diff)
}
