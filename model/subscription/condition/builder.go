package condition

type Builder interface {

	// Negation makes a resulting Condition negative.
	// May be set for any resulting Condition type.
	Negation() Builder

	// GroupLogic defines the logic for a resulting GroupCondition. Default is GroupLogicAnd.
	GroupLogic(l GroupLogic) Builder

	// GroupChildren defines the nested Condition set for a resulting GroupCondition.
	GroupChildren(children []Condition) Builder

	// BuildGroupCondition builds a GroupCondition.
	BuildGroupCondition() (c Condition)

	// MatchAttrKey defines the incoming messages attribute key to match for a resulting KeyCondition.
	// Default is empty string key. For a TextCondition empty key causes the matching against all attribute keys.
	MatchAttrKey(k string) Builder

	// MatchText defines the text search terms for a resulting TextCondition.
	MatchText(p string) Builder

	// MatchExact enables the exact text matching criteria for a resulting TextCondition.
	MatchExact() Builder

	// BuildTextCondition builds a TextCondition.
	BuildTextCondition() (c Condition)
}

type builder struct {
	not   bool
	gl    GroupLogic
	gc    []Condition
	key   string
	term  string
	exact bool
}

func NewBuilder() Builder {
	return &builder{}
}

func (b *builder) Negation() Builder {
	b.not = true
	return b
}

func (b *builder) GroupLogic(l GroupLogic) Builder {
	b.gl = l
	return b
}

func (b *builder) GroupChildren(children []Condition) Builder {
	b.gc = children
	return b
}

func (b *builder) MatchAttrKey(k string) Builder {
	b.key = k
	return b
}

func (b *builder) MatchText(term string) Builder {
	b.term = term
	return b
}

func (b *builder) MatchExact() Builder {
	b.exact = true
	return b
}

func (b *builder) BuildGroupCondition() (c Condition) {
	c = condition{
		Not: b.not,
	}
	c = groupCondition{
		Condition: c,
		Logic:     b.gl,
		Group:     b.gc,
	}
	return
}

func (b *builder) BuildTextCondition() (c Condition) {
	c = condition{
		b.not,
	}
	c = textCondition{
		KeyCondition: keyCondition{
			Condition: c,
			Key:       b.key,
		},
		Term:  b.term,
		Exact: b.exact,
	}
	return
}
