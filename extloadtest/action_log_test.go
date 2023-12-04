package extloadtest

import (
	"context"
	"github.com/steadybit/action-kit/go/action_kit_api/v2"
	"github.com/steadybit/extension-kit/extutil"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestPrepare(t *testing.T) {
	// Given
	request := extutil.JsonMangle(action_kit_api.PrepareActionRequestBody{
		Config: map[string]interface{}{
			"duration":        1000 * 60,
			"message":         "lorem ipsum %s",
			"errorEndpoint":   "stop",
			"latencyEndpoint": "start",
			"latencyDuration": 1000 * 60,
			"targetFilter":    "*",
			"targetName":      "loadtest",
		},
		Target: &action_kit_api.Target{
			Name: "example-target",
		},
		ExecutionContext: extutil.Ptr(action_kit_api.ExecutionContext{}),
	})
	action := NewLogAction("com.example.target", action_kit_api.TargetSelectionTemplate{})
	state := action.NewEmptyState()

	// When
	_, err := action.Prepare(context.TODO(), &state, request)

	// Then
	require.Nil(t, err)
	require.Equal(t, "lorem ipsum example-target", state.FormattedMessage)
	require.Equal(t, "example-target", state.TargetName)
	require.Equal(t, "stop", state.ErrorEndpoint)
	require.Equal(t, "start", state.LatencyEndpoint)
	require.Equal(t, 60*time.Second, state.LatencyDuration)
	require.Equal(t, "*", state.TargetFilter)
}
