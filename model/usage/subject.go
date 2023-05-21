package usage

type Subject int

const (
	SubjectUndefined Subject = iota
	SubjectPublishMessages
	SubjectSubscriptions
)

func (s Subject) String() string {
	return [...]string{
		"SubjectUndefined",
		"SubjectPublishMessages",
		"SubjectSubscriptions",
	}[s]
}
