package main

import (
	"testing"
	"time"
)

var now = time.Date(2017, 1, 1, 0, 0, 0, 0, time.UTC)

func TestNextDate(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		Interval RecurrenceInterval
		Expected string
	}{
		{RecurrenceInterval{Hours: 24}, "2017-01-02T00:00:00Z"},
		{RecurrenceInterval{Hours: 23, Minutes: 60}, "2017-01-02T00:00:00Z"},
		{RecurrenceInterval{Hours: 23, Minutes: 59, Seconds: 60}, "2017-01-02T00:00:00Z"},
		{RecurrenceInterval{Weeks: 4, Days: 3}, "2017-02-01T00:00:00Z"},
		{RecurrenceInterval{Months: 12}, "2018-01-01T00:00:00Z"},
		{RecurrenceInterval{Months: 3}, "2017-04-01T00:00:00Z"},
		{RecurrenceInterval{Months: 3, Days: 7}, "2017-04-08T00:00:00Z"},
		{RecurrenceInterval{Months: 12, Days: 30, Hours: 23, Minutes: 60, Seconds: 1}, "2018-02-01T00:00:01Z"},
		{RecurrenceInterval{Months: 1, Days: 28}, "2017-03-01T00:00:00Z"},
		{RecurrenceInterval{Years: 1}, "2018-01-01T00:00:00Z"},
		{RecurrenceInterval{Years: 3, Months: 1, Days: 28}, "2020-02-29T00:00:00Z"}, //leap year test
		{RecurrenceInterval{Years: 3, Months: 1, Weeks: 4}, "2020-02-29T00:00:00Z"}, //leap year test
	}

	for index, test := range tests {
		nextDate := test.Interval.NextDate(now)
		result := nextDate.Format(time.RFC3339)

		if result != test.Expected {
			t.Errorf("Test %d expected %s but got %s", index+1, test.Expected, result)
		}
	}
}
