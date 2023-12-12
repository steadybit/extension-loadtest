package extloadtest

import (
	"fmt"
	"github.com/steadybit/discovery-kit/go/discovery_kit_api"
	"github.com/steadybit/discovery-kit/go/discovery_kit_commons"
	"github.com/steadybit/extension-kit/extutil"
	"github.com/steadybit/extension-loadtest/config"
)

func getDiscoveryAzureInstance() discovery_kit_api.DiscoveryDescription {
	return discovery_kit_api.DiscoveryDescription{
		Id:         "com.steadybit.extension_azure.scale_set.instance",
		RestrictTo: extutil.Ptr(discovery_kit_api.LEADER),
		Discover: discovery_kit_api.DescribingEndpointReferenceWithCallInterval{
			CallInterval: extutil.Ptr("30s"),
		},
	}
}

func createAzureInstanceTargets(hosts []discovery_kit_api.Target) []discovery_kit_api.Target {
	result := make([]discovery_kit_api.Target, 0, len(hosts))
	for i, host := range hosts {
		instanceId := fmt.Sprintf("i-%s", host.Id)
		instanceName := fmt.Sprintf("loadtest-instance-%s", host.Id)
		hostname := host.Id
		zone := "westeurope-1"
		if i%2 == 0 {
			zone = "westeurope-2"
		}
		target := discovery_kit_api.Target{
			Id:         instanceId,
			TargetType: "com.steadybit.extension_azure.scale_set.instance",
			Label:      instanceName,
			Attributes: map[string][]string{
				"azure-vm.vm.name":          {instanceName},
				"azure.subscription.id":     {"00000000-0000-0000-0000-000000000000"},
				"azure-vm.vm.id":            {instanceId},
				"azure-vm.hostname":         {hostname},
				"azure-scale-set-instance.hostname":         {hostname},
				"azure-vm.vm.size":          {"Standard_B1s"},
				"azure-vm.os.name":          {"UbuntuServer"},
				"azure-vm.os.version":       {"20.04.202108190"},
				"azure-vm.os.type":          {"Linux"},
				"azure-vm.power.state":      {"VM running"},
				"azure-vm.network.id":       {"00000000-0000-0000-0000-000000000000"},
				"azure.location":            {"westeurope"},
				"azure.zone":                {zone},
				"azure.resource-group.name": {config.Config.ClusterName},
			},
		}
		result = append(result, target)
	}

	return discovery_kit_commons.ApplyAttributeExcludes(result, config.Config.DiscoveryAttributesExcludesEc2)
}
