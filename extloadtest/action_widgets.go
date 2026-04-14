/*
 * Copyright 2023 steadybit GmbH. All rights reserved.
 */

package extloadtest

import (
	"context"
	"fmt"
	"math/rand/v2"
	"strconv"
	"time"

	"github.com/steadybit/action-kit/go/action_kit_api/v2"
	"github.com/steadybit/action-kit/go/action_kit_sdk"
	"github.com/steadybit/extension-kit/extbuild"
	"github.com/steadybit/extension-kit/extutil"
)

type widgetAction struct {
}

var (
	_ action_kit_sdk.Action[WidgetActionState]           = (*widgetAction)(nil)
	_ action_kit_sdk.ActionWithStatus[WidgetActionState] = (*widgetAction)(nil)
)

const (
	TYPE_LOG_EXAMPLE              = "EXAMPLE-LOGS"
	TYPE_MARKDOWN_APPEND_EXAMPLE  = "EXAMPLE-MARKDOWN-APPEND"
	TYPE_MARKDOWN_REPLACE_EXAMPLE = "EXAMPLE-MARKDOWN-REPLACE"
)

type WidgetActionState struct {
}

type WidgetActionConfig struct {
}

func NewWidgetAction() action_kit_sdk.Action[WidgetActionState] {
	return &widgetAction{}
}

func (l *widgetAction) NewEmptyState() WidgetActionState {
	return WidgetActionState{}
}

func (l *widgetAction) Describe() action_kit_api.ActionDescription {
	return action_kit_api.ActionDescription{
		Id:          "com.steadybit.extension_loadtest.show_multiple_widgets",
		Label:       "Render Widgets",
		Description: "Showcase for multiple widgets in the run view",
		Version:     extbuild.GetSemverVersionStringOrUnknown(),
		TargetSelection: new(action_kit_api.TargetSelection{
			TargetType:  "com.steadybit.extension_host.host",
			TargetQuery: new("steadybit.loadtest=\"true\""),
			SelectionTemplates: new([]action_kit_api.TargetSelectionTemplate{
				{
					Label:       "host name",
					Description: new("Find by host name"),
					Query:       "host.hostname=\"\"",
				},
			}),
		}),
		Technology: new("Debug"),

		Kind:        action_kit_api.Other,
		TimeControl: action_kit_api.TimeControlExternal,
		Parameters: []action_kit_api.ActionParameter{
			{
				Name:         "duration",
				Label:        "Duration",
				Type:         action_kit_api.ActionParameterTypeDuration,
				DefaultValue: new("10s"),
				Required:     new(true),
			},
		},
		Status: new(action_kit_api.MutatingEndpointReferenceWithCallInterval{
			CallInterval: new("1s"),
		}),
		Widgets: new([]action_kit_api.Widget{
			action_kit_api.StateOverTimeWidget{
				Type:  action_kit_api.ComSteadybitWidgetStateOverTime,
				Title: "Example State Over Time Widget 1",
				Identity: action_kit_api.StateOverTimeWidgetIdentityConfig{
					From: "metric.type-1",
				},
				Label: action_kit_api.StateOverTimeWidgetLabelConfig{
					From: "metric.label",
				},
				State: action_kit_api.StateOverTimeWidgetStateConfig{
					From: "metric.state",
				},
				Tooltip: action_kit_api.StateOverTimeWidgetTooltipConfig{
					From: "metric.tooltip",
				},
				Url: new(action_kit_api.StateOverTimeWidgetUrlConfig{
					From: new("metric.url"),
				}),
				Value: new(action_kit_api.StateOverTimeWidgetValueConfig{
					Hide: new(true),
				}),
			},
			action_kit_api.StateOverTimeWidget{
				Type:  action_kit_api.ComSteadybitWidgetStateOverTime,
				Title: "Example State Over Time Widget 2",
				Identity: action_kit_api.StateOverTimeWidgetIdentityConfig{
					From: "metric.type-2",
				},
				Label: action_kit_api.StateOverTimeWidgetLabelConfig{
					From: "metric.label",
				},
				State: action_kit_api.StateOverTimeWidgetStateConfig{
					From: "metric.state",
				},
				Tooltip: action_kit_api.StateOverTimeWidgetTooltipConfig{
					From: "metric.tooltip",
				},
				Url: new(action_kit_api.StateOverTimeWidgetUrlConfig{
					From: new("metric.url"),
				}),
				Value: new(action_kit_api.StateOverTimeWidgetValueConfig{
					Hide: new(true),
				}),
			},
			action_kit_api.LogWidget{
				Type:    action_kit_api.ComSteadybitWidgetLog,
				Title:   "Log Widget",
				LogType: TYPE_LOG_EXAMPLE,
			},
			action_kit_api.MarkdownWidget{
				Type:        action_kit_api.ComSteadybitWidgetMarkdown,
				Title:       "Example Markdown Widget with appending content",
				MessageType: TYPE_MARKDOWN_APPEND_EXAMPLE,
				Append:      true,
			},
			action_kit_api.MarkdownWidget{
				Type:        action_kit_api.ComSteadybitWidgetMarkdown,
				Title:       "Example Markdown Widget with replacing content",
				MessageType: TYPE_MARKDOWN_REPLACE_EXAMPLE,
				Append:      false,
			},
		}),
	}
}

