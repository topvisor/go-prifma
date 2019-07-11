package proxy

type Condition struct {
}

func NewCondition(condition string) (*Condition, error) {
	return new(Condition), nil
}
