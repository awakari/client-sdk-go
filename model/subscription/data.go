package subscription

import "github.com/awakari/client-sdk-go/model/subscription/condition"

type Data struct {

	// Condition represents the certain criteria to select the Subscription for the further routing.
	// It's immutable once the Subscription is created.
	Condition condition.Condition

	// Description is a human-readable subscription description.
	Description string

	// Enabled defines whether subscription is active and may be used to deliver a message.
	Enabled bool
}
