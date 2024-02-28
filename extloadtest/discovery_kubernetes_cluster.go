package extloadtest

import (
	"fmt"
	"github.com/steadybit/discovery-kit/go/discovery_kit_api"
	"github.com/steadybit/extension-kit/extutil"
	"github.com/steadybit/extension-loadtest/config"
)

func getDiscoveryKubernetesCluster() discovery_kit_api.DiscoveryDescription {
	return discovery_kit_api.DiscoveryDescription{
		Id: "com.steadybit.extension_kubernetes.kubernetes-cluster",
		Discover: discovery_kit_api.DescribingEndpointReferenceWithCallInterval{
			CallInterval: extutil.Ptr("60m"),
		},
	}
}

func createKubernetesClusterTargets() []discovery_kit_api.Target {
	id := fmt.Sprintf("%s-%s", config.Config.ClusterName, config.Config.PodUID)
	return []discovery_kit_api.Target{
		{
			Id:         id,
			TargetType: "com.steadybit.extension_kubernetes.kubernetes-cluster",
			Label:      config.Config.ClusterName,
			Attributes: map[string][]string{
				"k8s.cluster-name": {config.Config.ClusterName},
			},
		},
	}
}
