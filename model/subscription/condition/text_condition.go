package condition

type (
	TextCondition interface {
		KeyCondition
		GetTerm() string
	}

	textCondition struct {
		KeyCondition KeyCondition
		Term         string
	}
)

func NewTextCondition(kc KeyCondition, pattern string) TextCondition {
	return textCondition{
		KeyCondition: kc,
		Term:         pattern,
	}
}

func (kc textCondition) IsNot() bool {
	return kc.KeyCondition.IsNot()
}

func (kc textCondition) GetKey() string {
	return kc.KeyCondition.GetKey()
}

func (kc textCondition) GetTerm() string {
	return kc.Term
}
