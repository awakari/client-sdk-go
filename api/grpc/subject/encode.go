package subject

import (
	"errors"
	"fmt"
	"github.com/awakari/client-sdk-go/model/usage"
)

var ErrInvalidSubject = errors.New("unrecognized subject")

func Encode(src usage.Subject) (dst Subject, err error) {
	switch src {
	case usage.SubjectPublishMessages:
		dst = Subject_PublishMessages
	case usage.SubjectSubscriptions:
		dst = Subject_Subscriptions
	default:
		err = fmt.Errorf("%w: %s", ErrInvalidSubject, src)
	}
	return
}
