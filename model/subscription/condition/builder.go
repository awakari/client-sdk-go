package condition

type Builder interface {
	Negation() Builder
	GroupLogic(l GroupLogic) Builder
	GroupChildren(children []Condition) Builder
	MatchAttrKey(k string) Builder
	MatchAttrValuePattern(p string) Builder
	MatchAttrValuePartial() Builder

	BuildGroupCondition() (c Condition)
	BuildKiwiTreeCondition() (c Condition)
}

type builder struct {
	not     bool
	gl      GroupLogic
	gc      []Condition
	key     string
	pattern string
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

func (b *builder) MatchAttrValuePattern(p string) Builder {
	b.pattern = p
	return b
}

func (b *builder) MatchAttrValuePartial() Builder {
	b.partial = true
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

func (b *builder) BuildKiwiTreeCondition() (c Condition) {
	c = condition{
		b.not,
	}
	c = kiwiTreeCondition{
		kiwiCondition{
			KeyCondition: keyCondition{
				Condition: c,
				Key:       b.key,
			},
			Partial: b.partial,
			Pattern: b.pattern,
		},
	}
	return
}
