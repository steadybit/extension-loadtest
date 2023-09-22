package extloadtest

import (
	"fmt"
	"github.com/steadybit/discovery-kit/go/discovery_kit_api"
	"github.com/steadybit/extension-kit/extutil"
	"github.com/steadybit/extension-loadtest/config"
)

func getDiscoveryEc2Instance() discovery_kit_api.DiscoveryDescription {
	return discovery_kit_api.DiscoveryDescription{
		Id:         "com.steadybit.extension_aws.ec2-instance",
		RestrictTo: extutil.Ptr(discovery_kit_api.LEADER),
		Discover: discovery_kit_api.DescribingEndpointReferenceWithCallInterval{
			Method:       "GET",
			Path:         "/discovery/ec2-instance/targets",
			CallInterval: extutil.Ptr("30s"),
		},
	}
}

func initEc2InstanceTargets() []discovery_kit_api.Target {
	result := make([]discovery_kit_api.Target, 0, config.Config.NodeCount)
	for i := 1; i <= config.Config.NodeCount; i++ {
		instanceId := fmt.Sprintf("i-000%d", i)
		instanceName := fmt.Sprintf("loadtest-instance-%d", i)
		hostname := getHostname(i)
		target := discovery_kit_api.Target{
			Id:         instanceId,
			TargetType: "com.steadybit.extension_aws.ec2-instance",
			Label:      instanceName,
			Attributes: map[string][]string{
				"aws-ec2.arn":                                     {fmt.Sprintf("arn:aws:ec2:eu-central-1:503171660203:instance/%s", instanceId)},
				"aws-ec2.hostname.internal":                       {hostname},
				"aws-ec2.hostname.public":                         {""},
				"aws-ec2.image":                                   {"ami-02fc9c535f43bbc91"},
				"aws-ec2.instance.id":                             {instanceId},
				"aws-ec2.instance.name":                           {instanceName},
				"aws-ec2.ipv4.private":                            {fmt.Sprintf("10.40.95.%d", i)},
				"aws-ec2.label.account_name":                      {"sandbox"},
				"aws-ec2.label.application":                       {"demo"},
				"aws-ec2.label.aws:autoscaling:groupname":         {"eks-sandbox-demo-ngroup2-c2c3879b-0659-aac4-0524-b06eedbf55b7"},
				"aws-ec2.label.aws:ec2:fleet-id":                  {"fleet-2606b70f-e835-ebb6-acb2-a48a9cb6cc6b"},
				"aws-ec2.label.aws:ec2launchtemplate:id":          {"lt-063fb519ebf2a336c"},
				"aws-ec2.label.aws:ec2launchtemplate:version":     {"3"},
				"aws-ec2.label.aws:eks:cluster-name":              {config.Config.ClusterName},
				"aws-ec2.label.eks:cluster-name":                  {config.Config.ClusterName},
				"aws-ec2.label.eks:nodegroup-name":                {"sandbox-demo-ngroup2"},
				"aws-ec2.label.environment":                       {"sandbox"},
				"aws-ec2.label.k8s.io/cluster-autoscaler/enabled": {"true"},
				fmt.Sprintf("aws-ec2.label.k8s.io/cluster-autoscaler/%s", config.Config.ClusterName): {"owned"},
				fmt.Sprintf("aws-ec2.label.kubernetes.io/cluster/%s", config.Config.ClusterName):     {"owned"},
				"aws-ec2.label.type": {"eks"},
				"aws-ec2.state":      {"running"},
				"aws-ec2.vpc":        {"vpc-00000ab91434cb717"},
				"aws.account":        {"503171660203"},
				"aws.region":         {"eu-central-1"},
				"aws.zone":           {"eu-central-1b"},
			},
		}
		result = append(result, target)
	}

	return discovery_kit_api.ApplyAttributeExcludes(result, config.Config.DiscoveryAttributeExcludesEc2)
}
