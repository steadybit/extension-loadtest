// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2024 Steadybit GmbH

package extloadtest

import (
	"github.com/rs/zerolog/log"
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
	log.Info().Msgf("Enrichment Host to Container enabled=%t", config.Config.EnrichmentHostToContainerEnabled)
	result := make([]discovery_kit_api.TargetEnrichmentRule, 0)
	if config.Config.EnrichmentHostToContainerEnabled {
		result = append(result, getHostToContainerEnrichmentRule())
	}
	log.Info().Msgf("Enrichment Container to Host enabled=%t", config.Config.EnrichmentContainerToHostEnabled)
	if config.Config.EnrichmentContainerToHostEnabled {
		result = append(result, getContainerToHostEnrichmentRule())
	}
	log.Info().Msgf("Enrichment EC2 to Host enabled=%t", config.Config.EnrichmentEc2ToHostEnabled)
	if config.Config.EnrichmentEc2ToHostEnabled {
		result = append(result, getEC2ToHostEnrichmentRule())
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
func getEC2ToHostEnrichmentRule() discovery_kit_api.TargetEnrichmentRule {
	return discovery_kit_api.TargetEnrichmentRule{
		Id:      "com.steadybit.extension_loadtest.ec2-to-host",
		Version: extbuild.GetSemverVersionStringOrUnknown(),
		Src: discovery_kit_api.SourceOrDestination{
			Type: "com.steadybit.extension_aws.ec2-instance",
			Selector: map[string]string{
				"aws-ec2.hostname.internal": "${dest.host.hostname}",
			},
		},
		Dest: discovery_kit_api.SourceOrDestination{
			Type: "com.steadybit.extension_host.host",
			Selector: map[string]string{
				"host.hostname": "${src.aws-ec2.hostname.internal}",
			},
		},
		Attributes: []discovery_kit_api.Attribute{
			{
				Matcher: discovery_kit_api.Equals,
				Name:    "aws.account",
			}, {
				Matcher: discovery_kit_api.Equals,
				Name:    "aws.region",
			},
			{
				Matcher: discovery_kit_api.Equals,
				Name:    "aws.zone",
			},
			{
				Matcher: discovery_kit_api.Equals,
				Name:    "aws-ec2.arn",
			},
			{
				Matcher: discovery_kit_api.Equals,
				Name:    "aws-ec2.instance.id",
			},
			{
				Matcher: discovery_kit_api.Equals,
				Name:    "aws-ec2.instance.name",
			},
			{
				Matcher: discovery_kit_api.StartsWith,
				Name:    "aws-ec2.label.",
			},
		},
	}
}
