package extloadtest

import (
	"fmt"
	"github.com/steadybit/discovery-kit/go/discovery_kit_api"
	"github.com/steadybit/extension-kit/exthttp"
	"github.com/steadybit/extension-kit/extutil"
	"github.com/steadybit/extension-loadtest/config"
)

func RegisterDiscoveryKubernetesDeployment() {
	exthttp.RegisterHttpHandler("/discovery/kubernetes-deployment", exthttp.GetterAsHandler(getDiscoveryKubernetesDeployment))
	exthttp.RegisterHttpHandler("/discovery/kubernetes-deployment/targets", exthttp.GetterAsHandler(getDiscoveryKubernetesDeploymentTargets))
}

func getDiscoveryKubernetesDeployment() discovery_kit_api.DiscoveryDescription {
	return discovery_kit_api.DiscoveryDescription{
		Id:         "com.steadybit.extension_kubernetes.kubernetes-deployment",
		RestrictTo: extutil.Ptr(discovery_kit_api.LEADER),
		Discover: discovery_kit_api.DescribingEndpointReferenceWithCallInterval{
			Method:       "GET",
			Path:         "/discovery/kubernetes-deployment/targets",
			CallInterval: extutil.Ptr("1m"),
		},
	}
}

var kubernetesDeployments []discovery_kit_api.Target

func getDiscoveryKubernetesDeploymentTargets() discovery_kit_api.DiscoveryData {
	if kubernetesDeployments == nil {
		kubernetesDeployments = initKubernetesDeploymentTargets()
	}
	return discovery_kit_api.DiscoveryData{
		Targets: &kubernetesDeployments,
	}
}

func initKubernetesDeploymentTargets() []discovery_kit_api.Target {
	count := config.Config.NodeCount * config.Config.DeploymentsPerNode
	result := make([]discovery_kit_api.Target, 0, count)
	for i := 1; i <= count; i++ {
		namespace := "loadtest-namespace"
		deployment := fmt.Sprintf("d-%d", i)
		deploymentLabel := fmt.Sprintf("loadtest-deployment-%d", i)
		id := fmt.Sprintf("%s/%s/%s", config.Config.ClusterName, namespace, deployment)

		pods := make([]string, 0, config.Config.PodsPerDeployment)
		containerCount := config.Config.PodsPerDeployment * config.Config.ContainerPerPod
		containers := make([]string, 0, containerCount)
		containersStripped := make([]string, 0, containerCount)
		for podIndex := 1; podIndex <= config.Config.PodsPerDeployment; podIndex++ {
			podName := fmt.Sprintf("%s-p-%d", deployment, podIndex)
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
	return discovery_kit_api.ApplyAttributeExcludes(result, config.Config.DiscoveryAttributeExcludesKubernetesDeployment)
}
