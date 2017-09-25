package main

import (
	"errors"
	"math"
	"regexp"
	"strconv"
	"time"
)

var (
	//ErrBadFormat bad duration or recurrence format
	ErrBadFormat = errors.New("bad duration or recurrence format")

	basicDateFormat    = "20060102T150405"
	extendedDateFormat = "2006-01-02T15:04:05"

	repeatRegexp = regexp.MustCompile(`^(R|R(\d+))$`)

	fullDurationRegexp = regexp.MustCompile(`P((?P<year>\d+)Y)?((?P<month>\d+)M)?((?P<day>\d+)D)?(T((?P<hour>\d+)H)?((?P<minute>\d+)M)?((?P<second>\d+)S)?)?`)
	weekDurationRegexp = regexp.MustCompile(`P(\d+)W`)
)

//DurationFromString parses duration from string
func DurationFromString(dur string) (*RecurrenceInterval, error) {
	d := &RecurrenceInterval{}

	if weekDurationRegexp.MatchString(dur) {
		match := weekDurationRegexp.FindStringSubmatch(dur)
		val, err := strconv.Atoi(match[1])

		if err != nil || val <= 0 {
			return nil, ErrBadFormat
		}

		d.Weeks = val

		return d, nil
	}

	var (
		match []string
		re    *regexp.Regexp
	)

	if fullDurationRegexp.MatchString(dur) {
		match = fullDurationRegexp.FindStringSubmatch(dur)
		re = fullDurationRegexp
	} else {
		return nil, ErrBadFormat
	}

	for i, name := range re.SubexpNames() {
		part := match[i]

		if i == 0 || name == "" || part == "" {
			continue
		}

		val, err := strconv.Atoi(part)

		if err != nil || val <= 0 {
			return nil, ErrBadFormat
		}

		switch name {
		case "year":
			d.Years = val
		case "month":
			d.Months = val
		case "week":
			d.Weeks = val
		case "day":
			d.Days = val
		case "hour":
			d.Hours = val
		case "minute":
			d.Minutes = val
		case "second":
			d.Seconds = val
		default:
			panic("unknown field " + name)
		}
	}

	if d.Years == 0 && d.Months == 0 && d.Weeks == 0 && d.Hours == 0 && d.Minutes == 0 && d.Seconds == 0 {
		return nil, ErrBadFormat
	}

	return d, nil
}

func RepeatFromString(repeatString string) (int, error) {
	if repeatRegexp.MatchString(repeatString) {
		match := repeatRegexp.FindStringSubmatch(repeatString)

		if match[1] == "R" { //unbounded repeats
			return -1, nil
		}

		val, err := strconv.Atoi(match[2])

		if err != nil || val < 0 {
			return math.MinInt32, ErrBadFormat
		}

		return val, nil
	} else {
		return math.MinInt32, ErrBadFormat
	}
}

func DateFromString(dateString string) (time.Time, error) {
	date, err := time.Parse(basicDateFormat, dateString)

	if err == nil {
		return date, nil
	}

	date, err = time.Parse(extendedDateFormat, dateString)

	if err == nil {
		return date, nil
	} else {
		return time.Time{}, ErrBadFormat
	}
}

//RecurrenceFromString parsing ISO8601 recurret intervals string
func RecurrenceFromString(recurrenceString string) (*Recurrence, error) {
	return nil, ErrBadFormat
}
