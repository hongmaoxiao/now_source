package now

import (
	"errors"
	"fmt"
	"regexp"
	"time"
)

func (now *Now) BeginningOfMinute() time.Time {
	return now.Truncate(time.Minute)
}

func (now *Now) BeginningOfHour() time.Time {
	return now.Truncate(time.Hour)
}

func (now *Now) BeginningOfDay() time.Time {
	d := time.Duration(-now.Hour()) * time.Hour
	return now.Truncate(time.Hour).Add(d)
}

func (now *Now) BeginningOfWeek() time.Time {
	t := now.BeginningOfDay()
	weekday := int(t.Weekday())
	if FirstDayMonday {
		if weekday == 0 {
			weekday = 7
		}
		weekday = weekday - 1
	}

	d := time.Duration(-weekday) * 24 * time.Hour
	return t.Truncate(time.Hour).Add(d)
}

func (now *Now) BeginningOfMonth() time.Time {
	t := now.BeginningOfDay()
	d := time.Duration(-int(t.Day())+1) * 24 * time.Hour
	return t.Truncate(time.Hour).Add(d)
}

func (now *Now) BeginningOfYear() time.Time {
	t := now.BeginningOfDay()
	d := time.Duration(-int(t.YearDay())+1) * 24 * time.Hour
	return t.Truncate(time.Hour).Add(d)
}

func (now *Now) EndOfMinute() time.Time {
	return now.BeginningOfMinute().Add(time.Minute - time.Nanosecond)
}

func (now *Now) EndOfHour() time.Time {
	return now.BeginningOfHour().Add(time.Hour - time.Nanosecond)
}

func (now *Now) EndOfDay() time.Time {
	return now.BeginningOfDay().Add(24*time.Hour - time.Nanosecond)
}

func (now *Now) EndOfWeek() time.Time {
	return now.BeginningOfWeek().AddDate(0, 0, 7).Add(-time.Nanosecond)
}

func (now *Now) EndOfMonth() time.Time {
	return now.BeginningOfMonth().AddDate(0, 1, 0).Add(-time.Nanosecond)
}

func (now *Now) EndOfYear() time.Time {
	return now.BeginningOfYear().AddDate(1, 0, 0).Add(-time.Nanosecond)
}

func (now *Now) Monday() time.Time {
	t := now.BeginningOfDay()
	weekday := int(t.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	d := time.Duration(-weekday+1) * 24 * time.Hour
	return t.Truncate(time.Hour).Add(d)
}

func (now *Now) Sunday() time.Time {
	t := now.BeginningOfDay()
	weekday := int(t.Weekday())
	if weekday == 0 {
		return t
	} else {
		d := time.Duration(7-weekday) * 24 * time.Hour
		return t.Truncate(time.Hour).Add(d)
	}
}

func (now *Now) EndOfSunday() time.Time {
	return now.Sunday().Add(24*time.Hour - time.Nanosecond)
}

func parseWithFormat(str string) (t time.Time, err error) {
	formats := []string{"2006-01-02 15:04:05", "2006-01-02 15:04", "2006-01-02", "01-02", "15:04:05", "15:04", "15"}
	for _, format := range formats {
		t, err = time.Parse(format, str)
		if err == nil {
			return
		}
	}
	err = errors.New("Can't parse string: " + str)
	return
}

func (now *Now) Parse(strs ...string) (t time.Time, err error) {
	var setCurrentTime bool
	parseTime := []int{}
	currentTime := []int{now.Second(), now.Minute(), now.Hour(), now.Day(), int(now.Month()), now.Year()}
	currentLocation := now.Location()

	for _, str := range strs {
		hasDate := regexp.MustCompile(`-\d`).MatchString(str)
		t, err = parseWithFormat(str)
		if err == nil {
			fmt.Println("t", t)
			parseTime = []int{t.Second(), t.Minute(), t.Hour(), t.Day(), int(t.Month()), t.Year()}

			fmt.Println("parseTime", parseTime)
			fmt.Println("currentTime", currentTime)
			for i, v := range parseTime {
				// Fill up missed information with current time
				if v == 0 {
					fmt.Println("parseTime[i]", parseTime[i])
					fmt.Println("currentTime", currentTime[i])
					if setCurrentTime {
						parseTime[i] = currentTime[i]
						fmt.Println("after parseTime[i]", parseTime[i])
						fmt.Println("after currentTime[i]", currentTime[i])
					}
				} else {
					setCurrentTime = true
				}

				// Default day and month is 1, fill up it if missing it
				if (i == 3 || i == 4) && !hasDate {
					parseTime[i] = currentTime[i]
					fmt.Println("after day parseTime[i]", parseTime[i])
					fmt.Println("after day currentTime[i]", currentTime[i])
				}
			}
		}

		fmt.Println("parseTime llll", parseTime)
		if len(parseTime) > 0 {
			fmt.Println("parseTime llll", parseTime)
			t = time.Date(parseTime[5], time.Month(parseTime[4]), parseTime[3], parseTime[2], parseTime[1], parseTime[0], 0, currentLocation)
			currentTime := []int{t.Second(), t.Minute(), t.Hour(), t.Day(), int(t.Month()), t.Year()}
			fmt.Println("parseTime t", t)
		}
	}
	return
}

func (now *Now) MustParse(strs ...string) (t time.Time) {
	t, err := now.Parse(strs...)
	if err != nil {
		panic(err)
	}
	return t
}
