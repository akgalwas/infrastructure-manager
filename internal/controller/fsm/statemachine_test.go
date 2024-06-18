package fsm

import (
	"context"
	"github.com/kyma-project/infrastructure-manager/internal/controller/fsm/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestStateMachine(t *testing.T) {
	t.Run("Test simple unconditional transitions", func(t *testing.T) {
		// given
		firstState := &mocks.State{}
		secondState := &mocks.State{}
		thirdState := &mocks.State{}

		firstState.On("Do", mock.Anything).Return(nil)
		secondState.On("Do", mock.Anything).Return(nil)
		thirdState.On("Do", mock.Anything).Return(nil)

		// when
		sm := NewStateMachine()

		result, err := sm.
			RegisterStates(firstState, secondState, thirdState).
			SetEntry(firstState).
			AddTransition(Immediate(firstState, secondState)).
			AddTransition(Immediate(secondState, thirdState)).
			AddTransition(Immediate(thirdState, Finished)).
			Run(context.TODO())

		// then
		require.NoError(t, err)
		assert.Equal(t, ResultFinished, result)
		firstState.AssertExpectations(t)
		secondState.AssertExpectations(t)
		thirdState.AssertExpectations(t)
	})

	t.Run("Test simple conditional transitions", func(t *testing.T) {
		// given
		firstState := &mocks.State{}
		secondState := &mocks.State{}
		thirdState := &mocks.State{}
		fourthState := &mocks.State{}

		alwaysTruePredicate := &mocks.Predicate{}
		alwaysFalsePredicate := &mocks.Predicate{}

		firstState.On("Do", mock.Anything).Return(nil)
		secondState.On("Do", mock.Anything).Return(nil)
		thirdState.On("Do", mock.Anything).Return(nil)
		alwaysTruePredicate.On("True").Return(true, nil)
		alwaysFalsePredicate.On("True").Return(false, nil)
		//when
		sm := NewStateMachine()

		result, err := sm.
			RegisterStates(firstState, secondState, thirdState, fourthState).
			SetEntry(firstState).
			AddTransition(Conditional(firstState, secondState, alwaysFalsePredicate)).
			AddTransition(Conditional(firstState, thirdState, alwaysTruePredicate)).
			AddTransition(Conditional(thirdState, Finished, alwaysTruePredicate)).
			AddTransition(Conditional(thirdState, fourthState, alwaysFalsePredicate)).
			AddTransition(Immediate(secondState, Finished)).
			AddTransition(Immediate(fourthState, Finished)).
			Run(context.TODO())

		//then
		require.NoError(t, err)
		assert.Equal(t, ResultFinished, result)
		firstState.AssertExpectations(t)
		secondState.AssertNotCalled(t, "Do")
		thirdState.AssertExpectations(t)
		fourthState.AssertNotCalled(t, "Do")
	})

	t.Run("Detect cycle", func(t *testing.T) {
		// given
		firstState := &mocks.State{}
		secondState := &mocks.State{}
		thirdState := &mocks.State{}

		firstState.On("Do", mock.Anything).Return(nil)
		secondState.On("Do", mock.Anything).Return(nil)
		thirdState.On("Do", mock.Anything).Return(nil)

		// when
		sm := NewStateMachine()

		result, err := sm.
			RegisterStates(firstState, secondState, thirdState).
			SetEntry(firstState).
			AddTransition(Immediate(firstState, secondState)).
			AddTransition(Immediate(secondState, thirdState)).
			AddTransition(Immediate(thirdState, firstState)).
			Run(context.Background())

		// then
		require.Error(t, err)
		assert.Equal(t, ResultConfigurationError, result)
	})
}

//func TestStateMachine2(t *testing.T) {
//	// given
//
//	// States
//	init := DummyState{}
//	createShoot := DummyState{}
//	setupCluster := DummyState{}
//	upgradeShoot := DummyState{}
//	deprovisionShoot := DummyState{}
//
//	// Predicates
//	clusterDoesnExists := fixAlwaysTruePredicate()
//	provisioningInProgress := fixAlwaysTruePredicate()
//	provisioningCompleted := fixAlwaysTruePredicate()
//	clusterToBeDeleted := fixAlwaysTruePredicate()
//	clusterDeletionInProgress := fixAlwaysTruePredicate()
//	clusterToBeUpgraded := fixAlwaysTruePredicate()
//	clusterDeleted := fixAlwaysTruePredicate()
//	clusterUpgradeInProgress := fixAlwaysTruePredicate()
//	clusterUpgraded := fixAlwaysTruePredicate()
//
//	// when
//	stateMachine := StateMachine{}
//
//	finalState, err := stateMachine.
//		RegisterStates(
//			init,
//			createShoot,
//			setupCluster,
//			upgradeShoot,
//			deprovisionShoot).
//		SetEntry(init).
//		AddTransitions("provisioning",
//			Conditional(init, createShoot, clusterDoesnExists),
//			Conditional(init, Postponed, provisioningInProgress),
//			Conditional(init, setupCluster, provisioningCompleted),
//			Immediate(createShoot, Postponed),
//			Immediate(setupCluster, Finished),
//		).
//		AddTransitions("deprovisioning",
//			Conditional(init, deprovisionShoot, clusterToBeDeleted),
//			Conditional(init, Postponed, clusterDeletionInProgress),
//			Conditional(init, Finished, clusterDeleted),
//			Immediate(deprovisionShoot, Postponed),
//		).
//		AddTransitions("upgrade",
//			Conditional(init, upgradeShoot, clusterToBeUpgraded),
//			Conditional(init, Postponed, clusterUpgradeInProgress),
//			Conditional(init, setupCluster, clusterUpgraded),
//			Immediate(setupCluster, Finished),
//			Immediate(upgradeShoot, Postponed),
//		).
//		Run()
//
//	require.NoError(t, err)
//	assert.Equal(t, Finished, finalState)
//}

//type DummyState struct {
//}
//
//func (su DummyState) Do() error {
//	return nil
//}
//
//func fixAlwaysTruePredicate() Predicate {
//	return AlwaysTruePredicate{}
//}
//
//type AlwaysTruePredicate struct {
//}
//
//func (a AlwaysTruePredicate) True() (bool, error) {
//	return true, nil
//}
