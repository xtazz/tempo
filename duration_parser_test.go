package main

import (
	"testing"
	"time"
)

func TestDurationFromString(t *testing.T) {
	t.Parallel()

	if _, err := DurationFromString("asdf"); err != ErrBadFormat {
		t.Errorf("Expected %s but got %s", ErrBadFormat, err)
	}

	// test without params
	if _, err := DurationFromString("PYMDTHMS"); err != ErrBadFormat {
		t.Errorf("Expected %s but got %#v", ErrBadFormat, err)
	}

	// test with good full string
	dur, err := DurationFromString("P1Y2M3DT4H5M6S")

	if err != nil {
		t.Errorf("Expected no error but got %s", err)
	}

	if dur.Years != 1 {
		t.Errorf("Expected year equal 1 but got %#v", dur.Years)
	}

	if dur.Months != 2 {
		t.Errorf("Expected month equal 2 but got %#v", dur.Months)
	}

	if dur.Days != 3 {
		t.Errorf("Expected days equal 3 but got %#v", dur.Days)
	}

	if dur.Hours != 4 {
		t.Errorf("Expected hours equal 4 but got %#v", dur.Hours)
	}

	if dur.Minutes != 5 {
		t.Errorf("Expected minutes equal 5 but got %#v", dur.Minutes)
	}

	if dur.Seconds != 6 {
		t.Errorf("Expected seconds equal 6 but got %#v", dur.Seconds)
	}

	if dur.Weeks != 0 {
		t.Errorf("Expected weeks equal 0 but got %#v", dur.Weeks)
	}

	// test with good week string
	dur, err = DurationFromString("P1W")

	if err != nil {
		t.Errorf("Expected no error but got %s", err)
	}

	if dur.Weeks != 1 {
		t.Errorf("Expected weeks equal 1 but got %#v", dur.Weeks)
	}

	// test with bad week string
	if _, err := DurationFromString("PW"); err != ErrBadFormat {
		t.Errorf("Expected %s but got %s", ErrBadFormat, err)
	}
}

func TestRepeatFromString(t *testing.T) {
	t.Parallel()

	repeat, err := RepeatFromString("R")

	if err != nil {
		t.Errorf("Expected correct format but got error %#v", err)
	} else if repeat != -1 {
		t.Errorf("Expected repeat to be -1 (unbounded) but got %#v", repeat)
	}

	repeat, err = RepeatFromString("R0")

	if err != nil {
		t.Errorf("Expected correct format but got error %#v", err)
	} else if repeat != 0 {
		t.Errorf("Expected repeat to be 0 (no repetitions) but got %#v", repeat)
	}

	repeat, err = RepeatFromString("R-1")

	if err == nil || err != ErrBadFormat {
		t.Errorf("Expected invalid format but got error %#v", err)
	}

	repeat, err = RepeatFromString("R999999")

	if err != nil {
		t.Errorf("Expected correct format but got error %#v", err)
	} else if repeat != 999999 {
		t.Errorf("Expected repeat to be 1 but got %#v", repeat)
	}
}

func TestDateFromString(t *testing.T) {
	t.Parallel()

	var formatTests = []struct {
		DateString    string
		Expected      string
		FormatCorrect bool
	}{
		/* correct formats */
		{"19850412T232050", "1985-04-12T23:20:50Z", true},     //basic
		{"1985-04-12T23:20:50", "1985-04-12T23:20:50Z", true}, //extended
		{"2020-02-29T23:20:50", "2020-02-29T23:20:50Z", true}, //leap year,
		/* incorrect formats */
		{"19850412232050", "", false},      //no T separator
		{"1985-31-12T23:20:50", "", false}, //month is more than 12
		{"1985-04-31T23:20:50", "", false}, //day is more than days in month
		{"2019-02-29T23:20:50", "", false}, //non leap year but 29 days in Feb
		{"1985-04-12T25:60:60", "", false}, //time is incorrect
	}

	for index, test := range formatTests {
		date, err := DateFromString(test.DateString)

		if test.FormatCorrect && err != nil {
			t.Errorf("Expected correct format but got error %#v", err)
		} else if !test.FormatCorrect {
			if err == nil {
				t.Errorf("Expected incorrect format but got %#v", date)
			} else if err != ErrBadFormat {
				t.Errorf("Expected incorrect format error but got error %#v", err)
			}
		} else {
			result := date.Format(time.RFC3339)

			if result != test.Expected {
				t.Errorf("Test %d expected %s but got %s", index+1, test.Expected, result)
			}
		}
	}
}
