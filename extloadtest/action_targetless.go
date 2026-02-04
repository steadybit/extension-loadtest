/*
 * Copyright 2023 steadybit GmbH. All rights reserved.
 */

package extloadtest

import (
	"context"

	"github.com/steadybit/action-kit/go/action_kit_api/v2"
	"github.com/steadybit/action-kit/go/action_kit_sdk"
	"github.com/steadybit/extension-kit/extbuild"
	"github.com/steadybit/extension-kit/extutil"
)

type targetlessAction struct {
}

// Make sure action implements all required interfaces
var (
	_ action_kit_sdk.Action[TargetlessActionState] = (*targetlessAction)(nil)
)

type TargetlessActionState struct {
}

func NewTargetlessAction() action_kit_sdk.Action[TargetlessActionState] {
	return &targetlessAction{}
}

func (l *targetlessAction) NewEmptyState() TargetlessActionState {
	return TargetlessActionState{}
}

func (l *targetlessAction) Describe() action_kit_api.ActionDescription {
	return action_kit_api.ActionDescription{
		Id:          "com.steadybit.extension_loadtest.targetless",
		Label:       "Do Nothing Without a Target",
		Description: "This action does nothing.",
		Version:     extbuild.GetSemverVersionStringOrUnknown(),
		Technology:  extutil.Ptr("Debug"),

		Kind:        action_kit_api.Other,
		TimeControl: action_kit_api.TimeControlInstantaneous,
	}
}

func (l *targetlessAction) Prepare(_ context.Context, _ *TargetlessActionState, request action_kit_api.PrepareActionRequestBody) (*action_kit_api.PrepareResult, error) {
	return &action_kit_api.PrepareResult{Messages: extutil.Ptr([]action_kit_api.Message{
		{
			Level:   extutil.Ptr(action_kit_api.Info),
			Message: "Prepared do nothing",
		},
	})}, nil
}

func (l *targetlessAction) Start(_ context.Context, _ *TargetlessActionState) (*action_kit_api.StartResult, error) {
	return &action_kit_api.StartResult{
		Messages: extutil.Ptr([]action_kit_api.Message{
			{
				Level:   extutil.Ptr(action_kit_api.Info),
				Message: "Started do nothing",
			},
		})}, nil
}
