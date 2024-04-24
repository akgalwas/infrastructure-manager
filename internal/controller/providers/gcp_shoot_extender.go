package providers

import (
	"encoding/json"
	gardener_gcp "github.com/gardener/gardener-extension-provider-gcp/pkg/apis/gcp/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const gcpAPIVersion = "gcp.provider.extensions.gardener.cloud/v1alpha1"

func GetRawGCPInfrastructureConfig() (jsonData []byte) {
	infraStructureConfigGCP := gardener_gcp.InfrastructureConfig{
		TypeMeta: v1.TypeMeta{
			Kind:       "InfrastructureConfig",
			APIVersion: gcpAPIVersion,
		},
		Networks: gardener_gcp.NetworkConfig{
			Workers: "10.180.0.0/16",
		},
	}
	jsonData, _ = json.Marshal(infraStructureConfigGCP)
	return jsonData
}

func GetGCPControlPlane() (jsonData []byte) {
	controlPlaneConfig := &gardener_gcp.ControlPlaneConfig{
		TypeMeta: v1.TypeMeta{
			Kind:       "ControlPlaneConfig",
			APIVersion: gcpAPIVersion,
		},
		Zone: "europe-central2-a",
	}

	jsonData, _ = json.Marshal(controlPlaneConfig)
	return jsonData
}

//admission webhook "validator.admission-gcp.extensions.gardener.cloud" denied the request: error decoding controlPlaneConfig: strict decoding error: unknown field "Zone", unknown field "CloudControllerManager", unknown field "Storage"
