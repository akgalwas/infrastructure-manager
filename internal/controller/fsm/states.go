package fsm

import "context"

var Unknown = SpecialState{}
var Finished = SpecialState{}
var Postponed = SpecialState{}

type SpecialState struct {
}

func (su SpecialState) Do(ctx context.Context) error {
	return nil
}

type ErrorOcurredState struct {
}

func (su ErrorOcurredState) Do() error {
	return nil
}
