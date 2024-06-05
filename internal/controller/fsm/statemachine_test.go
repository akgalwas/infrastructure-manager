package fsm

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestStateMachine(t *testing.T) {
	// given

	// States
	init := DummyState{}
	createShoot := DummyState{}
	setupCluster := DummyState{}
	upgradeShoot := DummyState{}
	deprovisionShoot := DummyState{}

	// Predicates
	clusterDoesnExists := fixAlwaysTruePredicate()
	provisioningInProgress := fixAlwaysTruePredicate()
	provisioningCompleted := fixAlwaysTruePredicate()
	clusterToBeDeleted := fixAlwaysTruePredicate()
	clusterDeletionInProgress := fixAlwaysTruePredicate()
	clusterToBeUpgraded := fixAlwaysTruePredicate()
	clusterDeleted := fixAlwaysTruePredicate()
	clusterUpgradeInProgress := fixAlwaysTruePredicate()
	clusterUpgraded := fixAlwaysTruePredicate()

	// when
	stateMachine := StateMachine{}

	finalState, err := stateMachine.
		RegisterStates(
			init,
			createShoot,
			setupCluster,
			upgradeShoot,
			deprovisionShoot).
		SetEntry(init).
		AddTransitions("provisioning",
			Conditional(init, createShoot, clusterDoesnExists),
			Conditional(init, Postponed, provisioningInProgress),
			Conditional(init, setupCluster, provisioningCompleted),
			Immediate(createShoot, Postponed),
			Immediate(setupCluster, Finished),
		).
		AddTransitions("deprovisioning",
			Conditional(init, deprovisionShoot, clusterToBeDeleted),
			Conditional(init, Postponed, clusterDeletionInProgress),
			Conditional(init, Finished, clusterDeleted),
			Immediate(deprovisionShoot, Postponed),
		).
		AddTransitions("upgrade",
			Conditional(init, upgradeShoot, clusterToBeUpgraded),
			Conditional(init, Postponed, clusterUpgradeInProgress),
			Conditional(init, setupCluster, clusterUpgraded),
			Immediate(setupCluster, Finished),
			Immediate(upgradeShoot, Postponed),
		).
		Run()

	require.NoError(t, err)
	assert.Equal(t, Finished, finalState)
}

type DummyState struct {
}

func (su DummyState) Do() error {
	return nil
}

func fixAlwaysTruePredicate() Predicate {
	return AlwaysTruePredicate{}
}

type AlwaysTruePredicate struct {
}

func (a AlwaysTruePredicate) True() (bool, error) {
	return true, nil
}
