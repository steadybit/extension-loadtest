package extloadtest

import (
	"context"
	"errors"
	"github.com/google/uuid"
	extension_kit "github.com/steadybit/extension-kit"
	"github.com/steadybit/extension-kit/extutil"
	"github.com/steadybit/preflight-kit/go/preflight_kit_api"
	"strings"
	"sync"
)

type ExamplePreflight struct {
}

var runningPreflights = sync.Map{}
var statusCount = sync.Map{}

func NewExamplePreflight() *ExamplePreflight {
	return &ExamplePreflight{}
}

func (preflight *ExamplePreflight) Describe() preflight_kit_api.PreflightDescription {
	return preflight_kit_api.PreflightDescription{
		Id:                      "com.steadybit.extension_loadtest.preflight.example",
		Version:                 "v0.0.1",
		Label:                   "ExamplePreflightId",
		Description:             "This is an Example Preflight",
		TargetAttributeIncludes: []string{"host.hostname", "k8s.deployment"},
		Start:                   preflight_kit_api.MutatingEndpointReference{},
		Status: preflight_kit_api.MutatingEndpointReferenceWithCallInterval{
			CallInterval: extutil.Ptr("1s"),
		},
		Cancel: &preflight_kit_api.MutatingEndpointReference{},
	}
}

func (preflight *ExamplePreflight) Start(_ context.Context, request preflight_kit_api.StartPreflightRequestBody) (*preflight_kit_api.StartResult, error) {
	runningPreflights.Store(request.PreflightActionExecutionId, request.ExperimentExecution)
	return &preflight_kit_api.StartResult{}, nil
}

func (preflight *ExamplePreflight) Status(_ context.Context, request preflight_kit_api.StatusPreflightRequestBody) (*preflight_kit_api.StatusResult, error) {
	count := incrementStatusCounter(request.PreflightActionExecutionId)
	if count < 2 {
		return &preflight_kit_api.StatusResult{Completed: false, Error: nil}, nil
	}
	loadedExecution, ok := runningPreflights.Load(request.PreflightActionExecutionId)
	if !ok {
		return nil, extutil.Ptr(extension_kit.ToError("Could not find preflight", errors.New("preflight not found")))
	}
	var execution = loadedExecution.(preflight_kit_api.ExperimentExecutionAO)
	if execution.Name != nil && strings.Contains(strings.ToLower(*execution.Name), "technicalerror") {
		return nil, extutil.Ptr(extension_kit.ToError("This is a test error", errors.New("with some details")))
	} else if execution.Name != nil && strings.Contains(strings.ToLower(*execution.Name), "failed") {
		return &preflight_kit_api.StatusResult{
			Completed: true,
			Error:     &preflight_kit_api.PreflightKitError{Title: "Preflight says: NO!", Detail: extutil.Ptr("because no"), Status: extutil.Ptr(preflight_kit_api.Failed)},
		}, nil
	} else if execution.Name != nil && strings.Contains(strings.ToLower(*execution.Name), "error") {
		return &preflight_kit_api.StatusResult{
			Completed: true,
			Error:     &preflight_kit_api.PreflightKitError{Title: "Preflight says: Oh NO. Error!", Detail: extutil.Ptr("because no"), Status: extutil.Ptr(preflight_kit_api.Errored)},
		}, nil
	} else if execution.Name != nil && strings.Contains(strings.ToLower(*execution.Name), "success") {
		return &preflight_kit_api.StatusResult{
			Completed: true,
			Error:     nil,
		}, nil
	} else if execution.Name != nil && strings.Contains(strings.ToLower(*execution.Name), "timeout") {
		return &preflight_kit_api.StatusResult{
			Completed: false,
			Error:     nil,
		}, nil
	}
	return &preflight_kit_api.StatusResult{Completed: true}, nil
}

func incrementStatusCounter(preflightActionExecutionId uuid.UUID) int {
	increment, _ := statusCount.LoadOrStore(preflightActionExecutionId, 0)
	count := increment.(int) + 1
	statusCount.Store(preflightActionExecutionId, count)
	return count
}

func (preflight *ExamplePreflight) Cancel(_ context.Context, request preflight_kit_api.CancelPreflightRequestBody) (*preflight_kit_api.CancelResult, error) {
	runningPreflights.Delete(request.PreflightActionExecutionId)
	return &preflight_kit_api.CancelResult{}, nil
}
