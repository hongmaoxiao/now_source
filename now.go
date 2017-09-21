package now

import "time"

func (now *Now) BeginningOfMinute() time.Time {
	return now.time.Truncate(time.Minute)
}

func (now *Now) BeginningOfHour() time.Time {
	return now.time.Truncate(time.Hour)
}

func (now *Now) BeginningOfDay() time.Time {
	d := time.Duration(-now.time.Hour()) * time.Hour
	return now.time.Truncate(time.Hour).Add(d)
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
