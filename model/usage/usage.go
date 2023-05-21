package usage

import "time"

type Usage struct {

	// Count represents the usage count since the last reset time (if any).
	Count int64

	// CountTotal represents the total usage count since the time in the Since field.
	CountTotal int64

	// Since represents the time when usage started for the 1st time. Together with CountTotal it may be used to
	// estimate the average usage rate.
	Since time.Time
}
