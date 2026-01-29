/*
 * Copyright 2026 steadybit GmbH. All rights reserved.
 */

package extloadtest

import (
	"context"

	"github.com/steadybit/action-kit/go/action_kit_api/v2"
	"github.com/steadybit/action-kit/go/action_kit_sdk"
	"github.com/steadybit/extension-kit/extbuild"
	"github.com/steadybit/extension-kit/extutil"
)

type multiOptionParameterAction struct {
}

var (
	_ action_kit_sdk.Action[MultiOptionParameterActionState] = (*multiOptionParameterAction)(nil)
)

type MultiOptionParameterActionState struct {
}

func NewMultiOptionParameterAction() action_kit_sdk.Action[MultiOptionParameterActionState] {
	return &multiOptionParameterAction{}
}

func (l *multiOptionParameterAction) NewEmptyState() MultiOptionParameterActionState {
	return MultiOptionParameterActionState{}
}

func (l *multiOptionParameterAction) Describe() action_kit_api.ActionDescription {
	return action_kit_api.ActionDescription{
		Id:          "com.steadybit.extension_loadtest.multi_parameters",
		Label:       "Do nothing but provide an option parameter",
		Description: "This action does nothing but provides an option parameter based on multiple target attributes (k8s.container.id, k8s.deployment).",
		Version:     extbuild.GetSemverVersionStringOrUnknown(),
		Technology:  extutil.Ptr("Debug"),
		Kind:        action_kit_api.Other,
		TimeControl: action_kit_api.TimeControlInstantaneous,
		TargetSelection: &action_kit_api.TargetSelection{
			TargetType: "com.steadybit.extension_kubernetes.kubernetes-pod",
		},
		Parameters: []action_kit_api.ActionParameter{
			{
				Name:         "Parameter with option",
				Label:        "Option",
				Description:  extutil.Ptr("Select on option (static, k8s.container_id, k8s.deployment)"),
				Type:         action_kit_api.ActionParameterTypeString,
				Advanced:     extutil.Ptr(true),
				DefaultValue: extutil.Ptr("none"),
				Options: extutil.Ptr([]action_kit_api.ParameterOption{
					action_kit_api.ExplicitParameterOption{
						Label: "static",
						Value: "static",
					},
					action_kit_api.ParameterOptionsFromTargetAttribute{
						Attribute: "k8s.container.id",
					},
					action_kit_api.ParameterOptionsFromTargetAttribute{
						Attribute: "k8s.deployment",
					},
				}),
			},
		},
	}
}

func (l *multiOptionParameterAction) Prepare(_ context.Context, _ *MultiOptionParameterActionState, request action_kit_api.PrepareActionRequestBody) (*action_kit_api.PrepareResult, error) {
	return &action_kit_api.PrepareResult{Messages: extutil.Ptr([]action_kit_api.Message{
		{
			Level:   extutil.Ptr(action_kit_api.Info),
			Message: "Prepared do nothing",
		},
	})}, nil
}

func (l *multiOptionParameterAction) Start(_ context.Context, _ *MultiOptionParameterActionState) (*action_kit_api.StartResult, error) {
	return &action_kit_api.StartResult{
		Messages: extutil.Ptr([]action_kit_api.Message{
			{
				Level:   extutil.Ptr(action_kit_api.Info),
				Message: "Started do nothing",
			},
		})}, nil
}
