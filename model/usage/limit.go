package usage

import "time"

// Limit represents the usage limit.
type Limit struct {

	// Count represents the maximum count of usage number.
	Count int64

	// UserId represents the Limit user association. If empty, the Limit is a group-level limit.
	UserId string

	// Expires represents the user-specific limit expiration deadline.
	Expires time.Time
}
