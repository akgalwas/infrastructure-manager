package fsm

type StateMachine struct {
	States          []State
	Entry           State
	EntryTransition Transition
	Transitions     map[State][]Transition
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

func (sm *StateMachine) Run() (State, error) {
	state := sm.Entry
	err := error(nil)

	for {
		if state == Finished || state == Postponed || state == Unknown {
			break
		}

		err = state.Do()
		if err != nil {
			break
		}

		state, err = sm.getNextState(state)
		if err != nil {
			break
		}
	}

	return state, err
}

func (sm *StateMachine) getNextState(currentState State) (State, error) {
	transitions, found := sm.Transitions[currentState]
	if !found {
		return Unknown, nil
	}

	for _, transition := range transitions {
		nextState, err := transition.Next()
		if err != nil {
			return Error, err
		}

		if nextState != Unknown {
			return nextState, nil
		}
	}

	return Unknown, nil
}
