/*
 * Copyright 2023 steadybit GmbH. All rights reserved.
 */

package extloadtest

import (
	"context"
	"fmt"
	"github.com/steadybit/action-kit/go/action_kit_api/v2"
	"github.com/steadybit/action-kit/go/action_kit_sdk"
	"github.com/steadybit/extension-kit/extbuild"
	"github.com/steadybit/extension-kit/extutil"
	"math/rand/v2"
	"strconv"
	"time"
)

type widgetAction struct {
}

var (
	_ action_kit_sdk.Action[WidgetActionState]           = (*widgetAction)(nil)
	_ action_kit_sdk.ActionWithStatus[WidgetActionState] = (*widgetAction)(nil)
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
		Id:          fmt.Sprintf("com.steadybit.extension_loadtest.show_multiple_widgets"),
		Label:       "Render Widgets",
		Description: "Showcase for multiple widgets in the run view",
		Version:     extbuild.GetSemverVersionStringOrUnknown(),
		Technology:  extutil.Ptr("Debug"),
		Kind:        action_kit_api.Other,
		TimeControl: action_kit_api.TimeControlExternal,
		Parameters: []action_kit_api.ActionParameter{
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
		Widgets: extutil.Ptr([]action_kit_api.Widget{
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
				Url: extutil.Ptr(action_kit_api.StateOverTimeWidgetUrlConfig{
					From: extutil.Ptr("metric.url"),
				}),
				Value: extutil.Ptr(action_kit_api.StateOverTimeWidgetValueConfig{
					Hide: extutil.Ptr(true),
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
				Url: extutil.Ptr(action_kit_api.StateOverTimeWidgetUrlConfig{
					From: extutil.Ptr("metric.url"),
				}),
				Value: extutil.Ptr(action_kit_api.StateOverTimeWidgetValueConfig{
					Hide: extutil.Ptr(true),
				}),
			},
			action_kit_api.LogWidget{
				Type:    action_kit_api.ComSteadybitWidgetLog,
				Title:   "Log Widget",
				LogType: "EXAMPLE-LOGS",
			},
		}),
	}
}

func (l *widgetAction) Prepare(_ context.Context, state *WidgetActionState, request action_kit_api.PrepareActionRequestBody) (*action_kit_api.PrepareResult, error) {
	return nil, nil
}

func (l *widgetAction) Start(_ context.Context, state *WidgetActionState) (*action_kit_api.StartResult, error) {
	return nil, nil
}

func (l *widgetAction) Status(_ context.Context, state *WidgetActionState) (*action_kit_api.StatusResult, error) {
	randomState1 := []string{"danger", "warn", "info", "success"}[rand.IntN(4)]
	randomState2 := []string{"danger", "warn", "info", "success"}[rand.IntN(4)]
	randomState3 := []string{"danger", "warn", "info", "success"}[rand.IntN(4)]

	metrics := []action_kit_api.Metric{
		{
			Name: extutil.Ptr("example-state-over-time-1"),
			Metric: map[string]string{
				"metric.type-1":  "example-row-1",
				"metric.label":   "Row 1",
				"metric.state":   randomState1,
				"metric.tooltip": fmt.Sprintf("State: %s", randomState1),
				"metric.url":     fmt.Sprintf("https://www.google.com/search?q=%s", randomState1),
			},
			Timestamp: time.Now(),
			Value:     0,
		},
		{
			Name: extutil.Ptr("example-state-over-time-1"),
			Metric: map[string]string{
				"metric.type-1":  "example-row-2",
				"metric.label":   "Row 2",
				"metric.state":   randomState2,
				"metric.tooltip": fmt.Sprintf("State: %s", randomState2),
				"metric.url":     fmt.Sprintf("https://www.google.com/search?q=%s", randomState2),
			},
			Timestamp: time.Now(),
			Value:     0,
		},
		{
			Name: extutil.Ptr("example-state-over-time-2"),
			Metric: map[string]string{
				"metric.type-2":  "example-row-1-for-second-widget",
				"metric.label":   "Row 1 - Widget 2",
				"metric.state":   randomState3,
				"metric.tooltip": fmt.Sprintf("State: %s", randomState3),
				"metric.url":     fmt.Sprintf("https://www.google.com/search?q=%s", randomState3),
			},
			Timestamp: time.Now(),
			Value:     0,
		},
	}

	randomLevel := []action_kit_api.MessageLevel{action_kit_api.Error, action_kit_api.Warn, action_kit_api.Info, action_kit_api.Debug}[rand.IntN(4)]
	messages := []action_kit_api.Message{
		{
			Message:         "This is an example log message",
			Type:            extutil.Ptr("EXAMPLE-LOGS"),
			Level:           extutil.Ptr(randomLevel),
			Timestamp:       extutil.Ptr(time.Now()),
			TimestampSource: extutil.Ptr(action_kit_api.TimestampSourceExternal),
			Fields: extutil.Ptr(action_kit_api.MessageFields{
				"tooltip-example": strconv.Itoa(rand.IntN(20)),
			}),
		},
	}

	return &action_kit_api.StatusResult{
		Completed: false,
		Metrics:   extutil.Ptr(metrics),
		Messages:  extutil.Ptr(messages),
	}, nil
}
