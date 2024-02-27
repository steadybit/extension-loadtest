package extloadtest

import (
	"fmt"
	"github.com/steadybit/discovery-kit/go/discovery_kit_api"
	"github.com/steadybit/discovery-kit/go/discovery_kit_commons"
	"github.com/steadybit/extension-kit/extutil"
	"github.com/steadybit/extension-loadtest/config"
)

func getDiscoveryKubernetesPods() discovery_kit_api.DiscoveryDescription {
	return discovery_kit_api.DiscoveryDescription{
		Id: "com.steadybit.extension_kubernetes.kubernetes-pod",
		Discover: discovery_kit_api.DescribingEndpointReferenceWithCallInterval{
			CallInterval: extutil.Ptr("1m"),
		},
	}
}

func createKubernetesPodTargets(hostTargets, deploymentTargets []discovery_kit_api.Target) []discovery_kit_api.Target {
	result := make([]discovery_kit_api.Target, 0, len(deploymentTargets)*config.Config.PodsPerDeployment)

	for hostIndex := 0; hostIndex < len(hostTargets); hostIndex++ {
		host := hostTargets[hostIndex]

		for deploymentIndex := 0; deploymentIndex < config.Config.DeploymentsPerNode; deploymentIndex++ {
			deployment := deploymentTargets[(hostIndex*config.Config.DeploymentsPerNode)+deploymentIndex]
			pods := deployment.Attributes["k8s.pod.name"]

			for podIndex := 0; podIndex < len(pods); podIndex++ {
				podName := pods[podIndex]

				containers := make([]string, 0, config.Config.ContainerPerPod)
				containersStripped := make([]string, 0, config.Config.ContainerPerPod)
				for containerIndex := 1; containerIndex <= config.Config.ContainerPerPod; containerIndex++ {
					containers = append(containers, fmt.Sprintf("containerd://%s-%s-c-%d", config.Config.PodUID, podName, containerIndex))
					containersStripped = append(containersStripped, fmt.Sprintf("%s-%s-c-%d", config.Config.PodUID, podName, containerIndex))
				}

				target := discovery_kit_api.Target{
					Id:         podName,
					TargetType: "com.steadybit.extension_kubernetes.kubernetes-pod",
					Label:      podName,
					Attributes: map[string][]string{
						"k8s.pod.name":                         {podName},
						"k8s.namespace":                        deployment.Attributes["k8s.namespace"],
						"k8s.cluster-name":                     {config.Config.ClusterName},
						"k8s.node.name":                        host.Attributes["host.hostname"],
						"host.hostname":                        host.Attributes["host.hostname"],
						"k8s.deployment":                       deployment.Attributes["k8s.deployment"],
						"k8s.label.domain":                     {"loadtest"},
						"k8s.label.run":                        {"loadtest-example"},
						"k8s.label.service-tier":               {"2"},
						"k8s.label.tags.datadoghq.com/service": {"loadtest"},
						"k8s.label.tags.datadoghq.com/version": {"1.0.0"},
						"k8s.replicaset":                       deployment.Attributes["k8s.deployment"],
						"k8s.container.id":                     containers,
						"k8s.container.id.stripped":            containersStripped,
					},
				}
				deployment.Attributes["host.hostname"] = append(deployment.Attributes["host.hostname"], host.Attributes["host.hostname"]...)
				result = append(result, target)
			}
		}
	}
	return discovery_kit_commons.ApplyAttributeExcludes(result, config.Config.DiscoveryAttributesExcludesKubernetesPod)
}
