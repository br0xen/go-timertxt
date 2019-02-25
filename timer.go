package timertxt

import (
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"
)

var (
	// DateLayout is used for formatting time.Time into timer.txt date format and vice-versa.
	DateLayout = time.RFC3339

	addonTagRx = regexp.MustCompile(`(^|\s+)([\w-]+):(\S+)`) // Match additional tags date: '... due:2012-12-12 ...'
	contextRx  = regexp.MustCompile(`(^|\s+)@(\S+)`)         // Match contexts: '@Context ...' or '... @Context ...'
	projectRx  = regexp.MustCompile(`(^|\s+)\+(\S+)`)        // Match projects: '+Project...' or '... +Project ...')
)

type Timer struct {
	Id             int    // Internal timer id
	Original       string // Original raw timer text
	StartDate      time.Time
	FinishDate     time.Time
	Finished       bool
	Notes          string // Notes part of timer text
	Projects       []string
	Contexts       []string
	AdditionalTags map[string]string // Addon tags will be available here
}

// String returns a complete timer string in timer.txt format.
//
// Contexts, Projects, and additional tags are alphabetically sorted,
// and appended at the end in the following order:
// Contexts, Projects, Tags
//
// For example:
// "2019-02-15T11:43:00-0600 Working on Go Library @home @personal +timertxt customTag1:Important! due:Today"
// "x 2019-02-15T10:00:00-0600 2019-02-15T06:00:00-0600 Creating Go Library Repo @home @personal +timertxt customTag1:Important! due:Today"
func (timer Timer) String() string {
	var text string
	if timer.Finished {
		text += "x "
		if !timer.FinishDate.IsZero() {
			text += fmt.Sprintf("%s ", timer.FinishDate.Format(DateLayout))
		}
	}
	text += fmt.Sprintf("%s ", timer.StartDate.Format(DateLayout))
	text += timer.Notes
	if len(timer.Contexts) > 0 {
		sort.Strings(timer.Contexts)
		for _, context := range timer.Contexts {
			text += fmt.Sprintf(" @%s", context)
		}
	}
	if len(timer.Projects) > 0 {
		sort.Strings(timer.Projects)
		for _, project := range timer.Projects {
			text += fmt.Sprintf(" +%s", project)
		}
	}
	if len(timer.AdditionalTags) > 0 {
		// Sort map alphabetically by keys
		keys := make([]string, 0, len(timer.AdditionalTags))
		for key := range timer.AdditionalTags {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		for _, key := range keys {
			text += fmt.Sprintf(" %s:%s", key, timer.AdditionalTags[key])
		}
	}
	return text
}

// NewTimer creates a new empty Timer with default values. (StartDate is set to Now())
func NewTimer() *Timer {
	timer := Timer{}
	timer.StartDate = time.Now()
	return &timer
}

// ParseTimer parses the input text string into a Timer struct
func ParseTimer(text string) (*Timer, error) {
	var err error
	timer := Timer{}
	timer.Original = strings.Trim(text, "\t\n\r ")
	originalParts := strings.Fields(timer.Original)

	// Check for finished
	if originalParts[0] == "x" {
		timer.Finished = true
		// If it's finished, there _must_ be a finished date
		if timer.FinishDate, err = time.Parse(DateLayout, originalParts[1]); err != nil {
			return nil, errors.New("Timer marked finished, but failed to parse FinishDate: " + err.Error())
		}
		originalParts = originalParts[2:]
	}
	if timer.StartDate, err = time.Parse(DateLayout, originalParts[0]); err != nil {
		return nil, errors.New("Unable to parse StartDate: " + err.Error())
	}
	originalParts = originalParts[1:]
	var notes []string
	for _, v := range originalParts {
		if strings.HasPrefix(v, "@") {
			v = strings.TrimPrefix(v, "@")
			// Contexts
			timer.Contexts = append(timer.Contexts, v)
		} else if strings.HasPrefix(v, "+") {
			// Projects
			v = strings.TrimPrefix(v, "+")
			timer.Projects = append(timer.Projects, v)
		} else if strings.Contains(v, ":") {
			// Additional tags
			tagPts := strings.Split(v, ":")
			if tagPts[0] != "" && tagPts[1] != "" {
				timer.AdditionalTags[tagPts[0]] = tagPts[1]
			}
		} else {
			notes = append(notes, v)
		}
	}
	timer.Notes = strings.Join(notes, " ")

	return &timer, nil
}

// Timer returns a complete timer string in timer.txt format.
// See *Timer.String() for further information
func (timer *Timer) Timer() string {
	return timer.String()
}

// Finish sets Timer.Finished to true if the timer hasn't already been finished.
// Also sets Timer.FinishDate to time.Now()
func (timer *Timer) Finish() {
	if !timer.Finished {
		timer.Finished = true
		timer.FinishDate = time.Now()
	}
}

// Reopen sets Timer.Finished to 'false' if the timer was finished
// Also resets Timer.FinishDate
func (timer *Timer) Reopen() {
	if timer.Finished {
		timer.Finished = false
		timer.FinishDate = time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC) // time.IsZero() value
	}
}

func (timer *Timer) Duration() time.Duration {
	end := time.Now()
	if !timer.FinishDate.IsZero() {
		end = timer.FinishDate
	}
	return end.Sub(timer.StartDate)
}

func (timer *Timer) ActiveToday() bool {
	return timer.ActiveOnDay(time.Now())
}

func (timer *Timer) ActiveOnDay(t time.Time) bool {
	f := "2006/01/02"
	tStr := t.Format(f)
	// If StartDate or FinishDate is _on_ t, true
	if timer.StartDate.Format(f) == tStr || timer.FinishDate.Format(f) == tStr {
		return true
	}
	// Otherwise, if StartDate is before t and FinishDate is after t
	return timer.StartDate.Before(t) && timer.FinishDate.After(t)
}

func (timer *Timer) HasContext(context string) bool {
	for _, v := range timer.Contexts {
		if v == context {
			return true
		}
	}
	return false
}

func (timer *Timer) HasProject(project string) bool {
	for _, v := range timer.Projects {
		if v == project {
			return true
		}
	}
	return false
}
