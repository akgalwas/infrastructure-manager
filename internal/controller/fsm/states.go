package fsm

var Unknown = SpecialState{}
var Finished = SpecialState{}
var Postponed = SpecialState{}
var Error = SpecialState{}

type SpecialState struct {
}

func (su SpecialState) Do() error {
	return nil
}

type ErrorOcurredState struct {
}

func (su ErrorOcurredState) Do() error {
	return nil
}
