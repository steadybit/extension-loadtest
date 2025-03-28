/*
 * Copyright 2023 steadybit GmbH. All rights reserved.
 */

package main

import (
	_ "github.com/KimMachineGun/automemlimit" // By default, it sets `GOMEMLIMIT` to 90% of cgroup's memory limit.
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/steadybit/action-kit/go/action_kit_api/v2"
	"github.com/steadybit/action-kit/go/action_kit_sdk"
	"github.com/steadybit/advice-kit/go/advice_kit_api"
	"github.com/steadybit/discovery-kit/go/discovery_kit_api"
	"github.com/steadybit/discovery-kit/go/discovery_kit_sdk"
	"github.com/steadybit/event-kit/go/event_kit_api"
	"github.com/steadybit/extension-kit/extbuild"
	"github.com/steadybit/extension-kit/exthealth"
	"github.com/steadybit/extension-kit/exthttp"
	"github.com/steadybit/extension-kit/extlogging"
	"github.com/steadybit/extension-kit/extruntime"
	"github.com/steadybit/extension-kit/extsignals"
	"github.com/steadybit/extension-kit/extutil"
	"github.com/steadybit/extension-loadtest/config"
	"github.com/steadybit/extension-loadtest/extloadtest"
	"github.com/steadybit/preflight-kit/go/preflight_kit_api"
	"github.com/steadybit/preflight-kit/go/preflight_kit_sdk"
	_ "go.uber.org/automaxprocs" // Importing automaxprocs automatically adjusts GOMAXPROCS.
)

func main() {
	extlogging.InitZeroLog()
	extbuild.PrintBuildInformation()
	extruntime.LogRuntimeInformation(zerolog.DebugLevel)

	exthealth.SetReady(false)
	exthealth.StartProbes(8083)

	config.ParseConfiguration()
	config.ValidateConfiguration()

	exthttp.RegisterHttpHandler("/", exthttp.GetterAsHandler(func() ExtensionListResponse {
		return getExtensionList(config.Config.EventListenerEnabled)
	}))

	exthttp.RegisterHttpHandler("/advice/loadtest-advice", exthttp.GetterAsHandler(
		extloadtest.GetAdviceDescriptionKubernetesDeploymentLoadtest))
	exthttp.RegisterHttpHandler("/advice/loadtest-dependencies", exthttp.GetterAsHandler(
		extloadtest.GetAdviceDescriptionKubernetesDeploymentDependencies))
	exthttp.RegisterHttpHandler("/advice/checkout-dependency", exthttp.GetterAsHandler(
		extloadtest.GetAdviceDescriptionCheckoutDependency))

	extloadtest.RegisterEventListenerHandlers()

	targetData := extloadtest.NewTargetData()
	targetData.ScheduleUpdates()
	targetData.RegisterDiscoveries()
	targetData.RegisterRecreateActions()
	targetData.RegisterConfigUpdateHandlers()

	discovery_kit_sdk.Register(extloadtest.NewEnrichmentRuleProvider())

	extsignals.ActivateSignalHandlers()
	action_kit_sdk.RegisterCoverageEndpoints()

	action_kit_sdk.RegisterAction(extloadtest.NewTargetlessAction())
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
	action_kit_sdk.RegisterAction(extloadtest.NewWidgetAction())

	preflight_kit_sdk.RegisterPreflight(extloadtest.NewGitHubActionPreflight())
	preflight_kit_sdk.RegisterPreflight(extloadtest.NewSimplePreflight())

	if config.IsPodZero() {
		log.Info().Msg("I am pod zero")
	}
	exthealth.SetReady(true)

	exthttp.Listen(exthttp.ListenOpts{
		Port: 8082,
	})
}

type ExtensionListResponse struct {
	action_kit_api.ActionList       `json:",inline"`
	discovery_kit_api.DiscoveryList `json:",inline"`
	event_kit_api.EventListenerList `json:",inline"`
	advice_kit_api.AdviceList       `json:",inline"`
	preflight_kit_api.PreflightList `json:",inline"`
}

func getExtensionList(eventListenerEnabled bool) ExtensionListResponse {

	eventListener := []event_kit_api.EventListener{}
	if eventListenerEnabled {
		eventListener = append(eventListener, event_kit_api.EventListener{
			Method: "POST",
			Path:   "/events/log",
			ListenTo: []string{
				"experiment.execution.created",
				"experiment.execution.completed",
				"experiment.execution.failed",
				"experiment.execution.canceled",
				"experiment.execution.errored",
				"experiment.execution.step-started",
				"experiment.execution.step-completed",
				"experiment.execution.step-canceled",
				"experiment.execution.step-errored",
				"experiment.execution.step-failed",
				"experiment.execution.target-started",
				"experiment.execution.target-completed",
				"experiment.execution.target-canceled",
				"experiment.execution.target-errored",
				"experiment.execution.target-failed",
			},
		})
	}

	return ExtensionListResponse{
		ActionList:    action_kit_sdk.GetActionList(),
		PreflightList: preflight_kit_sdk.GetPreflightList(),
		DiscoveryList: discovery_kit_sdk.GetDiscoveryList(),
		EventListenerList: event_kit_api.EventListenerList{
			EventListeners: eventListener,
		},
		AdviceList: advice_kit_api.AdviceList{
			Advice: []advice_kit_api.DescribingEndpointReference{{
				Method: "GET",
				Path:   "/advice/loadtest-advice",
			}, {
				Method: "GET",
				Path:   "/advice/loadtest-dependencies",
			}, {
				Method: "GET",
				Path:   "/advice/checkout-dependency",
			}},
		},
	}
}
