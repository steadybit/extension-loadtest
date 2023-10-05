/*
 * Copyright 2023 steadybit GmbH. All rights reserved.
 */

package extloadtest

import (
	"context"
	"fmt"
	"github.com/steadybit/action-kit/go/action_kit_api/v2"
	"github.com/steadybit/action-kit/go/action_kit_sdk"
	extension_kit "github.com/steadybit/extension-kit"
	"github.com/steadybit/extension-kit/extbuild"
	"github.com/steadybit/extension-kit/extconversion"
	"github.com/steadybit/extension-kit/extutil"
)

type recreateAction struct {
	targetId          string
	selectionTemplate action_kit_api.TargetSelectionTemplate
	callbackFn        func(name string)
}

// Make sure action implements all required interfaces
var (
	_ action_kit_sdk.Action[RecreateActionState] = (*recreateAction)(nil)
)

type RecreateActionState struct {
	Name string `json:"name"`
}

type RecreateActionConfig struct {
}

func NewRecreateAction(targetId string, selectionTemplate action_kit_api.TargetSelectionTemplate, callbackFn func(name string)) action_kit_sdk.Action[RecreateActionState] {
	return &recreateAction{
		targetId:          targetId,
		selectionTemplate: selectionTemplate,
		callbackFn:        callbackFn,
	}
}

func (r *recreateAction) NewEmptyState() RecreateActionState {
	return RecreateActionState{}
}

func (r *recreateAction) Describe() action_kit_api.ActionDescription {
	return action_kit_api.ActionDescription{
		Id:          fmt.Sprintf("%s.recreate", r.targetId),
		Label:       "Recreate targets",
		Description: "Simulate targets removal and creation by altering the id",
		Version:     extbuild.GetSemverVersionStringOrUnknown(),
		TargetSelection: extutil.Ptr(action_kit_api.TargetSelection{
			TargetType: r.targetId,
			SelectionTemplates: extutil.Ptr([]action_kit_api.TargetSelectionTemplate{
				r.selectionTemplate,
			}),
		}),
		Category:    extutil.Ptr("internal"),
		Kind:        action_kit_api.Attack,
		TimeControl: action_kit_api.TimeControlInstantaneous,
	}
}

func (r *recreateAction) Prepare(_ context.Context, state *RecreateActionState, request action_kit_api.PrepareActionRequestBody) (*action_kit_api.PrepareResult, error) {
	var config RecreateActionConfig
	if err := extconversion.Convert(request.Config, &config); err != nil {
		return nil, extension_kit.ToError("Failed to unmarshal the config.", err)
	}

	state.Name = request.Target.Name

	return &action_kit_api.PrepareResult{}, nil
}

func (r *recreateAction) Start(_ context.Context, state *RecreateActionState) (*action_kit_api.StartResult, error) {
	r.callbackFn(state.Name)
	return &action_kit_api.StartResult{
		Messages: extutil.Ptr([]action_kit_api.Message{
			{
				Level:   extutil.Ptr(action_kit_api.Info),
				Message: fmt.Sprintf("Recreated %s", state.Name),
			},
		})}, nil
}
