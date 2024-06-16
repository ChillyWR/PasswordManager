package pmtime

import "time"

func TruncateToMillisecond(t time.Time) time.Time {
	return t.Truncate(time.Millisecond)
}
