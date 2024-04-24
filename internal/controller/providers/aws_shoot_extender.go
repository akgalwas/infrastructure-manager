package providers

import (
	"encoding/json"
	gardener_aws "github.com/gardener/gardener-extension-provider-aws/pkg/apis/aws/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetRawAWSInfrastructureConfig() (jsonData []byte) {
	zones := make([]gardener_aws.Zone, 0)
	zone := gardener_aws.Zone{
		Name:     "europe-central2-a",
		Internal: "10.180.48.0/20",
		Public:   "10.180.32.0/20",
		Workers:  "10.180.0.0/19",
	}
	//        internal: 10.250.112.0/22
	//        public: 10.250.96.0/22
	//        workers: 10.250.0.0/19
	//

	//
	zones = append(zones, zone)
	infrastructureConfigAWS := gardener_aws.InfrastructureConfig{
		TypeMeta: v1.TypeMeta{
			Kind:       "InfrastructureConfig",
			APIVersion: "aws.provider.extensions.gardener.cloud/v1alpha1",
		},
		EnableECRAccess: PtrTo(true),
		DualStack: &gardener_aws.DualStack{
			Enabled: false,
		},
		Networks: gardener_aws.Networks{
			VPC: gardener_aws.VPC{
				CIDR: PtrTo("10.180.0.0/16"),
			},
			Zones: zones,
		},
	}
	jsonData, _ = json.Marshal(infrastructureConfigAWS)
	return jsonData
}

func PtrTo[T any](v T) *T {
	return &v
}
