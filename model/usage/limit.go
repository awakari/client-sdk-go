package usage

// Limit represents the usage limit.
type Limit struct {

	// Count represents the maximum count of usage number.
	Count int64

	// UserId represents the Limit user association. If empty, the Limit is a group-level limit.
	UserId string
}
