package misc

import "time"

func TrackLatency(start time.Time) string {
	return time.Since(start).String()
}
