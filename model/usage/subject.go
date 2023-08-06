package usage

type Subject int

const (
	SubjectUndefined Subject = iota
	SubjectSubscriptions
	SubjectPublishEvents
)

func (s Subject) String() string {
	return [...]string{
		"SubjectUndefined",
		"SubjectSubscriptions",
		"SubjectPublishEvents",
	}[s]
}
