package shoot

import (
	"github.com/gardener/gardener-extension-shoot-dns-service/pkg/apis/service/v1alpha1"
	"github.com/kyma-project/infrastructure-manager/hack/shoot-comparator/pkg/runtime"
	"github.com/kyma-project/infrastructure-manager/hack/shoot-comparator/pkg/utilz"
	"github.com/onsi/gomega"
	"github.com/onsi/gomega/gstruct"
	"github.com/onsi/gomega/types"
	api_runtime "k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/yaml"
)

type providerCfgMatcher struct {
	toMatch        interface{} // *runtime.RawExtension
	pcType         string
	failed         string
	negativeFailed string
}

func newProviderCfgMatcher(pcType string, v interface{}) types.GomegaMatcher {
	return &providerCfgMatcher{
		toMatch: v,
		pcType:  pcType,
	}
}

func (m *providerCfgMatcher) NegatedFailureMessage(_ interface{}) string {
	return m.failed
}

func (m *providerCfgMatcher) FailureMessage(_ interface{}) string {
	return m.negativeFailed
}

func (m *providerCfgMatcher) Match(actual interface{}) (success bool, err error) {
	if actual == nil && m.toMatch == nil {
		return true, nil
	}

	actualProviderCfg, err := utilz.Get[*api_runtime.RawExtension](actual)
	if err != nil {
		return false, err
	}

	toMatchProviderCfg, err := utilz.Get[*api_runtime.RawExtension](m.toMatch)
	if err != nil {
		return false, err
	}

	if m.pcType != "shoot-dns-service" {
		return runtime.NewRawExtensionMatcher(m.toMatch).Match(actual)
	}

	var actualCfg v1alpha1.DNSConfig
	if err := yaml.Unmarshal(actualProviderCfg.Raw, &actualCfg); err != nil {
		return false, err
	}

	var toMatchCfg v1alpha1.DNSConfig
	if err := yaml.Unmarshal(toMatchProviderCfg.Raw, &toMatchCfg); err != nil {
		return false, err
	}

	matcher := gstruct.MatchFields(
		gstruct.IgnoreMissing,
		gstruct.Fields{
			"SyncProvidersFromShootSpecDNS": gomega.BeComparableTo(actualCfg.SyncProvidersFromShootSpecDNS),
			"DNSProviderReplication":        gomega.BeComparableTo(actualCfg.DNSProviderReplication),
			"Providers":                     gstruct.MatchAllElements(idDNSProvider, dnsProviders(actualCfg.Providers)),
			"TypeMeta":                      gstruct.Ignore(),
		})
	match, err := matcher.Match(toMatchCfg)
	if !match {
		m.failed = matcher.FailureMessage(actualCfg)
		m.negativeFailed = matcher.NegatedFailureMessage(actualCfg)
	}
	return match, err
}

func idDNSProvider(v interface{}) string {
	p, ok := v.(v1alpha1.DNSProvider)
	if !ok {
		panic("invalid type")
	}

	return *p.Type
}

func dnsProviders(ps []v1alpha1.DNSProvider) gstruct.Elements {
	out := map[string]types.GomegaMatcher{}
	for _, p := range ps {
		ID := idDNSProvider(p)
		secretNameMatcher := gomega.BeNil()
		if p.SecretName != nil {
			secretNameMatcher = gstruct.PointTo(gomega.HaveSuffix(*p.SecretName))
		}

		domainsMatcher := gomega.BeNil()
		if p.Domains != nil {
			domainsMatcher = gstruct.PointTo(gstruct.MatchFields(
				gstruct.IgnoreMissing,
				gstruct.Fields{
					"Include": gomega.BeComparableTo(p.Domains.Include),
					"Exclude": gstruct.Ignore(),
				}))
		}

		zonesMatcher := gomega.BeNil()
		if p.Zones != nil {
			zonesMatcher = gstruct.PointTo(gstruct.MatchFields(
				gstruct.IgnoreMissing,
				gstruct.Fields{
					"Include": gomega.BeComparableTo(p.Zones.Include),
					"Exclude": gstruct.Ignore(),
				}))
		}

		out[ID] = gstruct.MatchAllFields(gstruct.Fields{
			"Domains":    domainsMatcher,
			"SecretName": secretNameMatcher,
			"Type":       gomega.BeComparableTo(p.Type),
			"Zones":      zonesMatcher,
		})
	}
	return out
}
