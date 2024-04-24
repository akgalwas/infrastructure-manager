package providers

//import (
//	"encoding/json"
//	gardener_azure "github.com/gardener/gardener-extension-provider-azure/pkg/apis/azure/v1alpha1"
//	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
//)
//
//func GetRawAzureInfrastructureConfig() (jsonData []byte) {
//	infraStructureConfigAzure := gardener_azure.InfrastructureConfig{
//		TypeMeta: v1.TypeMeta{
//			Kind:       "gcp.provider.extensions.gardener.cloud/v1alpha1",
//			APIVersion: "InfrastructureConfig",
//		},
//		Networks: gardener_azure.NetworkConfig{
//			Worker:  "10.180.0.0/16",
//			Workers: "10.180.0.0/16",
//		},
//	}
//	jsonData, _ = json.Marshal(infraStructureConfigAzure)
//	return jsonData
//}
