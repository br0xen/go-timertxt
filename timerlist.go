package timertxt

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// TimerList represents a list of timer.txt timer entries.
// It is usually loasded from a whole timer.txt file.
type TimerList []Timer

// NewTimerList creates a new empty TimerList.
func NewTimerList() *TimerList {
	return &TimerList{}
}

func (timerlist *TimerList) GetTimersInRange(start, end time.Time) *TimerList {
	t := *NewTimerList()
	for _, v := range *timerlist {
		if v.FinishDate.Before(end) && 
	}
	return &t
}

func (timerlist *TimerList) GetActiveTimers() *TimerList {
	t := *NewTimerList()
	for _, v := range *timerlist {
		if v.FinishDate.IsZero() {
			t = append(t, v)
		}
	}
	return &t
}

// String returns a complete list of timers in timer.txt format.
func (timerlist *TimerList) String() string {
	var ret string
	for _, timer := range *timerlist {
		ret += fmt.Sprintf("%s\n", timer.String())
	}
	return ret
}

// AddTimer appends a Timer to the current TimerList and takes care to set the Timer.Id correctly
func (timerlist *TimerList) AddTimer(timer *Timer) {
	timer.Id = 0
	for _, t := range *timerlist {
		if t.Id > timer.Id {
			timer.Id = t.Id
		}
	}
	timer.Id += 1
	*timerlist = append(*timerlist, *timer)
}

// GetTimer returns the Timer with the given timer 'id' from the TimerList.
// Returns an error if Timer could not be found.
func (timerlist *TimerList) GetTimer(id int) (*Timer, error) {
	for i := range *timerlist {
		if ([]Timer(*timerlist))[i].Id == id {
			return &([]Timer(*timerlist))[i], nil
		}
	}
	return nil, errors.New("timer not found")
}

// RemoveTimerById removes any Timer with given Timer 'id' from the TimerList.
// Returns an error if no Timer was removed.
func (timerlist *TimerList) RemoveTimerById(id int) error {
	var newList TimerList
	found := false
	for _, t := range *timerlist {
		if t.Id != id {
			newList = append(newList, t)
		} else {
			found = true
		}
	}
	if !found {
		return errors.New("timer not found")
	}
	*timerlist = newList
	return nil
}

// RemoveTimer removes any Timer from the TimerList with the same String representation as the given Timer.
// Returns an error if no Timer was removed.
func (timerlist *TimerList) RemoveTimer(timer Timer) error {
	var newList TimerList
	found := false
	for _, t := range *timerlist {
		if t.String() != timer.String() {
			newList = append(newList, t)
		} else {
			found = true
		}
	}
	if !found {
		return errors.New("timer not found")
	}
	*timerlist = newList
	return nil
}

// Filter filters the current TimerList for the given predicate (a function that takes a timer as input and returns a
// bool), and returns a new TimerList. The original TimerList is not modified.
func (timerlist *TimerList) Filter(predicate func(Timer) bool) *TimerList {
	var newList TimerList
	for _, t := range *timerlist {
		if predicate(t) {
			newList = append(newList, t)
		}
	}
	return &newList
}

// LoadFromFile loads a TimerList from *os.File.
// Note: This will clear the current TimerList and overwrite it's contents with whatever is in *os.File.
func (timerlist *TimerList) LoadFromFile(file *os.File) error {
	*timerlist = []Timer{} // Empty timerlist
	timerId := 1
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := strings.Trim(scanner.Text(), "\t\n\r") // Read Line
		// Ignore blank lines
		if text == "" {
			continue
		}
		timer, err := ParseTimer(text)
		if err != nil {
			return err
		}
		timer.Id = timerId
		*timerlist = append(*timerlist, *timer)
		timerId++
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

// WriteToFile writes a TimerList to *os.File
func (timerlist *TimerList) WriteToFile(file *os.File) error {
	writer := bufio.NewWriter(file)
	_, err := writer.WriteString(timerlist.String())
	writer.Flush()
	return err
}

// WriteToFile writes a TimerList to *os.File.
func (timerlist *TimerList) LoadFromFilename(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	return timerlist.LoadFromFile(file)
}

// WriteToFilename writes a TimerList to the specified file (most likely called "timer.txt").
func (timerlist *TimerList) WriteToFilename(filename string) error {
	return ioutil.WriteFile(filename, []byte(timerlist.String()), 0640)
}

// LoadFromFile loads and returns a TimerList from *os.File.
func LoadFromFile(file *os.File) (TimerList, error) {
	timerlist := TimerList{}
	if err := timerlist.LoadFromFile(file); err != nil {
		return nil, err
	}
	return timerlist, nil
}

// WriteToFile writes a TimerList to *os.File.
func WriteToFile(timerlist *TimerList, file *os.File) error {
	return timerlist.WriteToFile(file)
}

// LoadFromFilename loads and returns a TimerList from a file (most likely called "timer.txt")
func LoadFromFilename(filename string) (TimerList, error) {
	timerlist := TimerList{}
	if err := timerlist.LoadFromFilename(filename); err != nil {
		return nil, err
	}
	return timerlist, nil
}

// WriteToFilename write a TimerList to the specified file (most likely called "timer.txt")
func WriteToFilename(timerlist *TimerList, filename string) error {
	return timerlist.WriteToFilename(filename)
}
