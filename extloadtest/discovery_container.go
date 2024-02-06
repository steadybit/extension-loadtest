package extloadtest

import (
	"github.com/steadybit/discovery-kit/go/discovery_kit_api"
	"github.com/steadybit/discovery-kit/go/discovery_kit_commons"
	"github.com/steadybit/extension-kit/extutil"
	"github.com/steadybit/extension-loadtest/config"
)

func getDiscoveryContainer() discovery_kit_api.DiscoveryDescription {
	return discovery_kit_api.DiscoveryDescription{
		Id: "com.steadybit.extension_container.container",
		Discover: discovery_kit_api.DescribingEndpointReferenceWithCallInterval{
			CallInterval: extutil.Ptr("1m"),
		},
	}
}

func createContainerTargets(kubernetesContainers []discovery_kit_api.EnrichmentData) []discovery_kit_api.Target {
	result := make([]discovery_kit_api.Target, 0, len(kubernetesContainers))

	for _, kubernetesContainer := range kubernetesContainers {
		target := discovery_kit_api.Target{
			Id:         kubernetesContainer.Id,
			TargetType: "com.steadybit.extension_container.container",
			Label:      kubernetesContainer.Id,
			Attributes: map[string][]string{
				"container.engine":                       {"containerd"},
				"container.engine.version":               {"1.6.6"},
				"container.host":                         kubernetesContainer.Attributes["k8s.node.name"],
				"host.hostname":                          kubernetesContainer.Attributes["k8s.node.name"],
				"container.id":                           kubernetesContainer.Attributes["k8s.container.id"],
				"container.id.stripped":                  kubernetesContainer.Attributes["k8s.container.id.stripped"],
				"container.image":                        {"docker.io/steadybit/LoadTest:latest"},
				"container.image.registry":               {"docker.io"},
				"container.image.repository":             {"docker.io/steadybit/LoadTest"},
				"container.image.tag":                    {"latest"},
				"container.label.io.cri-containerd.kind": {"container"},
				"container.label.io.kubernetes.pod.uid":  {"6418b03c-147c-4685-854b-9ffc324216f2"},
				"k8s.container.name":                     kubernetesContainer.Attributes["k8s.container.name"],
				"k8s.namespace":                          kubernetesContainer.Attributes["k8s.namespace"],
				"k8s.pod.name":                           kubernetesContainer.Attributes["k8s.pod.name"],
			},
		}
		result = append(result, target)
	}

	return discovery_kit_commons.ApplyAttributeExcludes(result, config.Config.DiscoveryAttributesExcludesContainer)
}
