package fsm

import (
	"context"

	imv1 "github.com/kyma-project/infrastructure-manager/api/v1"
	gardener_shoot "github.com/kyma-project/infrastructure-manager/internal/gardener/shoot"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

func sFnCreateShoot(ctx context.Context, m *fsm, s *systemState) (stateFn, *ctrl.Result, error) {
	converterConfig := FixConverterConfig()
	converter := gardener_shoot.NewConverter(converterConfig)
	shoot, err := converter.ToShoot(s.instance)

	if err != nil {
		m.log.Error(err, "unable to convert Runtime CR to a shoot object")

		s.instance.UpdateStatePending(
			imv1.ConditionTypeRuntimeProvisioned,
			imv1.ConditionReasonConversionError,
			"False",
			"Runtime conversion error",
		)

		return updateStatusAndStop()
	}

	m.log.Info("Shoot converted successfully", "Name", shoot.Name, "Namespace", shoot.Namespace, "Shoot", shoot)

	s.shoot, err = m.ShootClient.Create(ctx, &shoot, v1.CreateOptions{})

	if err != nil {
		m.log.Error(err, "Failed to create new gardener Shoot")

		s.instance.UpdateStatePending(
			imv1.ConditionTypeRuntimeProvisioned,
			imv1.ConditionReasonGardenerError,
			"False",
			"Gardener API create error",
		)

		return updateStatusAndRequeueAfter(gardenerRequeueDuration)
	}

	m.log.Info("Gardener shoot for runtime initialised successfully", "Name", s.shoot.Name, "Namespace", s.shoot.Namespace)

	s.instance.UpdateStatePending(
		imv1.ConditionTypeRuntimeProvisioned,
		imv1.ConditionReasonShootCreationPending,
		"Unknown",
		"Shoot is pending",
	)

	shouldPersistShoot := m.PVCPath != ""
	if shouldPersistShoot {
		return switchState(sFnPersistShoot)
	}

	return updateStatusAndRequeueAfter(gardenerRequeueDuration)
}

func FixConverterConfig() gardener_shoot.ConverterConfig {
	return gardener_shoot.ConverterConfig{
		Kubernetes: gardener_shoot.KubernetesConfig{
			DefaultVersion: "1.29", //nolint:godox TODO: Should be parametrised
		},

		DNS: gardener_shoot.DNSConfig{
			SecretName:   "aws-route53-secret-dev",
			DomainPrefix: "dev.kyma.ondemand.com",
			ProviderType: "aws-route53",
		},
		Provider: gardener_shoot.ProviderConfig{
			AWS: gardener_shoot.AWSConfig{
				EnableIMDSv2: true, //nolint:godox TODO: Should be parametrised
			},
		},
		Gardener: gardener_shoot.GardenerConfig{
			ProjectName: "kyma-dev", //nolint:godox TODO: should be parametrised
		},
	}
}