func (l *widgetAction) Prepare(_ context.Context, _ *WidgetActionState, _ action_kit_api.PrepareActionRequestBody) (*action_kit_api.PrepareResult, error) {
	return nil, nil
}

func (l *widgetAction) Start(_ context.Context, _ *WidgetActionState) (*action_kit_api.StartResult, error) {
	return &action_kit_api.StartResult{
		Messages: &[]action_kit_api.Message{
			{
				Message: "# This will be a header",
				Type:    new(TYPE_MARKDOWN_APPEND_EXAMPLE),
			},
			{
				Message: "## And a nice sub-header",
				Type:    new(TYPE_MARKDOWN_APPEND_EXAMPLE),
			},
			{
				Message: "# Hello - I just started",
				Type:    new(TYPE_MARKDOWN_REPLACE_EXAMPLE),
			},
		},
	}, nil
}

func (l *widgetAction) Status(_ context.Context, _ *WidgetActionState) (*action_kit_api.StatusResult, error) {
	randomState1 := []string{"danger", "warn", "info", "success"}[rand.IntN(4)]
	randomState2 := []string{"danger", "warn", "info", "success"}[rand.IntN(4)]
	randomState3 := []string{"danger", "warn", "info", "success"}[rand.IntN(4)]
	randomEmoji := []string{"✅", "🔴", "🙂"}[rand.IntN(3)]
	now := time.Now()

	metrics := []action_kit_api.Metric{
		{
			Name: new("example-state-over-time-1"),
			Metric: map[string]string{
				"metric.type-1":  "example-row-1",
				"metric.label":   "Row 1",
				"metric.state":   randomState1,
				"metric.tooltip": fmt.Sprintf("State: %s", randomState1),
				"metric.url":     fmt.Sprintf("https://www.google.com/search?q=%s", randomState1),
			},
			Timestamp: now,
			Value:     0,
		},
		{
			Name: new("example-state-over-time-1"),
			Metric: map[string]string{
				"metric.type-1":  "example-row-2",
				"metric.label":   "Row 2",
				"metric.state":   randomState2,
				"metric.tooltip": fmt.Sprintf("State: %s", randomState2),
				"metric.url":     fmt.Sprintf("https://www.google.com/search?q=%s", randomState2),
			},
			Timestamp: now,
			Value:     0,
		},
		{
			Name: new("example-state-over-time-2"),
			Metric: map[string]string{
				"metric.type-2":  "example-row-1-for-second-widget",
				"metric.label":   "Row 1 - Widget 2",
				"metric.state":   randomState3,
				"metric.tooltip": fmt.Sprintf("State: %s", randomState3),
				"metric.url":     fmt.Sprintf("https://www.google.com/search?q=%s", randomState3),
			},
			Timestamp: now,
			Value:     0,
		},
	}

	randomLevel := []action_kit_api.MessageLevel{action_kit_api.Error, action_kit_api.Warn, action_kit_api.Info, action_kit_api.Debug}[rand.IntN(4)]
	messages := []action_kit_api.Message{
		{
			Message:         "This is an example log message",
			Type:            new(TYPE_LOG_EXAMPLE),
			Level:           new(randomLevel),
			Timestamp:       new(now),
			TimestampSource: extutil.Ptr(action_kit_api.TimestampSourceExternal),
			Fields: new(action_kit_api.MessageFields{
				"tooltip-example": strconv.Itoa(rand.IntN(20)),
			}),
		},
		{
			Message: "- Status " + now.Format("01-02-2006 15:04:05.000000"),
			Type:    new(TYPE_MARKDOWN_APPEND_EXAMPLE),
		},
		{
			Message:   "- Status " + now.Format("01-02-2006 15:04:05.000000"),
			Type:      new(TYPE_MARKDOWN_REPLACE_EXAMPLE),
			Timestamp: new(now),
		},
		{
			Message:   "- Condition " + randomEmoji,
			Type:      new(TYPE_MARKDOWN_REPLACE_EXAMPLE),
			Timestamp: new(now),
		},
	}

	return &action_kit_api.StatusResult{
		Completed: false,
		Metrics:   new(metrics),
		Messages:  new(messages),
	}, nil
}
