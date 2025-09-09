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
	"time"
)

type logAction struct {
	targetId          string
	selectionTemplate *action_kit_api.TargetSelectionTemplate
	actionId          string
	actionLabel       string
}

// Make sure action implements all required interfaces
var (
	_ action_kit_sdk.Action[LogActionState]           = (*logAction)(nil)
	_ action_kit_sdk.ActionWithStatus[LogActionState] = (*logAction)(nil) // Optional, needed when the action needs a status method
	_ action_kit_sdk.ActionWithStop[LogActionState]   = (*logAction)(nil) // Optional, needed when the action needs a stop method
)

type LogActionState struct {
	FormattedMessage string
	ErrorEndpoint    string
	LatencyEndpoint  string
	LatencyDuration  time.Duration
	TargetFilter     string
	TargetName       string
}

type LogActionConfig struct {
	Message          string
	ErrorEndpoint    string
	LatencyEndpoint  string
	LatencyDuration  int64
	TargetFilter     string
	BooleanParameter bool
}

func NewLogAction(actionId string, targetId string, selectionTemplate action_kit_api.TargetSelectionTemplate) action_kit_sdk.Action[LogActionState] {
	return NewLogActionWithLabel(actionId, targetId, &selectionTemplate, "Log message")
}

func NewLogActionWithLabel(actionId string, targetId string, selectionTemplate *action_kit_api.TargetSelectionTemplate, actionLabel string) action_kit_sdk.Action[LogActionState] {
	return &logAction{
		actionId:          actionId,
		targetId:          targetId,
		selectionTemplate: selectionTemplate,
		actionLabel:       actionLabel,
	}
}

func (l *logAction) NewEmptyState() LogActionState {
	return LogActionState{}
}

