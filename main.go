/*
 * Copyright 2023 steadybit GmbH. All rights reserved.
 */

package main

import (
	"github.com/rs/zerolog"
	"github.com/steadybit/action-kit/go/action_kit_api/v2"
	"github.com/steadybit/action-kit/go/action_kit_sdk"
	"github.com/steadybit/discovery-kit/go/discovery_kit_api"
	"github.com/steadybit/event-kit/go/event_kit_api"
	"github.com/steadybit/extension-kit/extbuild"
	"github.com/steadybit/extension-kit/exthealth"
	"github.com/steadybit/extension-kit/exthttp"
	"github.com/steadybit/extension-kit/extlogging"
	"github.com/steadybit/extension-kit/extruntime"
	"github.com/steadybit/extension-kit/extutil"
	"github.com/steadybit/extension-loadtest/config"
	"github.com/steadybit/extension-loadtest/extloadtest"
)

func main() {
	extlogging.InitZeroLog()
	extbuild.PrintBuildInformation()
	extruntime.LogRuntimeInformation(zerolog.DebugLevel)

	exthealth.SetReady(false)
	exthealth.StartProbes(8083)

	config.ParseConfiguration()
	config.ValidateConfiguration()

	exthttp.RegisterHttpHandler("/", exthttp.GetterAsHandler(getExtensionList))

	extloadtest.RegisterAllDiscoveryHandlers()

	action_kit_sdk.InstallSignalHandler()
	action_kit_sdk.RegisterCoverageEndpoints()

	action_kit_sdk.RegisterAction(extloadtest.NewLogAction("com.steadybit.extension_host.host", action_kit_api.TargetSelectionTemplate{
		Label:       "by host name",
		Description: extutil.Ptr("Find host by host name."),
		Query:       "host.hostname=\"\"",
	}))
	action_kit_sdk.RegisterAction(extloadtest.NewLogAction("com.steadybit.extension_aws.ec2-instance", action_kit_api.TargetSelectionTemplate{
		Label:       "by instance-id",
		Description: extutil.Ptr("Find ec2-instance by instance-id"),
		Query:       "aws-ec2.instance.id=\"\"",
	}))
	action_kit_sdk.RegisterAction(extloadtest.NewLogAction("com.steadybit.extension_container.container", action_kit_api.TargetSelectionTemplate{
		Label:       "by kubernetes deployment",
		Description: extutil.Ptr("Find container by kubernetes deployment."),
		Query:       "k8s.cluster-name=\"\" and k8s.namespace=\"\" and k8s.deployment=\"\"",
	}))
	action_kit_sdk.RegisterAction(extloadtest.NewLogAction("com.steadybit.extension_kubernetes.kubernetes-deployment", action_kit_api.TargetSelectionTemplate{
		Label:       "default",
		Description: extutil.Ptr("Find deployment by cluster, namespace and deployment"),
		Query:       "k8s.cluster-name=\"\" AND k8s.namespace=\"\" AND k8s.deployment=\"\"",
	}))
  action_kit_sdk.RegisterAction(extloadtest.NewDoNothingAction("com.steadybit.extension_container.container", action_kit_api.TargetSelectionTemplate{
		Label:       "by kubernetes deployment",
		Description: extutil.Ptr("Find container by kubernetes deployment."),
		Query:       "k8s.cluster-name=\"\" and k8s.namespace=\"\" and k8s.deployment=\"\"",
	}))

	exthealth.SetReady(true)

	exthttp.Listen(exthttp.ListenOpts{
		Port: 8082,
	})
}

type ExtensionListResponse struct {
	action_kit_api.ActionList       `json:",inline"`
	discovery_kit_api.DiscoveryList `json:",inline"`
	event_kit_api.EventListenerList `json:",inline"`
}

func getExtensionList() ExtensionListResponse {
	return ExtensionListResponse{
		ActionList: action_kit_sdk.GetActionList(),
		DiscoveryList: discovery_kit_api.DiscoveryList{
			Discoveries: []discovery_kit_api.DescribingEndpointReference{
				{
					Method: "GET",
					Path:   "/discovery/host",
				},
				{
					Method: "GET",
					Path:   "/discovery/container",
				},
				{
					Method: "GET",
					Path:   "/discovery/ec2-instance",
				},
				{
					Method: "GET",
					Path:   "/discovery/kubernetes-cluster",
				},
				{
					Method: "GET",
					Path:   "/discovery/kubernetes-deployment",
				},
				{
					Method: "GET",
					Path:   "/discovery/kubernetes-container",
				},
			},
		},
	}
}
