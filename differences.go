package jman

import (
	"fmt"
	"strings"
)

type Differences []Difference

func (d Differences) Report() string {
	var report string
	for _, diff := range d {
		report += fmt.Sprintf("%s\n", diff.String())
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

const (
	Expected Prefix = "expected"
	Actual   Prefix = "actual"
	Both     Prefix = "$"
)

type Difference struct {
	prefix Prefix
	diff   string
	path   string
}

func (d Difference) String() string {
	var prefixes []string
	if d.prefix != "" {
		prefixes = append(prefixes, string(d.prefix))
	}
	if d.path != "" {
		prefixes = append(prefixes, d.path)
	}
	prefix := strings.Join(prefixes, ".")
	return fmt.Sprintf("%s %s", prefix, d.diff)
}
