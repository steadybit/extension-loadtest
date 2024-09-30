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
	"math/rand"
	"time"
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
	Name             string  `json:"name"`
	StartFailureRate float64 `json:"startFailureRate"`
}

type RecreateActionConfig struct {
	PrepareFailureRate int `json:"prepareFailureRate"`
	StartFailureRate   int `json:"startFailureRate"`
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
		Technology:  extutil.Ptr("Debug"),
		Kind:        action_kit_api.Attack,
		TimeControl: action_kit_api.TimeControlInstantaneous,
		Parameters: []action_kit_api.ActionParameter{
			{
				Name:         "prepareFailureRate",
				Label:        "Prepare Failure Rate",
				Description:  extutil.Ptr("The rate of prepare calls to fail"),
				Type:         action_kit_api.Percentage,
				DefaultValue: extutil.Ptr("0"),
				Required:     extutil.Ptr(true),
				Order:        extutil.Ptr(1),
			},
			{
				Name:         "startFailureRate",
				Label:        "Start Failure Rate",
				Description:  extutil.Ptr("The rate of Start calls to fail"),
				Type:         action_kit_api.Percentage,
				DefaultValue: extutil.Ptr("0"),
				Required:     extutil.Ptr(true),
				Order:        extutil.Ptr(2),
			},
		},
	}
}

func (r *recreateAction) Prepare(_ context.Context, state *RecreateActionState, request action_kit_api.PrepareActionRequestBody) (*action_kit_api.PrepareResult, error) {
	var config RecreateActionConfig
	if err := extconversion.Convert(request.Config, &config); err != nil {
		return nil, extension_kit.ToError("Failed to unmarshal the config.", err)
	}

	var prepareFailureRate = float64(config.PrepareFailureRate) / 100.0
	state.StartFailureRate = float64(config.StartFailureRate) / 100.0

	sleepRandom()
	randomPanic(prepareFailureRate)
	sleepRandom()

	return &action_kit_api.PrepareResult{}, nil
}

func (r *recreateAction) Start(_ context.Context, state *RecreateActionState) (*action_kit_api.StartResult, error) {
	r.callbackFn(state.Name)

	sleepRandom()
	randomPanic(state.StartFailureRate)
	sleepRandom()

	return &action_kit_api.StartResult{
		Messages: extutil.Ptr([]action_kit_api.Message{
			{
				Level:   extutil.Ptr(action_kit_api.Info),
				Message: fmt.Sprintf("Recreated %s", state.Name),
			},
		})}, nil
}

func sleepRandom() {
	time.Sleep(time.Duration(rand.Intn(3)) * time.Second)
}
func randomPanic(rate float64) {
	if rand.Float64() < rate {
		panic("Random panic!")
	}
}
