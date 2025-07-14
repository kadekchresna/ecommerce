package helper_time

import "time"

type TimeHelper struct {
	t *time.Time
}

func NewTime(t *time.Time) TimeHelper {
	return TimeHelper{
		t: t,
	}
}

func (t *TimeHelper) Now() time.Time {
	if t.t != nil && !t.t.IsZero() {
		return *t.t
	}

	return time.Now()
}

func IsWeekend(day time.Weekday) bool {
	return day == time.Saturday || day == time.Sunday
}
