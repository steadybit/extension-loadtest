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
	"strings"
)

type doNothingAction struct {
	targetId          string
	selectionTemplate action_kit_api.TargetSelectionTemplate
}

// Make sure action implements all required interfaces
var (
	_ action_kit_sdk.Action[DoNothingActionState] = (*doNothingAction)(nil)
)

type DoNothingActionState struct {
}

type DoNothingActionConfig struct {
}

func NewDoNothingAction(targetId string, selectionTemplate action_kit_api.TargetSelectionTemplate) action_kit_sdk.Action[DoNothingActionState] {
	return &doNothingAction{
		targetId:          targetId,
		selectionTemplate: selectionTemplate,
	}
}

func (l *doNothingAction) NewEmptyState() DoNothingActionState {
	return DoNothingActionState{}
}

func (l *doNothingAction) Describe() action_kit_api.ActionDescription {
	targetTypeShort := l.targetId[strings.LastIndex(l.targetId, ".")+1:]
	return action_kit_api.ActionDescription{
		Id:          fmt.Sprintf("com.steadybit.extension_loadtest.nothing.%s", targetTypeShort),
		Label:       "Do Nothing",
		Description: "This action does nothing.",
		Version:     extbuild.GetSemverVersionStringOrUnknown(),
		TargetSelection: extutil.Ptr(action_kit_api.TargetSelection{
			TargetType: l.targetId,
			SelectionTemplates: extutil.Ptr([]action_kit_api.TargetSelectionTemplate{
				l.selectionTemplate,
			}),
		}),
		Technology:  extutil.Ptr("Debug"),
		Kind:        action_kit_api.Attack,
		TimeControl: action_kit_api.TimeControlInstantaneous,
		Parameters: []action_kit_api.ActionParameter{
			{
				Name:        "foo",
				Label:       "Example Parameter 1",
				Description: extutil.Ptr("Example parameter"),
				Type:        action_kit_api.ActionParameterTypeString,
			},
			{
				Name:        "message",
				Label:       "Example target selection",
				Description: extutil.Ptr("Dummy target selection."),
				Type:        action_kit_api.ActionParameterTypeTargetSelection,
			},
			{
				Name:        "bar",
				Label:       "Example Parameter 2",
				Description: extutil.Ptr("Example parameter"),
				Type:        action_kit_api.ActionParameterTypeString,
			},
		},
	}
}

func (l *doNothingAction) Prepare(_ context.Context, _ *DoNothingActionState, request action_kit_api.PrepareActionRequestBody) (*action_kit_api.PrepareResult, error) {
	var config DoNothingActionConfig
	if err := extconversion.Convert(request.Config, &config); err != nil {
		return nil, extension_kit.ToError("Failed to unmarshal the config.", err)
	}

	return &action_kit_api.PrepareResult{Messages: extutil.Ptr([]action_kit_api.Message{
		{
			Level:   extutil.Ptr(action_kit_api.Info),
			Message: "Prepared do nothing",
		},
	})}, nil
}

func (l *doNothingAction) Start(_ context.Context, _ *DoNothingActionState) (*action_kit_api.StartResult, error) {
	return &action_kit_api.StartResult{
		Messages: extutil.Ptr([]action_kit_api.Message{
			{
				Level:   extutil.Ptr(action_kit_api.Info),
				Message: "Started do nothing",
			},
		})}, nil
}
