package extender

import (
	gardener "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	imv1 "github.com/kyma-project/infrastructure-manager/api/v1"
)

// Provisioner was setting the following annotations:
//- kcp.provisioner.kyma-project.io/licence-type
//- kcp.provisioner.kyma-project.io/runtime-id
//- support.gardener.cloud/eu-access-for-cluster-nodes

// support.gardener.cloud/eu-access-for-cluster-nodes is no longer set as it was added temporarily, and is no longer used by the Gardener

// Provisioner was setting the following labels:
//- accout
//- subaccount

const (
	ShootRuntimeIDAnnotation   = "infrastructuremanager.kyma-project.io/runtime-id"
	ShootLicenceTypeAnnotation = "infrastructuremanager.kyma-project.io/licence-type"
	ShootGlobalAccountLabel    = "account"
	ShootSubAccountLabel       = "subaccount"
	RuntimeIDLabel             = "kyma-project.io/runtime-id"
	RuntimeGlobalAccountLabel  = "kyma-project.io/global-account-id"
	RuntimeSubaccountLabel     = "kyma-project.io/subaccount-id"
)

func ExtendWithAnnotationsAndLabels(runtime imv1.Runtime, shoot *gardener.Shoot) error {
	shoot.Labels = getLabels(runtime)
	shoot.Annotations = getAnnotations(runtime)

	return nil
}

func getAnnotations(runtime imv1.Runtime) map[string]string {
	annotations := map[string]string{
		ShootRuntimeIDAnnotation: runtime.Labels[RuntimeIDLabel],
	}

	if runtime.Spec.Shoot.LicenceType != nil && *runtime.Spec.Shoot.LicenceType != "" {
		annotations[ShootLicenceTypeAnnotation] = *runtime.Spec.Shoot.LicenceType
	}

	return annotations
}

func getLabels(runtime imv1.Runtime) map[string]string {
	labels := map[string]string{
		ShootGlobalAccountLabel: runtime.Labels[RuntimeGlobalAccountLabel],
		ShootSubAccountLabel:    runtime.Labels[RuntimeSubaccountLabel],
	}

	return labels
}