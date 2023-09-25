/*
 * Copyright 2023 steadybit GmbH. All rights reserved.
 */

package extloadtest

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/steadybit/action-kit/go/action_kit_api/v2"
	"github.com/steadybit/action-kit/go/action_kit_sdk"
	extension_kit "github.com/steadybit/extension-kit"
	"github.com/steadybit/extension-kit/extbuild"
	"github.com/steadybit/extension-kit/extconversion"
	"github.com/steadybit/extension-kit/extutil"
)

type logAction struct {
	targetId          string
	selectionTemplate action_kit_api.TargetSelectionTemplate
}

// Make sure action implements all required interfaces
var (
	_ action_kit_sdk.Action[LogActionState]           = (*logAction)(nil)
	_ action_kit_sdk.ActionWithStatus[LogActionState] = (*logAction)(nil) // Optional, needed when the action needs a status method
	_ action_kit_sdk.ActionWithStop[LogActionState]   = (*logAction)(nil) // Optional, needed when the action needs a stop method
)

type LogActionState struct {
	FormattedMessage string
}

type LogActionConfig struct {
	Message string
}

func NewLogAction(targetId string, selectionTemplate action_kit_api.TargetSelectionTemplate) action_kit_sdk.Action[LogActionState] {
	return &logAction{
		targetId:          targetId,
		selectionTemplate: selectionTemplate,
	}
}

func (l *logAction) NewEmptyState() LogActionState {
	return LogActionState{}
}

func (l *logAction) Describe() action_kit_api.ActionDescription {
	return action_kit_api.ActionDescription{
		Id:          fmt.Sprintf("%s.log", l.targetId),
		Label:       "log message",
		Description: "logs",
		Version:     extbuild.GetSemverVersionStringOrUnknown(),
		TargetSelection: extutil.Ptr(action_kit_api.TargetSelection{
			TargetType: l.targetId,
			SelectionTemplates: extutil.Ptr([]action_kit_api.TargetSelectionTemplate{
				l.selectionTemplate,
			}),
		}),
		Category:    extutil.Ptr("internal"),
		Kind:        action_kit_api.Attack,
		TimeControl: action_kit_api.TimeControlExternal,
		Parameters: []action_kit_api.ActionParameter{
			{
				Name:         "message",
				Label:        "Message",
				Description:  extutil.Ptr("What should we log to the console? Use %s to insert the target name."),
				Type:         action_kit_api.String,
				DefaultValue: extutil.Ptr("Hello from %s"),
				Required:     extutil.Ptr(true),
				Order:        extutil.Ptr(0),
			},
			{
				Name:         "duration",
				Label:        "Duration",
				Type:         action_kit_api.Duration,
				DefaultValue: extutil.Ptr("10s"),
				Required:     extutil.Ptr(true),
				Order:        extutil.Ptr(0),
			},
		},
		Status: extutil.Ptr(action_kit_api.MutatingEndpointReferenceWithCallInterval{
			CallInterval: extutil.Ptr("1s"),
		}),
		Stop: extutil.Ptr(action_kit_api.MutatingEndpointReference{}),
	}
}

func (l *logAction) Prepare(_ context.Context, state *LogActionState, request action_kit_api.PrepareActionRequestBody) (*action_kit_api.PrepareResult, error) {
	log.Info().Str("message", state.FormattedMessage).Msg("Logging in log action **prepare**")

	var config LogActionConfig
	if err := extconversion.Convert(request.Config, &config); err != nil {
		return nil, extension_kit.ToError("Failed to unmarshal the config.", err)
	}

	state.FormattedMessage = fmt.Sprintf(config.Message, request.Target.Name)

	return &action_kit_api.PrepareResult{
		Messages: extutil.Ptr([]action_kit_api.Message{
			{
				Level:   extutil.Ptr(action_kit_api.Info),
				Message: fmt.Sprintf("Prepared logging '%s'", state.FormattedMessage),
			},
		})}, nil
}

func (l *logAction) Start(_ context.Context, state *LogActionState) (*action_kit_api.StartResult, error) {
	log.Info().Str("message", state.FormattedMessage).Msg("Logging in log action **start**")

	return &action_kit_api.StartResult{
		Messages: extutil.Ptr([]action_kit_api.Message{
			{
				Level:   extutil.Ptr(action_kit_api.Info),
				Message: fmt.Sprintf("Started logging '%s'", state.FormattedMessage),
			},
		})}, nil
}

func (l *logAction) Status(_ context.Context, state *LogActionState) (*action_kit_api.StatusResult, error) {
	log.Info().Str("message", state.FormattedMessage).Msg("Logging in log action **status**")

	return &action_kit_api.StatusResult{
		//indicate that the action is still running
		Completed: false,
		//These messages will show up in agent log
		Messages: extutil.Ptr([]action_kit_api.Message{
			{
				Level:   extutil.Ptr(action_kit_api.Info),
				Message: fmt.Sprintf("Status for logging '%s'", state.FormattedMessage),
			},
		})}, nil
}

func (l *logAction) Stop(_ context.Context, state *LogActionState) (*action_kit_api.StopResult, error) {
	log.Info().Str("message", state.FormattedMessage).Msg("Logging in log action **status**")

	return &action_kit_api.StopResult{
		//These messages will show up in agent log
		Messages: extutil.Ptr([]action_kit_api.Message{
			{
				Level:   extutil.Ptr(action_kit_api.Info),
				Message: fmt.Sprintf("Stopped logging '%s'", state.FormattedMessage),
			},
		})}, nil
}
