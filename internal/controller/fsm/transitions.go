package fsm

type State interface {
	Do() error
}

type Transition interface {
	Current() State
	Next() (State, error)
}

func Immediate(from, to State) Transition {
	return immediateTransition{
		From: from,
		To:   to,
	}
}

type Predicate interface {
	True() (bool, error)
}

type immediateTransition struct {
	From State
	To   State
}

func (it immediateTransition) Current() State {
	return it.From
}

func (it immediateTransition) Next() (State, error) {
	return it.From, nil
}

func Conditional(from, to State, predicate Predicate) Transition {
	return conditionalTransition{
		from:      from,
		To:        to,
		predicate: predicate,
	}
}

type conditionalTransition struct {
	from      State
	To        State
	predicate Predicate
}

func (it conditionalTransition) Current() State {
	return it.from
}

func (it conditionalTransition) Next() (State, error) {
	predicateTrue, err := it.predicate.True()

	if err != nil {
		return Unknown, err
	}

	if predicateTrue {
		return it.To, nil
	}

	return Unknown, nil
}
