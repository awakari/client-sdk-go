package subscription

import "github.com/awakari/client-sdk-go/model/subscription/condition"

type Data struct {

	// Metadata represents a mutable Subscription data.
	Metadata Metadata

	// Condition represents the certain criteria to select the Subscription for the further routing.
	// It's immutable once the Subscription is created.
	Condition condition.Condition
}
