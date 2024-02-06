package extloadtest

import (
	"fmt"
	"github.com/steadybit/discovery-kit/go/discovery_kit_api"
	"github.com/steadybit/discovery-kit/go/discovery_kit_commons"
	"github.com/steadybit/extension-kit/extutil"
	"github.com/steadybit/extension-loadtest/config"
)

func getDiscoveryKubernetesContainer() discovery_kit_api.DiscoveryDescription {
	return discovery_kit_api.DiscoveryDescription{
		Id: "com.steadybit.extension_kubernetes.kubernetes-container",
		Discover: discovery_kit_api.DescribingEndpointReferenceWithCallInterval{
			CallInterval: extutil.Ptr("1m"),
		},
	}
}

func createKubernetesContainerTargets(podTargets []discovery_kit_api.Target) []discovery_kit_api.EnrichmentData {
	count := len(podTargets) * config.Config.ContainerPerPod
	result := make([]discovery_kit_api.EnrichmentData, 0, count)

	for _, pod := range podTargets {
		for containerIndex := 1; containerIndex <= config.Config.ContainerPerPod; containerIndex++ {
			containerIdStripped := fmt.Sprintf("%s-c-%d", pod.Id, containerIndex)
			containerId := fmt.Sprintf("containerd://%s", containerIdStripped)
			target := discovery_kit_api.EnrichmentData{
				Id:                 containerId,
				EnrichmentDataType: "com.steadybit.extension_kubernetes.kubernetes-container",
				Attributes: map[string][]string{
					"k8s.cluster-name":                         {config.Config.ClusterName},
					"k8s.container.id":                         {containerId},
					"k8s.container.id.stripped":                {containerIdStripped},
					"k8s.container.image":                      {"docker.io/steadybit/loadtest-example:latest"},
					"k8s.container.name":                       {"loadtest-example"},
					"k8s.container.ready":                      {"true"},
					"k8s.deployment":                           pod.Attributes["k8s.deployment"],
					"k8s.distribution":                         {"kubernetes"},
					"k8s.label.domain":                         {"loadtest"},
					"k8s.label.run":                            {"loadtest-example"},
					"k8s.label.service-tier":                   {"2"},
					"k8s.label.tags.datadoghq.com/service":     {"loadtest"},
					"k8s.label.tags.datadoghq.com/version":     {"1.0.0"},
					"k8s.namespace":                            pod.Attributes["k8s.namespace"],
					"k8s.node.name":                            pod.Attributes["host.hostname"],
					"k8s.pod.label.domain":                     {"loadtest"},
					"k8s.pod.label.run":                        {"loadtest-example"},
					"k8s.pod.label.service-tier":               {"2"},
					"k8s.pod.label.tags.datadoghq.com/service": {"loadtest"},
					"k8s.pod.label.tags.datadoghq.com/version": {"1.0.0"},
					"k8s.pod.name":                             {pod.Id},
					"k8s.replicaset":                           pod.Attributes["k8s.replicaset"],
					"k8s.service.name":                         {fmt.Sprintf("%s-service", pod.Attributes["k8s.deployment"][0])},
				},
			}
			result = append(result, target)
		}
	}

	return discovery_kit_commons.ApplyAttributeExcludesToEnrichmentData(result, config.Config.DiscoveryAttributesExcludesKubernetesContainer)
}
