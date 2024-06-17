package fsm

import (
	"context"
	"github.com/pkg/errors"
	"slices"
)

type StateMachine struct {
	States          []State
	Entry           State
	EntryTransition Transition
	Transitions     map[State][]Transition
}

func NewStateMachine() StateMachine {
	return StateMachine{
		Transitions: map[State][]Transition{},
	}
}

func (sm *StateMachine) RegisterStates(state ...State) *StateMachine {
	sm.States = append(sm.States, state...)

	return sm
}

func (sm *StateMachine) SetEntry(state State) *StateMachine {
	sm.Entry = state

	return sm
}

func (sm *StateMachine) AddTransitions(_ string, transitions ...Transition) *StateMachine {

	for _, transition := range transitions {
		sm.AddTransition(transition)
	}

	return sm
}

func (sm *StateMachine) AddTransition(transition Transition) *StateMachine {
	sm.Transitions[transition.Current()] = append(sm.Transitions[transition.Current()], transition)

	return sm
}

func (sm *StateMachine) Run(ctx context.Context) (Result, error) {
	err := sm.validate()

	if err != nil {
		return ResultConfigurationError, err
	}

	state := sm.Entry

	for {
		select {
		case <-ctx.Done():
			err = ctx.Err()
			return ResultCancelled, err
		default:
			if state == Finished {
				return ResultFinished, nil
			}

			if state == Postponed {
				return ResultPostponed, nil
			}

			err = state.Do(ctx)
			if err != nil {
				return ResultError, err
			}

			state, err = sm.getNextState(state)
			if err != nil {
				return ResultError, nil
			}
		}
	}
}

func (sm *StateMachine) getNextState(currentState State) (State, error) {
	transitions, found := sm.Transitions[currentState]
	if !found {
		return Unknown, nil
	}

	for _, transition := range transitions {
		nextState, err := transition.Next()
		if err != nil {
			return Unknown, err
		}

		if nextState != Unknown {
			return nextState, nil
		}
	}

	return Unknown, nil
}

func (sm *StateMachine) validate() error {
	if sm.Entry == nil {
		errors.New("entry state not set")
	}

	var traversedStaes []State

	exitReached, err := sm.traverseStates(sm.Entry, traversedStaes)
	if err != nil {
		return err
	}

	if !exitReached {
		return errors.New("transitions graph doesn't reach exit state")
	}

	return nil
}

func (sm *StateMachine) traverseStates(currentState State, traversedStates []State) (bool, error) {
	if currentState == Finished || currentState == Postponed {
		return true, nil
	}
	if slices.Contains(traversedStates, currentState) {
		return false, errors.New("cycle detected")
	}

	traversedStates = append(traversedStates, currentState)

	transitions, found := sm.Transitions[currentState]
	if !found {
		return false, errors.New("state not found")
	}

	var exitReached bool
	for _, transition := range transitions {
		nextState, err := transition.Next()
		if err != nil {
			return false, errors.New("transition failed")
		}
		exitReached, err = sm.traverseStates(nextState, traversedStates)
		if err != nil {
			return false, err
		}
	}

	return exitReached, nil
}
