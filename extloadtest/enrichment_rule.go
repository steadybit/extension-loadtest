package extloadtest

import (
	"github.com/steadybit/discovery-kit/go/discovery_kit_api"
	"github.com/steadybit/discovery-kit/go/discovery_kit_sdk"
	"github.com/steadybit/extension-kit/extbuild"
	"github.com/steadybit/extension-loadtest/config"
)

type ltEnrichmentRuleProvider struct {
}

var (
	_ discovery_kit_sdk.EnrichmentRulesDescriber = (*ltEnrichmentRuleProvider)(nil)
)

func NewEnrichmentRuleProvider() discovery_kit_sdk.EnrichmentRulesDescriber {
	return &ltEnrichmentRuleProvider{}
}

func (d *ltEnrichmentRuleProvider) DescribeEnrichmentRules() []discovery_kit_api.TargetEnrichmentRule {
	result := make([]discovery_kit_api.TargetEnrichmentRule, 0)
	if config.Config.EnrichmentHostToContainerEnabled {
		result = append(result, getHostToContainerEnrichmentRule())
	}
	if config.Config.EnrichmentContainerToHostEnabled {
		result = append(result, getContainerToHostEnrichmentRule())
	}
	return result
}

func getHostToContainerEnrichmentRule() discovery_kit_api.TargetEnrichmentRule {
	return discovery_kit_api.TargetEnrichmentRule{
		Id:      "com.steadybit.extension_loadtest.host-to-container",
		Version: extbuild.GetSemverVersionStringOrUnknown(),
		Src: discovery_kit_api.SourceOrDestination{
			Type: "com.steadybit.extension_host.host",
			Selector: map[string]string{
				"host.hostname": "${dest.host.hostname}",
			},
		},
		Dest: discovery_kit_api.SourceOrDestination{
			Type: "com.steadybit.extension_container.container",
			Selector: map[string]string{
				"host.hostname": "${src.host.hostname}",
			},
		},
		Attributes: []discovery_kit_api.Attribute{
			{
				Matcher: discovery_kit_api.StartsWith,
				Name:    "host.os.",
			},
		},
	}
}

func getContainerToHostEnrichmentRule() discovery_kit_api.TargetEnrichmentRule {
	return discovery_kit_api.TargetEnrichmentRule{
		Id:      "com.steadybit.extension_loadtest.container-to-host",
		Version: extbuild.GetSemverVersionStringOrUnknown(),
		Src: discovery_kit_api.SourceOrDestination{
			Type: "com.steadybit.extension_container.container",
			Selector: map[string]string{
				"host.hostname": "${dest.host.hostname}",
			},
		},
		Dest: discovery_kit_api.SourceOrDestination{
			Type: "com.steadybit.extension_host.host",
			Selector: map[string]string{
				"host.hostname": "${src.host.hostname}",
			},
		},
		Attributes: []discovery_kit_api.Attribute{
			{
				Matcher: discovery_kit_api.Equals,
				Name:    "container.id",
			},
		},
	}
}
