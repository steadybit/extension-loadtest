package extloadtest

import (
	"fmt"
	"github.com/steadybit/discovery-kit/go/discovery_kit_api"
	"github.com/steadybit/discovery-kit/go/discovery_kit_commons"
	"github.com/steadybit/extension-kit/extutil"
	"github.com/steadybit/extension-loadtest/config"
)

func getDiscoveryHost() discovery_kit_api.DiscoveryDescription {
	return discovery_kit_api.DiscoveryDescription{
		Id:         "com.steadybit.extension_host.host",
		RestrictTo: extutil.Ptr(discovery_kit_api.LEADER),
		Discover: discovery_kit_api.DescribingEndpointReferenceWithCallInterval{
			CallInterval: extutil.Ptr("1m"),
		},
	}
}

func getHostname(i int) string {
	return fmt.Sprintf("host-%d", i)
}

func createHostTargets() []discovery_kit_api.Target {
	result := make([]discovery_kit_api.Target, 0, config.Config.NodeCount)
	for i := 1; i <= config.Config.NodeCount; i++ {
		hostname := getHostname(i)
		target := discovery_kit_api.Target{
			Id:         hostname,
			TargetType: "com.steadybit.extension_host.host",
			Label:      hostname,
			Attributes: map[string][]string{
				"host.hostname":   {hostname},
				"host.domainname": {hostname},
				"host.ipv4":       {fmt.Sprintf("10.3.85.%d", i)},
				"host.ipv6": {
					"fe80::28:7eff:fe61:9b77",
					"fe80::f87a:cdff:fe08:46aa",
					"fe80::8ac:5aff:fe29:4eb3",
					"fe80::f06d:9fff:fe4e:cd4e",
					"fe80::f446:34ff:fe82:5b6c",
					"fe80::6488:aeff:fef2:7d33",
					"fe80::50aa:8dff:feff:e63d",
					"fe80::184c:9fff:fe79:e67b",
					"fe80::40d6:fff:fe40:ce1f",
					"fe80::644c:eeff:fe55:f16e",
					"fe80::14be:eff:fed6:f8e2",
					"fe80::1425:44ff:fead:ea4b",
					"fe80::784c:f9ff:fe48:a552",
					"fe80::801b:8ff:fe66:9e3f",
				},
				"host.nic": {
					"lo",
					"eth0",
					"eni5d9f55c199b",
					"eni0e36d046f38",
					"eni8e49267bbcc",
					"eni2df9d204cb4",
					"enic8f6de69406",
					"enia74d9db256a",
					"eni5889d5e5d4e",
					"enia527b9c1b53",
					"enif0723310b32",
					"enie22af453ffc",
					"eni245c7f5d214",
					"eni0efde140147",
					"enid1d3582b592",
				},
				"host.os.family":       {"debian"},
				"host.os.manufacturer": {"Debian GNU/Linux"},
				"host.os.version":      {"12 (bookworm)"},
			},
		}
		result = append(result, target)
	}
	return discovery_kit_commons.ApplyAttributeExcludes(result, config.Config.DiscoveryAttributesExcludesEc2)
}