func (l *logAction) Describe() action_kit_api.ActionDescription {
	selectionTemplates := []action_kit_api.TargetSelectionTemplate{}
	if l.selectionTemplate != nil {
		selectionTemplates = []action_kit_api.TargetSelectionTemplate{*l.selectionTemplate}
	}

	return action_kit_api.ActionDescription{
		Id:          l.actionId,
		Label:       l.actionLabel,
		Description: "Logs a message for the given duration to the agent log",
		Version:     extbuild.GetSemverVersionStringOrUnknown(),
		TargetSelection: extutil.Ptr(action_kit_api.TargetSelection{
			TargetType:         l.targetId,
			SelectionTemplates: extutil.Ptr(selectionTemplates),
		}),
		Technology:  extutil.Ptr("Debug"),
		Kind:        action_kit_api.Attack,
		TimeControl: action_kit_api.TimeControlExternal,
		Parameters: []action_kit_api.ActionParameter{
			{
				Name:         "duration",
				Label:        "Duration",
				Type:         action_kit_api.ActionParameterTypeDuration,
				DefaultValue: extutil.Ptr("10s"),
				Required:     extutil.Ptr(true),
			},
			{
				Name:         "message",
				Label:        "Message",
				Description:  extutil.Ptr("What should we log to the console? Use %s to insert the target name."),
				Type:         action_kit_api.ActionParameterTypeString,
				DefaultValue: extutil.Ptr("Hello from %s"),
				Required:     extutil.Ptr(true),
			},
			{
				Name:         "errorEndpoint",
				Label:        "Error in endpoint",
				Description:  extutil.Ptr("Should we throw an error in the selected endpoint?"),
				Type:         action_kit_api.ActionParameterTypeString,
				Advanced:     extutil.Ptr(true),
				DefaultValue: extutil.Ptr("none"),
				Options: extutil.Ptr([]action_kit_api.ParameterOption{
					action_kit_api.ExplicitParameterOption{
						Label: "no error",
						Value: "none",
					},
					action_kit_api.ExplicitParameterOption{
						Label: "prepare",
						Value: "prepare",
					},
					action_kit_api.ExplicitParameterOption{
						Label: "start",
						Value: "start",
					},
					action_kit_api.ExplicitParameterOption{
						Label: "status",
						Value: "status",
					},
					action_kit_api.ExplicitParameterOption{
						Label: "stop",
						Value: "stop",
					},
				}),
			},
			{
				Name:         "latencyEndpoint",
				Label:        "Latency in endpoint",
				Description:  extutil.Ptr("Should we add latency in the selected endpoint?"),
				Type:         action_kit_api.ActionParameterTypeString,
				Advanced:     extutil.Ptr(true),
				DefaultValue: extutil.Ptr("none"),
				Options: extutil.Ptr([]action_kit_api.ParameterOption{
					action_kit_api.ExplicitParameterOption{
						Label: "no latency",
						Value: "none",
					},
					action_kit_api.ExplicitParameterOption{
						Label: "prepare",
						Value: "prepare",
					},
					action_kit_api.ExplicitParameterOption{
						Label: "start",
						Value: "start",
					},
					action_kit_api.ExplicitParameterOption{
						Label: "status",
						Value: "status",
					},
					action_kit_api.ExplicitParameterOption{
						Label: "stop",
						Value: "stop",
					},
				}),
			},
			{
				Name:         "latencyDuration",
				Label:        "Latency",
				Type:         action_kit_api.ActionParameterTypeDuration,
				DefaultValue: extutil.Ptr("5s"),
				Required:     extutil.Ptr(false),
				Advanced:     extutil.Ptr(true),
			},
			{
				Name:         "targetFilter",
				Label:        "Target Filter for error / latency",
				Description:  extutil.Ptr("For which target should we throw an error / add latency? '*' throws for all targets."),
				DefaultValue: extutil.Ptr("*"),
				Type:         action_kit_api.ActionParameterTypeString,
				Advanced:     extutil.Ptr(true),
			},
			{
				Name:         "booleanParameter",
				Label:        "Just a dummy boolean parameter",
				Description:  extutil.Ptr("This is not used."),
				DefaultValue: extutil.Ptr("false"),
				Type:         action_kit_api.ActionParameterTypeBoolean,
				Advanced:     extutil.Ptr(true),
			},
			{
				Name:        "integerParameter",
				Label:       "Just a dummy integer parameter",
				Description: extutil.Ptr("This is not used."),
				Type:        action_kit_api.ActionParameterTypeInteger,
				Advanced:    extutil.Ptr(true),
			},
			{
				Name:        "keyValueParameter",
				Label:       "Just a dummy key value parameter",
				Description: extutil.Ptr("This is not used."),
				Type:        action_kit_api.ActionParameterTypeKeyValue,
				Advanced:    extutil.Ptr(true),
			},
			{
				Name:        "stringArrayParameter",
				Label:       "Just a dummy string array parameter",
				Description: extutil.Ptr("This is not used."),
				Type:        action_kit_api.ActionParameterTypeStringArray,
				Advanced:    extutil.Ptr(true),
				Options: extutil.Ptr([]action_kit_api.ParameterOption{
					action_kit_api.ExplicitParameterOption{
						Label: "value1",
						Value: "value1",
					},
					action_kit_api.ExplicitParameterOption{
						Label: "value2",
						Value: "value2",
					},
					action_kit_api.ExplicitParameterOption{
						Label: "value3",
						Value: "value3",
					},
				}),
			},
			{
				Name:        "stringArrayWithoutOptionsParameter",
				Label:       "Just a dummy string array without options parameter",
				Description: extutil.Ptr("This is not used."),
				Type:        action_kit_api.ActionParameterTypeStringArray,
				Advanced:    extutil.Ptr(true),
			},
			{
				Name:        "fileParameter",
				Label:       "Just a dummy file parameter",
				Description: extutil.Ptr("This is not used."),
				Type:        action_kit_api.ActionParameterTypeFile,
				Advanced:    extutil.Ptr(true),
			},
			{
				Name:        "regexParameter",
				Label:       "Just a dummy regex parameter",
				Description: extutil.Ptr("This is not used."),
				Type:        action_kit_api.ActionParameterTypeRegex,
				Advanced:    extutil.Ptr(true),
			},
			{
				Name:        "textareaParameter",
				Label:       "Just a dummy textarea parameter",
				Description: extutil.Ptr("This is not used."),
				Type:        action_kit_api.ActionParameterTypeTextarea,
				Advanced:    extutil.Ptr(true),
			},
			{
				Name:        "urlParameter",
				Label:       "Just a dummy url parameter",
				Description: extutil.Ptr("This is not used."),
				Type:        action_kit_api.ActionParameterTypeUrl,
				Advanced:    extutil.Ptr(true),
			},
			{
				Name:        "percentageParameter",
				Label:       "Just a dummy percentage parameter",
				Description: extutil.Ptr("This is not used."),
				Type:        action_kit_api.ActionParameterTypePercentage,
				Advanced:    extutil.Ptr(true),
			},
			{
				Name:        "headerParameter",
				Label:       "Just a dummy header parameter",
				Description: extutil.Ptr("This is not used."),
				Type:        action_kit_api.ActionParameterTypeHeader,
				Advanced:    extutil.Ptr(true),
			},
			{
				Name:        "bitrateParameter",
				Label:       "Just a dummy bitrate parameter",
				Description: extutil.Ptr("This is not used."),
				Type:        action_kit_api.ActionParameterTypeBitrate,
				Advanced:    extutil.Ptr(true),
			},
			{
				Name:        "stressngworkersParameter",
				Label:       "Just a dummy stressng workers parameter",
				Description: extutil.Ptr("This is not used."),
				Type:        action_kit_api.ActionParameterTypeStressngWorkers,
				Advanced:    extutil.Ptr(true),
			},
		},
		Status: extutil.Ptr(action_kit_api.MutatingEndpointReferenceWithCallInterval{
			CallInterval: extutil.Ptr("1s"),
		}),
		Stop: extutil.Ptr(action_kit_api.MutatingEndpointReference{}),
	}
}

