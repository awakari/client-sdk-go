package condition

type Builder interface {
	Negation() Builder
	GroupLogic(l GroupLogic) Builder
	GroupChildren(children []Condition) Builder
	MatchAttrKey(k string) Builder
	MatchText(p string) Builder

	BuildGroupCondition() (c Condition)
	BuildTextCondition() (c Condition)
}

type builder struct {
	not     bool
	gl      GroupLogic
	gc      []Condition
	key     string
	term    string
	partial bool
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
		Term: b.term,
	}
	return
}
