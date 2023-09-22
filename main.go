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
