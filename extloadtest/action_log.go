/*
 * Copyright 2023 steadybit GmbH. All rights reserved.
 */

package extloadtest

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/steadybit/action-kit/go/action_kit_api/v2"
	"github.com/steadybit/action-kit/go/action_kit_sdk"
	extension_kit "github.com/steadybit/extension-kit"
	"github.com/steadybit/extension-kit/extbuild"
	"github.com/steadybit/extension-kit/extconversion"
	"github.com/steadybit/extension-kit/extfile"
	"github.com/steadybit/extension-kit/extutil"
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
	AddArtifact      bool
	Step             string
	StatusCount      int
}

type LogActionConfig struct {
	Message          string
	ErrorEndpoint    string
	LatencyEndpoint  string
	LatencyDuration  int64
	TargetFilter     string
	AddArtifact      bool
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
		TargetSelection: new(action_kit_api.TargetSelection{
			TargetType:         l.targetId,
			TargetQuery:        new("steadybit.loadtest=\"true\""),
			SelectionTemplates: new(selectionTemplates),
		}),
		Technology:  new("Debug"),
		Kind:        action_kit_api.Attack,
		TimeControl: action_kit_api.TimeControlExternal,
		Parameters: []action_kit_api.ActionParameter{
			{
				Name:         "duration",
				Label:        "Duration",
				Type:         action_kit_api.ActionParameterTypeDuration,
				DefaultValue: new("10s"),
				Required:     new(true),
			},
			{
				Name:         "message",
				Label:        "Message",
				Description:  new("What should we log to the console? Use %s to insert the target name."),
				Type:         action_kit_api.ActionParameterTypeString,
				DefaultValue: new("Hello from %s"),
				Required:     new(true),
			},
			{
				Name:         "errorEndpoint",
				Label:        "Error in endpoint",
				Description:  new("Should we throw an error in the selected endpoint?"),
				Type:         action_kit_api.ActionParameterTypeString,
				Advanced:     new(true),
				DefaultValue: new("none"),
				Options: new([]action_kit_api.ParameterOption{
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
						Label: "status (invalid property update)",
						Value: "statusPropertyUpdate",
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
				Description:  new("Should we add latency in the selected endpoint?"),
				Type:         action_kit_api.ActionParameterTypeString,
				Advanced:     new(true),
				DefaultValue: new("none"),
				Options: new([]action_kit_api.ParameterOption{
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
				DefaultValue: new("5s"),
				Required:     new(false),
				Advanced:     new(true),
			},
			{
				Name:         "targetFilter",
				Label:        "Target Filter for error / latency",
				Description:  new("For which target should we throw an error / add latency? '*' throws for all targets."),
				DefaultValue: new("*"),
				Type:         action_kit_api.ActionParameterTypeString,
				Advanced:     new(true),
			},
			{
				Name:         "addArtifact",
				Label:        "Add a dummy artifact to all results",
				DefaultValue: new("false"),
				Type:         action_kit_api.ActionParameterTypeBoolean,
				Advanced:     new(true),
			},
			{
				Name:         "booleanParameter",
				Label:        "Just a dummy boolean parameter",
				Description:  new("This is not used."),
				DefaultValue: new("false"),
				Type:         action_kit_api.ActionParameterTypeBoolean,
				Advanced:     new(true),
			},
			{
				Name:        "integerParameter",
				Label:       "Just a dummy integer parameter",
				Description: new("This is not used."),
				Type:        action_kit_api.ActionParameterTypeInteger,
				Advanced:    new(true),
			},
			{
				Name:        "keyValueParameter",
				Label:       "Just a dummy key value parameter",
				Description: new("This is not used."),
				Type:        action_kit_api.ActionParameterTypeKeyValue,
				Advanced:    new(true),
			},
			{
				Name:        "stringArrayParameter",
				Label:       "Just a dummy string array parameter",
				Description: new("This is not used."),
				Type:        action_kit_api.ActionParameterTypeStringArray,
				Advanced:    new(true),
				Options: new([]action_kit_api.ParameterOption{
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
				Description: new("This is not used."),
				Type:        action_kit_api.ActionParameterTypeStringArray,
				Advanced:    new(true),
			},
			{
				Name:        "fileParameter",
				Label:       "Just a dummy file parameter",
				Description: new("This is not used."),
				Type:        action_kit_api.ActionParameterTypeFile,
				Advanced:    new(true),
			},
			{
				Name:        "regexParameter",
				Label:       "Just a dummy regex parameter",
				Description: new("This is not used."),
				Type:        action_kit_api.ActionParameterTypeRegex,
				Advanced:    new(true),
			},
			{
				Name:        "textareaParameter",
				Label:       "Just a dummy textarea parameter",
				Description: new("This is not used."),
				Type:        action_kit_api.ActionParameterTypeTextarea,
				Advanced:    new(true),
			},
			{
				Name:        "urlParameter",
				Label:       "Just a dummy url parameter",
				Description: new("This is not used."),
				Type:        action_kit_api.ActionParameterTypeUrl,
				Advanced:    new(true),
			},
			{
				Name:        "percentageParameter",
				Label:       "Just a dummy percentage parameter",
				Description: new("This is not used."),
				Type:        action_kit_api.ActionParameterTypePercentage,
				Advanced:    new(true),
			},
			{
				Name:        "headerParameter",
				Label:       "Just a dummy header parameter",
				Description: new("This is not used."),
				Type:        action_kit_api.ActionParameterTypeHeader,
				Advanced:    new(true),
			},
			{
				Name:        "bitrateParameter",
				Label:       "Just a dummy bitrate parameter",
				Description: new("This is not used."),
				Type:        action_kit_api.ActionParameterTypeBitrate,
				Advanced:    new(true),
			},
			{
				Name:        "stressngworkersParameter",
				Label:       "Just a dummy stressng workers parameter",
				Description: new("This is not used."),
				Type:        action_kit_api.ActionParameterTypeStressngWorkers,
				Advanced:    new(true),
			},
		},
		Status: new(action_kit_api.MutatingEndpointReferenceWithCallInterval{
			CallInterval: new("1s"),
		}),
		Stop: new(action_kit_api.MutatingEndpointReference{}),
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
	state.AddArtifact = config.AddArtifact
	state.TargetName = request.Target.Name
	state.Step = "prepare"
	state.StatusCount = 0

	log.Info().Str("message", state.FormattedMessage).Msg("Logging in log action **prepare**")
	log.Info().Bool("booleanParameter", config.BooleanParameter).Msg("Value of booleanParameter in log action **prepare**")

	log.Info().Msgf("Received current State of execution properties: %+v", request.Properties)

	if state.ErrorEndpoint == "prepare" && (state.TargetFilter == "*" || state.TargetFilter == state.TargetName) {
		return nil, extension_kit.ToError("Simulated error thrown in prepare endpoint", nil)
	}
	if state.LatencyEndpoint == "prepare" && (state.TargetFilter == "*" || state.TargetFilter == state.TargetName) {
		log.Info().Int64("latency", state.LatencyDuration.Milliseconds()).Msg("Adding latency in log action **prepare**")
		time.Sleep(state.LatencyDuration)
	}

	return &action_kit_api.PrepareResult{
		Messages: new([]action_kit_api.Message{
			{
				Level:   extutil.Ptr(action_kit_api.Info),
				Message: fmt.Sprintf("Called `prepare` for logging '%+v'", state),
			},
		}),
		Modifications: new(
			[]action_kit_api.ExecutionModification{
				action_kit_api.ExecutionModificationSetPropertyValue{
					Type:        action_kit_api.SetPropertyValue,
					PropertyKey: "extensionLoadtestShowcaseDelete",
					Value:       4711,
				},
				action_kit_api.ExecutionModificationAddValueToListProperty{
					Type:        action_kit_api.AddValueToListProperty,
					PropertyKey: "extensionLoadtestShowcaseList",
					Value:       "Prepared",
				},
			},
		),
	}, nil
}

func (l *logAction) Start(_ context.Context, state *LogActionState) (*action_kit_api.StartResult, error) {
	log.Info().Str("message", state.FormattedMessage).Msg("Logging in log action **start**")
	state.Step = "start"

	if state.ErrorEndpoint == "start" && (state.TargetFilter == "*" || state.TargetFilter == state.TargetName) {
		return nil, extension_kit.ToError("Simulated error thrown in start endpoint", nil)
	}
	if state.LatencyEndpoint == "start" && (state.TargetFilter == "*" || state.TargetFilter == state.TargetName) {
		log.Info().Int64("latency", state.LatencyDuration.Milliseconds()).Msg("Adding latency in log action **start**")
		time.Sleep(state.LatencyDuration)
	}

	return &action_kit_api.StartResult{
		Messages: new([]action_kit_api.Message{
			{
				Level:   extutil.Ptr(action_kit_api.Info),
				Message: fmt.Sprintf("Called `start` for logging '%+v'", state),
			},
		}),
		Modifications: new(
			[]action_kit_api.ExecutionModification{
				action_kit_api.ExecutionModificationAddValueToListProperty{
					Type:        action_kit_api.AddValueToListProperty,
					PropertyKey: "extensionLoadtestShowcaseList",
					Value:       "Start",
				},
			},
		),
	}, nil
}

func (l *logAction) Status(_ context.Context, state *LogActionState) (*action_kit_api.StatusResult, error) {
	log.Info().Str("message", state.FormattedMessage).Msg("Logging in log action **status**")

	state.Step = "status"
	state.StatusCount++
	if state.ErrorEndpoint == "status" && (state.TargetFilter == "*" || state.TargetFilter == state.TargetName) {
		return nil, extension_kit.ToError("Simulated error thrown in status endpoint", nil)
	}
	if state.LatencyEndpoint == "status" && (state.TargetFilter == "*" || state.TargetFilter == state.TargetName) {
		log.Info().Int64("latency", state.LatencyDuration.Milliseconds()).Msg("Adding latency in log action **status**")
		time.Sleep(state.LatencyDuration)
	}

	modifications := []action_kit_api.ExecutionModification{
		action_kit_api.ExecutionModificationAddValueToListProperty{
			Type:        action_kit_api.AddValueToListProperty,
			PropertyKey: "extensionLoadtestShowcaseList",
			Value:       "Status " + time.Now().Format("15:04:05.000"),
		},
	}
	if (state.ErrorEndpoint == "statusPropertyUpdate") && (state.TargetFilter == "*" || state.TargetFilter == state.TargetName) {
		modifications = append(modifications, action_kit_api.ExecutionModificationSetPropertyValue{
			Type:        action_kit_api.SetPropertyValue,
			PropertyKey: "extensionLoadtestShowcaseNotEditable",
			Value:       "This will fail",
		})
	}

	var summary *action_kit_api.Summary
	if state.StatusCount == 1 && rand.Intn(2) == 0 {
		level := action_kit_api.SummaryLevelInfo
		if rand.Intn(2) == 0 {
			level = action_kit_api.SummaryLevelWarning
		}
		summary = &action_kit_api.Summary{
			Level: level,
			Text:  "Uuuuuh, lucky you are! On the first status call, you got a summary!",
		}
	}

	return &action_kit_api.StatusResult{
		//indicate that the action is still running
		Completed: false,
		//These messages will show up in agent log
		Messages: new([]action_kit_api.Message{
			{
				Level:   extutil.Ptr(action_kit_api.Info),
				Message: fmt.Sprintf("Called `status` for logging '%+v'", state),
			},
		}),
		Modifications: new(
			modifications,
		),
		Summary: summary,
	}, nil
}

func (l *logAction) Stop(_ context.Context, state *LogActionState) (*action_kit_api.StopResult, error) {
	previousStep := state.Step
	log.Info().Str("message", state.FormattedMessage).Str("previousStep", previousStep).Msg("Logging in log action **stop**")

	if state.ErrorEndpoint == "stop" && (state.TargetFilter == "*" || state.TargetFilter == state.TargetName) {
		return nil, extension_kit.ToError("Simulated error thrown in stop endpoint", nil)
	}
	if state.LatencyEndpoint == "stop" && (state.TargetFilter == "*" || state.TargetFilter == state.TargetName) {
		log.Info().Int64("latency", state.LatencyDuration.Milliseconds()).Msg("Adding latency in log action **stop**")
		time.Sleep(state.LatencyDuration)
	}

	artifacts := make([]action_kit_api.Artifact, 0)
	if state.AddArtifact {
		content, err := extfile.File2Base64("./revision.txt")
		if err != nil {
			return nil, extension_kit.ToError("Failed to encode dummy-artifact.txt.", err)
		}
		artifacts = append(artifacts, action_kit_api.Artifact{
			Label: "dummy-artifact.txt",
			Data:  content,
		})
	}

	return &action_kit_api.StopResult{
		//These messages will show up in agent log
		Messages: new([]action_kit_api.Message{
			{
				Level:   extutil.Ptr(action_kit_api.Info),
				Message: fmt.Sprintf("Called `stop` for logging '%+v' - previous step: '%s'", state, previousStep),
			},
		}),
		Modifications: new(
			[]action_kit_api.ExecutionModification{
				action_kit_api.ExecutionModificationAddValueToListProperty{
					Type:        action_kit_api.AddValueToListProperty,
					PropertyKey: "extensionLoadtestShowcaseList",
					Value:       "Stop",
				},
				action_kit_api.ExecutionModificationSetPropertyValue{
					Type:        action_kit_api.SetPropertyValue,
					PropertyKey: "extensionLoadtestShowcaseDelete",
					Value:       nil,
				},
				action_kit_api.ExecutionModificationSetPropertyValue{
					Type:        action_kit_api.SetPropertyValue,
					PropertyKey: "extensionLoadtestShowcaseMarkdown",
					Value:       "# Whoop whoop\n## This property was not assigned before\n\n- Now it should be editable.",
				},
			},
		),
		Artifacts: new(artifacts),
	}, nil
}
