package main

import "time"

type Recurrence struct {
	Repetitions int
	Start       *time.Time
	End         *time.Time
	Duration    *RecurrenceInterval
}

type RecurrenceInterval struct {
	Years   int
	Months  int
	Weeks   int
	Days    int
	Hours   int
	Minutes int
	Seconds int
}

// NextDate returns the next time by adding recurrence interval.
func (d RecurrenceInterval) NextDate(fromDate time.Time) time.Time {
	nextTime := fromDate.AddDate(d.Years, d.Months, d.Weeks*7+d.Days)

	nextTime = nextTime.Add(time.Duration(d.Hours) * time.Hour)
	nextTime = nextTime.Add(time.Duration(d.Minutes) * time.Minute)
	nextTime = nextTime.Add(time.Duration(d.Seconds) * time.Second)

	return nextTime
}