func (l *logAction) Prepare(_ context.Context, state *LogActionState, request action_kit_api.PrepareActionRequestBody) (*action_kit_api.PrepareResult, error) {
	var config LogActionConfig
	if err := extconversion.Convert(request.Config, &config); err != nil {
		return nil, extension_kit.ToError("Failed to unmarshal the config.", err)
	}
	state.FormattedMessage = fmt.Sprintf(config.Message, request.Target.Name)
	state.ErrorEndpoint = config.ErrorEndpoint
	state.LatencyEndpoint = config.LatencyEndpoint
	state.LatencyDuration = time.Duration(config.LatencyDuration * int64(time.Millisecond))
	state.TargetFilter = config.TargetFilter
	state.TargetName = request.Target.Name

	log.Info().Str("message", state.FormattedMessage).Msg("Logging in log action **prepare**")
	log.Info().Bool("booleanParameter", config.BooleanParameter).Msg("Value of booleanParameter in log action **prepare**")

	if state.ErrorEndpoint == "prepare" && (state.TargetFilter == "*" || state.TargetFilter == state.TargetName) {
		return nil, extension_kit.ToError("Simulated error thrown in prepare endpoint", nil)
	}
	if state.LatencyEndpoint == "prepare" && (state.TargetFilter == "*" || state.TargetFilter == state.TargetName) {
		log.Info().Int64("latency", state.LatencyDuration.Milliseconds()).Msg("Adding latency in log action **prepare**")
		time.Sleep(state.LatencyDuration)
	}

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

	if state.ErrorEndpoint == "start" && (state.TargetFilter == "*" || state.TargetFilter == state.TargetName) {
		return nil, extension_kit.ToError("Simulated error thrown in start endpoint", nil)
	}
	if state.LatencyEndpoint == "start" && (state.TargetFilter == "*" || state.TargetFilter == state.TargetName) {
		log.Info().Int64("latency", state.LatencyDuration.Milliseconds()).Msg("Adding latency in log action **start**")
		time.Sleep(state.LatencyDuration)
	}

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

	if state.ErrorEndpoint == "status" && (state.TargetFilter == "*" || state.TargetFilter == state.TargetName) {
		return nil, extension_kit.ToError("Simulated error thrown in status endpoint", nil)
	}
	if state.LatencyEndpoint == "status" && (state.TargetFilter == "*" || state.TargetFilter == state.TargetName) {
		log.Info().Int64("latency", state.LatencyDuration.Milliseconds()).Msg("Adding latency in log action **status**")
		time.Sleep(state.LatencyDuration)
	}

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
	log.Info().Str("message", state.FormattedMessage).Msg("Logging in log action **stop**")

	if state.ErrorEndpoint == "stop" && (state.TargetFilter == "*" || state.TargetFilter == state.TargetName) {
		return nil, extension_kit.ToError("Simulated error thrown in stop endpoint", nil)
	}
	if state.LatencyEndpoint == "stop" && (state.TargetFilter == "*" || state.TargetFilter == state.TargetName) {
		log.Info().Int64("latency", state.LatencyDuration.Milliseconds()).Msg("Adding latency in log action **stop**")
		time.Sleep(state.LatencyDuration)
	}

	return &action_kit_api.StopResult{
		//These messages will show up in agent log
		Messages: extutil.Ptr([]action_kit_api.Message{
			{
				Level:   extutil.Ptr(action_kit_api.Info),
				Message: fmt.Sprintf("Stopped logging '%s'", state.FormattedMessage),
			},
		})}, nil
}
