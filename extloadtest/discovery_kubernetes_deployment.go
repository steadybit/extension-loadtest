package extloadtest

import (
	"fmt"
	"github.com/steadybit/discovery-kit/go/discovery_kit_api"
	"github.com/steadybit/discovery-kit/go/discovery_kit_commons"
	"github.com/steadybit/extension-kit/extutil"
	"github.com/steadybit/extension-loadtest/config"
)

func getDiscoveryKubernetesDeployment() discovery_kit_api.DiscoveryDescription {
	return discovery_kit_api.DiscoveryDescription{
		Id: "com.steadybit.extension_kubernetes.kubernetes-deployment",
		Discover: discovery_kit_api.DescribingEndpointReferenceWithCallInterval{
			CallInterval: extutil.Ptr("1m"),
		},
	}
}

func createKubernetesDeploymentTargets(nodeCount int, suffix string) []discovery_kit_api.Target {
	count := nodeCount * config.Config.DeploymentsPerNode
	result := make([]discovery_kit_api.Target, 0, count)
	for i := 1; i <= count; i++ {
		namespace := "loadtest-namespace"
		deploymentWithOutPodUid := fmt.Sprintf("d-%d-%s", i, suffix)
		deployment := fmt.Sprintf("%s-%s", config.Config.PodUID, deploymentWithOutPodUid)
		deploymentLabel := fmt.Sprintf("%s-loadtest-deployment-%d-%s", config.Config.PodUID, i, suffix)
		id := fmt.Sprintf("%s/%s/%s/%s", config.Config.PodUID, config.Config.ClusterName, namespace, deploymentWithOutPodUid)

		pods := make([]string, 0, config.Config.PodsPerDeployment)
		containerCount := config.Config.PodsPerDeployment * config.Config.ContainerPerPod
		containers := make([]string, 0, containerCount)
		containersStripped := make([]string, 0, containerCount)
		for podIndex := 1; podIndex <= config.Config.PodsPerDeployment; podIndex++ {
			podName := fmt.Sprintf("%s-%s-Pod-%d", config.Config.PodUID, deployment, podIndex)
			pods = append(pods, podName)
			for containerIndex := 1; containerIndex <= config.Config.ContainerPerPod; containerIndex++ {
				containers = append(containers, fmt.Sprintf("containerd://%s-c-%d", podName, containerIndex))
				containersStripped = append(containersStripped, fmt.Sprintf("%s-c-%d", podName, containerIndex))
			}
		}

		target := discovery_kit_api.Target{
			Id:         id,
			TargetType: "com.steadybit.extension_kubernetes.kubernetes-deployment",
			Label:      deploymentLabel,
			Attributes: map[string][]string{
				"k8s.cluster-name":                                {config.Config.ClusterName},
				"k8s.container.id":                                containers,
				"k8s.container.id.stripped":                       containersStripped,
				"k8s.deployment":                                  {deployment},
				"k8s.deployment.label.domain":                     {"shop-products"},
				"k8s.deployment.label.run":                        {"loadtest"},
				"k8s.deployment.label.service-tier":               {"2"},
				"k8s.deployment.label.tags.datadoghq.com/service": {"shop-products"},
				"k8s.deployment.label.tags.datadoghq.com/version": {"1.0.0"},
				"k8s.distribution":                                {"kubernetes"},
				"k8s.label.domain":                                {"shop-products"},
				"k8s.label.run":                                   {"loadtest"},
				"k8s.label.service-tier":                          {"2"},
				"k8s.label.tags.datadoghq.com/service":            {"shop-products"},
				"k8s.label.tags.datadoghq.com/version":            {"1.0.0"},
				"k8s.namespace":                                   {namespace},
				"k8s.pod.name":                                    pods,
			},
		}
		result = append(result, target)
	}
	return discovery_kit_commons.ApplyAttributeExcludes(result, config.Config.DiscoveryAttributesExcludesKubernetesDeployment)
}
