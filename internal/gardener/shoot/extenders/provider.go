package extenders

import (
	gardener "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	imv1 "github.com/kyma-project/infrastructure-manager/api/v1"
	"github.com/kyma-project/infrastructure-manager/internal/gardener/shoot/hyperscaler/aws"
	"github.com/kyma-project/infrastructure-manager/internal/gardener/shoot/hyperscaler/azure"
	"github.com/kyma-project/infrastructure-manager/internal/gardener/shoot/hyperscaler/gcp"
	"github.com/kyma-project/infrastructure-manager/internal/gardener/shoot/hyperscaler/sapconvergedcloud"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime"
)

func ExtendWithProvider(runtimeShoot imv1.RuntimeShoot, shoot *gardener.Shoot) error {
	provider := &shoot.Spec.Provider
	provider.Type = runtimeShoot.Provider.Type
	provider.Workers = runtimeShoot.Provider.Workers

	var err error
	provider.InfrastructureConfig, provider.ControlPlaneConfig, err = getProviderSpecificConfig(runtimeShoot)

	return err
}

type InfrastructureProviderFunc func(workersCidr string, zones []string) ([]byte, error)
type ControlPlaneProviderFunc func() ([]byte, error)

func getProviderSpecificConfig(runtimeShoot imv1.RuntimeShoot) (infrastructureConfig *runtime.RawExtension, controlPlaneConfig *runtime.RawExtension, err error) {
	getConfig := func(runtimeShoot imv1.RuntimeShoot, infrastructureConfigFunc InfrastructureProviderFunc, controlPlaneConfigFunc ControlPlaneProviderFunc) (*runtime.RawExtension, *runtime.RawExtension, error) {
		infrastructureConfigBytes, err := infrastructureConfigFunc(runtimeShoot.Networking.Nodes, runtimeShoot.Provider.Workers[0].Zones)
		if err != nil {
			return nil, nil, err
		}

		controlPlaneConfigBytes, err := controlPlaneConfigFunc()
		if err != nil {
			return nil, nil, err
		}

		return &runtime.RawExtension{Raw: infrastructureConfigBytes}, &runtime.RawExtension{Raw: controlPlaneConfigBytes}, nil
	}

	switch runtimeShoot.Provider.Type {
	case "aws":
		{
			return getConfig(runtimeShoot, aws.GetInfrastructureConfig, aws.GetControlPlaneConfig)
		}
	case "azure":
		{
			return getConfig(runtimeShoot, azure.GetInfrastructureConfig, azure.GetControlPlaneConfig)
		}
	case "gcp":
		{
			return getConfig(runtimeShoot, gcp.GetInfrastructureConfig, gcp.GetControlPlaneConfig)
		}
	case "sapconvergedcloud":
		{
			return getConfig(runtimeShoot, sapconvergedcloud.GetInfrastructureConfig, sapconvergedcloud.GetControlPlaneConfig)
		}
	default:
		return nil, nil, errors.New("provider not supported")
	}
}