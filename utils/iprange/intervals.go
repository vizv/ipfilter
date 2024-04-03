// Adapted from https://github.com/antongulenko/merge-intervals
// MIT License
// Copyright (c) 2020 Anton Gulenko

package iprange

import "sort"

var _ sort.Interface = (*Intervals)(nil)

// Intervals contain a slice of Interval data. Intervals implements sort.Interface to sort all intervals based on their
// lower bounds (From field).
type Intervals []Interval

func (intervals Intervals) Len() int {
	return len(intervals)
}

func (intervals Intervals) Swap(x, y int) {
	intervals[x], intervals[y] = intervals[y], intervals[x]
}

func (intervals Intervals) Less(x, y int) bool {
	return intervals[x].From.ToInt().Cmp(intervals[y].From.ToInt().Int) < 0
}

// Merge is the core method of this module. The output is a list of intervals, where all overlapping input intervals
// are merged into one output interval.
// Merge first sorts the input intervals based on their lower bounds. Afterwards, it iterates over the sorted intervals
// and produces a new output interval every time an input interval does not overlap with its predecessor.
func (intervals Intervals) Merge() Intervals {
	if len(intervals) == 0 {
		return intervals
	}
	sort.Sort(intervals)
	var result Intervals

	// Initialize the aggregation variable to the lowest interval
	current := intervals[0]

	// Iterate the sorted intervals and keep merging them, until encountering a non-overlapping interval.
	// Since the intervals are sorted, a non-overlapping interval indicates the beginning of a new output interval.
	for _, interval := range intervals[1:] {
		merged, overlap := current.Merge(interval)
		if overlap {
			// Intervals overlap: continue merging.
			current = merged
		} else {
			result = append(result, current.Fix())
			current = interval
		}
	}

	// At the end, the current aggregation variable is the last result interval.
	return append(result, current.Fix())
}

func (intervals *Intervals) Append(f string, t string) {
	from, err := ParseIP(f)
	if err != nil {
		return
	}

	to, err := ParseIP(t)
	if err != nil {
		return
	}

	*intervals = append(*intervals, Interval{from, to})
}
