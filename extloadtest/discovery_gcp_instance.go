package extloadtest

import (
	"fmt"
	"github.com/steadybit/discovery-kit/go/discovery_kit_api"
	"github.com/steadybit/discovery-kit/go/discovery_kit_commons"
	"github.com/steadybit/extension-kit/extutil"
	"github.com/steadybit/extension-loadtest/config"
)

func getDiscoveryGcpInstance() discovery_kit_api.DiscoveryDescription {
	return discovery_kit_api.DiscoveryDescription{
		Id: "com.steadybit.extension_gcp.vm",
		Discover: discovery_kit_api.DescribingEndpointReferenceWithCallInterval{
			CallInterval: extutil.Ptr("30s"),
		},
	}
}

func createGcpInstanceTargets(hosts []discovery_kit_api.Target) []discovery_kit_api.Target {
	result := make([]discovery_kit_api.Target, 0, len(hosts))
	for i, host := range hosts {
		instanceId := fmt.Sprintf("i-%s", host.Id)
		instanceName := fmt.Sprintf("loadtest-instance-%s", host.Id)
		hostname := host.Id
		zones := []string{"europe-west1-a"}
		if i%2 == 0 {
			zones = []string{"europe-west1-a", "europe-west1-b"}
		}
		target := discovery_kit_api.Target{
			Id:         instanceId,
			TargetType: "com.steadybit.extension_gcp.vm",
			Label:      instanceName,
			Attributes: map[string][]string{
				"gcp-vm.name":        {instanceName},
				"gcp-vm.id":          {instanceId},
				"gcp-vm.hostname":    {hostname},
				"gcp-vm.description": {"loadtest"},

				"gcp-vm.cpu-platform":                    {"Intel Broadwell"},
				"gcp-vm.machine-type":                    {"n1-standard-1"},
				"gcp-vm.source-machine-image":            {"projects/debian-cloud/global/images/family/debian-11"},
				"gcp-vm.status":                          {"RUNNING"},
				"gcp-vm.status-message":                  {"RUNNING"},
				"gcp.zone-url":                           {"https://www.googleapis.com/compute/v1/projects/steadybit/zones/europe-west1-b"},
				"gcp.zone":                               zones,
				"gcp.project.id":                         {"steadybit-loadtest"},
				"gcp-kubernetes-engine.cluster.name":     {config.Config.ClusterName},
				"gcp-kubernetes-engine.cluster.location": zones,
			},
		}
		result = append(result, target)
	}

	return discovery_kit_commons.ApplyAttributeExcludes(result, config.Config.DiscoveryAttributesExcludesEc2)
}
