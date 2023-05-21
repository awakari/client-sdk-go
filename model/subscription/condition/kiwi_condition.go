package condition

type (
	KiwiCondition interface {
		KeyCondition
		IsPartial() bool
		GetPattern() string
	}

	kiwiCondition struct {
		KeyCondition KeyCondition
		Partial      bool
		Pattern      string
	}
)

func NewKiwiCondition(kc KeyCondition, partial bool, pattern string) KiwiCondition {
	return kiwiCondition{
		KeyCondition: kc,
		Partial:      partial,
		Pattern:      pattern,
	}
}

func (kc kiwiCondition) IsNot() bool {
	return kc.KeyCondition.IsNot()
}

func (kc kiwiCondition) GetKey() string {
	return kc.KeyCondition.GetKey()
}

func (kc kiwiCondition) IsPartial() bool {
	return kc.Partial
}

func (kc kiwiCondition) GetPattern() string {
	return kc.Pattern
}
