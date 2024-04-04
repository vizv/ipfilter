// Adapted from https://github.com/antongulenko/merge-intervals
// MIT License
// Copyright (c) 2020 Anton Gulenko

package iprange

import (
	"fmt"
	"math/big"
)

// Interval contains the lower and upper bounds of an integer interval
type Interval struct {
	From *IP
	To   *IP
}

// Merge attempts to merge the receiving interval with the argument. The merge succeeds when the intervals overlap,
// in this case the merged interval is returned, along with a true-valued flag indicating the success. Otherwise,
// a false-valued flag indicates that the intervals do not overlap. This operation is symmetrical (the receiver and the
// argument can be exchanged with the same result).
//
// In case of a successful merge, the result is always a correct interval, i.e. result.From <= result.To
func (i Interval) Merge(other Interval) (result Interval, overlap bool) {
	other = other.Fix()
	overlap = i.Contains(other.From.ToInt().Int) || i.Contains(other.To.ToInt().Int)
	if overlap {
		result = i.Fix()
		if result.From.ToInt().Cmp(other.From.ToInt().Int) > 0 {
			result.From = other.From
		}
		if result.To.ToInt().Cmp(other.To.ToInt().Int) < 0 {
			result.To = other.To
		}
	}
	return
}

// Contains returns whether the receiving interval contains the argument integer (both the lower and upper bounds
// work inclusively).
func (i Interval) Contains(num *big.Int) bool {
	return big.NewInt(0).Add(num, big.NewInt(1)).Cmp(i.From.ToInt().Int) >= 0 && big.NewInt(0).Sub(num, big.NewInt(1)).Cmp(i.To.ToInt().Int) <= 0
}

// Fix swaps the From and To fields of the receiving interval, if To < From. This corrects wrong input intervals, where
// the bounds are exchanged.
func (i Interval) Fix() Interval {
	if i.To.ToInt().Cmp(i.From.ToInt().Int) < 0 {
		i.From, i.To = i.To, i.From
	}
	return i
}

// String returns a human-readable representation of the receiving interval.
func (i Interval) String() string {
	return fmt.Sprintf("[%v, %v]", i.From, i.To)
}
