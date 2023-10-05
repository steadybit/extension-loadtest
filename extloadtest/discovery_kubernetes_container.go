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
		Id:         "com.steadybit.extension_kubernetes.kubernetes-container",
		RestrictTo: extutil.Ptr(discovery_kit_api.LEADER),
		Discover: discovery_kit_api.DescribingEndpointReferenceWithCallInterval{
			Method:       "GET",
			Path:         "/discovery/kubernetes-container/targets",
			CallInterval: extutil.Ptr("1m"),
		},
	}
}

func createKubernetesContainerTargets(hostTargets, deploymentTargets []discovery_kit_api.Target) []discovery_kit_api.EnrichmentData {
	count := len(deploymentTargets) * config.Config.PodsPerDeployment * config.Config.ContainerPerPod
	result := make([]discovery_kit_api.EnrichmentData, 0, count)

	for hostIndex := 0; hostIndex < len(hostTargets); hostIndex++ {
		host := hostTargets[hostIndex]

		for deploymentIndex := 0; deploymentIndex < config.Config.DeploymentsPerNode; deploymentIndex++ {
			deployment := deploymentTargets[(hostIndex*config.Config.DeploymentsPerNode)+deploymentIndex]
			pods := deployment.Attributes["k8s.pod.name"]
			containerIdsStripped := deployment.Attributes["k8s.container.id.stripped"]

			for podIndex := 0; podIndex < len(pods); podIndex++ {
				podName := pods[podIndex]

				for containerIndex := 0; containerIndex < config.Config.ContainerPerPod; containerIndex++ {
					containerIdStripped := containerIdsStripped[(podIndex*config.Config.ContainerPerPod)+containerIndex]
					containerId := fmt.Sprintf("dummy://%s", containerIdStripped)
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
							"k8s.deployment":                           deployment.Attributes["k8s.deployment"],
							"k8s.distribution":                         {"kubernetes"},
							"k8s.label.domain":                         {"loadtest"},
							"k8s.label.run":                            {"loadtest-example"},
							"k8s.label.service-tier":                   {"2"},
							"k8s.label.tags.datadoghq.com/service":     {"loadtest"},
							"k8s.label.tags.datadoghq.com/version":     {"1.0.0"},
							"k8s.namespace":                            deployment.Attributes["k8s.namespace"],
							"k8s.node.name":                            host.Attributes["host.hostname"],
							"k8s.pod.label.domain":                     {"loadtest"},
							"k8s.pod.label.run":                        {"loadtest-example"},
							"k8s.pod.label.service-tier":               {"2"},
							"k8s.pod.label.tags.datadoghq.com/service": {"loadtest"},
							"k8s.pod.label.tags.datadoghq.com/version": {"1.0.0"},
							"k8s.pod.name":                             {podName},
							"k8s.replicaset":                           deployment.Attributes["k8s.deployment"],
							"k8s.service.name":                         {fmt.Sprintf("%s-service", deployment.Attributes["k8s.deployment"][0])},
						},
					}
					result = append(result, target)
				}
			}
		}
	}

	return discovery_kit_commons.ApplyAttributeExcludesToEnrichmentData(result, config.Config.DiscoveryAttributesExcludesKubernetesContainer)
}
