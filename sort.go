package timertxt

import (
	"errors"
	"sort"
	"time"
)

// Flags for defining sort element and order.
const (
	SORT_UNFINISHED_START = iota
	SORT_START_DATE_ASC
	SORT_START_DATE_DESC
	SORT_FINISH_DATE_ASC
	SORT_FINISH_DATE_DESC
)

// Sort allows a TimerList to be sorted by certain predefined fields.
// See constants SORT_* for fields and sort order.
func (timerlist *TimerList) Sort(sortFlag int) error {
	switch sortFlag {
	case SORT_UNFINISHED_START:
		timerlist.sortByUnfinishedThenStart()
	case SORT_START_DATE_ASC, SORT_START_DATE_DESC:
		timerlist.sortByStartDate(sortFlag)
	case SORT_FINISH_DATE_ASC, SORT_FINISH_DATE_DESC:
		timerlist.sortByFinishDate(sortFlag)
	default:
		return errors.New("Unrecognized sort option")
	}
	return nil
}

type timerlistSort struct {
	timerlists TimerList
	by         func(t1, t2 *Timer) bool
}

func (ts *timerlistSort) Len() int {
	return len(ts.timerlists)
}

func (ts *timerlistSort) Swap(l, r int) {
	ts.timerlists[l], ts.timerlists[r] = ts.timerlists[r], ts.timerlists[l]
}

func (ts *timerlistSort) Less(l, r int) bool {
	return ts.by(&ts.timerlists[l], &ts.timerlists[r])
}

func (timerlist *TimerList) sortBy(by func(t1, t2 *Timer) bool) *TimerList {
	ts := &timerlistSort{
		timerlists: *timerlist,
		by:         by,
	}
	sort.Sort(ts)
	return timerlist
}

func sortByDate(asc bool, date1, date2 time.Time) bool {
	if asc { // ASC
		if !date1.IsZero() && !date2.IsZero() {
			return date1.Before(date2)
		}
		return !date2.IsZero()
	}
	// DESC
	if !date1.IsZero() && !date2.IsZero() {
		return date1.After(date2)
	}
	return date2.IsZero()
}

func (timerlist *TimerList) sortByStartDate(order int) *TimerList {
	timerlist.sortBy(func(t1, t2 *Timer) bool {
		return sortByDate(order == SORT_START_DATE_ASC, t1.StartDate, t2.StartDate)
	})
	return timerlist
}

func (timerlist *TimerList) sortByFinishDate(order int) *TimerList {
	timerlist.sortBy(func(t1, t2 *Timer) bool {
		return sortByDate(order == SORT_FINISH_DATE_ASC, t1.FinishDate, t2.FinishDate)
	})
	return timerlist
}

func (timerlist *TimerList) sortByUnfinishedThenStart() *TimerList {
	timerlist.sortBy(func(t1, t2 *Timer) bool {
		if t1.FinishDate.IsZero() && !t2.FinishDate.IsZero() {
			return true
		} else if t2.FinishDate.IsZero() && !t1.FinishDate.IsZero() {
			return false
		}
		return sortByDate(false, t1.StartDate, t2.StartDate)
	})
	return timerlist
}
