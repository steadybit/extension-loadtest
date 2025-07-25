package extloadtest

import (
	"context"
	"errors"
	"github.com/google/uuid"
	extension_kit "github.com/steadybit/extension-kit"
	"github.com/steadybit/extension-kit/extbuild"
	"github.com/steadybit/extension-kit/extutil"
	"github.com/steadybit/preflight-kit/go/preflight_kit_api"
	"github.com/steadybit/preflight-kit/go/preflight_kit_sdk"
	"strings"
	"sync"
	"time"
)

type GithubActionPreflight struct {
}

var runningPreflights = sync.Map{}
var statusCount = sync.Map{}

func NewGitHubActionPreflight() *GithubActionPreflight {
	return &GithubActionPreflight{}
}

// Make sure action implements all required interfaces
var (
	_ preflight_kit_sdk.Preflight = (*GithubActionPreflight)(nil)
)

func (preflight *GithubActionPreflight) Describe() preflight_kit_api.PreflightDescription {
	return preflight_kit_api.PreflightDescription{
		Id:                      "com.steadybit.extension_loadtest.preflight.github-action",
		Version:                 extbuild.GetSemverVersionStringOrUnknown(),
		Label:                   "Github Action Preflight",
		Description:             "This is a Preflight for the nightly Github Action which fails depending on the experiment name.",
		TargetAttributeIncludes: []string{"host.hostname", "k8s.deployment"},
		Start:                   preflight_kit_api.MutatingEndpointReference{},
		Status: preflight_kit_api.MutatingEndpointReferenceWithCallInterval{
			CallInterval: extutil.Ptr("1s"),
		},
		Cancel: &preflight_kit_api.MutatingEndpointReference{},
		Icon:   extutil.Ptr("data:image/svg+xml,%3Csvg width='24' height='24' viewBox='0 0 24 24' fill='none' xmlns='http://www.w3.org/2000/svg'%3E%3Cpath d='M21.284 12.306V12.0427C21.284 9.96499 19.5099 8.274 17.3285 8.274H12.6367V6.66334C13.7014 6.39431 14.4883 5.46881 14.4883 4.37114C14.4883 3.06325 13.3716 2 11.9996 2C10.6277 2 9.511 3.06325 9.511 4.37114C9.511 5.46881 10.2979 6.39431 11.3626 6.66334V8.274H6.67072C4.48932 8.274 2.71528 9.96427 2.71528 12.0427V12.306C1.17014 12.5901 0 13.8887 0 15.4419C0 16.9952 1.17316 18.2966 2.71979 18.5793C2.82295 20.5702 4.55558 22.16 6.67072 22.16H17.3293C19.4444 22.16 21.1763 20.5702 21.2802 18.5793C22.8268 18.2966 24 16.9973 24 15.4419C24 13.8865 22.8299 12.5908 21.2847 12.306H21.284ZM2.71452 17.3259C1.87946 17.0691 1.27406 16.3222 1.27406 15.4419C1.27406 14.5616 1.87946 13.8148 2.71452 13.5586V17.3259ZM10.7851 4.37185C10.7851 3.73405 11.3295 3.21534 11.9996 3.21534C12.6698 3.21534 13.2134 3.73405 13.2134 4.37185C13.2134 5.00966 12.669 5.52836 11.9996 5.52836C11.3302 5.52836 10.7851 5.00966 10.7851 4.37185ZM20.0092 18.3913C20.0092 19.7996 18.8066 20.9454 17.3285 20.9454H6.66997C5.19185 20.9454 3.98858 19.7996 3.98858 18.3913V12.0434C3.98858 10.6351 5.1911 9.48933 6.66997 9.48933H17.3285C18.8066 9.48933 20.0092 10.6351 20.0092 12.0434V18.3913ZM21.284 17.3266V13.5594C22.119 13.8162 22.7244 14.5631 22.7244 15.4426C22.7244 16.3222 22.119 17.0698 21.284 17.3266Z' fill='%231D2632'/%3E%3Cpath d='M16.6471 12.0994C15.6102 12.0994 14.7669 12.9029 14.7669 13.8908C14.7669 14.8787 15.6102 15.6823 16.6471 15.6823C17.6839 15.6823 18.5273 14.8787 18.5273 13.8908C18.5273 12.9029 17.6839 12.0994 16.6471 12.0994ZM16.6471 14.4676C16.3135 14.4676 16.0417 14.2094 16.0417 13.8908C16.0417 13.5723 16.3135 13.314 16.6471 13.314C16.9807 13.314 17.2525 13.573 17.2525 13.8908C17.2525 14.2086 16.9807 14.4676 16.6471 14.4676Z' fill='%231D2632'/%3E%3Cpath d='M7.35067 12.0994C6.31381 12.0994 5.47046 12.9029 5.47046 13.8908C5.47046 14.8787 6.31381 15.6823 7.35067 15.6823C8.38754 15.6823 9.23088 14.8787 9.23088 13.8908C9.23088 12.9029 8.38754 12.0994 7.35067 12.0994ZM7.35067 14.4676C7.0171 14.4676 6.74527 14.2094 6.74527 13.8908C6.74527 13.5723 7.0171 13.314 7.35067 13.314C7.68425 13.314 7.95607 13.573 7.95607 13.8908C7.95607 14.2086 7.68425 14.4676 7.35067 14.4676Z' fill='%231D2632'/%3E%3Cpath d='M13.736 17.5698H10.2625C9.91009 17.5698 9.62546 17.8418 9.62546 18.1768C9.62546 18.5118 9.91084 18.7845 10.2625 18.7845H13.736C14.0877 18.7845 14.373 18.5126 14.373 18.1768C14.373 17.841 14.0877 17.5698 13.736 17.5698Z' fill='%231D2632'/%3E%3C/svg%3E%0A"),
	}
}

func (preflight *GithubActionPreflight) Start(_ context.Context, request preflight_kit_api.StartPreflightRequestBody) (*preflight_kit_api.StartResult, error) {
	runningPreflights.Store(request.PreflightActionExecutionId, request.ExperimentExecution)
	if request.ExperimentExecution.Name != nil && strings.Contains(strings.ToLower(*request.ExperimentExecution.Name), "startfailure") {
		return &preflight_kit_api.StartResult{Error: extutil.Ptr(preflight_kit_api.PreflightKitError{
			Title:  "Some start error",
			Status: extutil.Ptr(preflight_kit_api.Failed),
		})}, nil
	}
	return &preflight_kit_api.StartResult{}, nil
}

func (preflight *GithubActionPreflight) Status(_ context.Context, request preflight_kit_api.StatusPreflightRequestBody) (*preflight_kit_api.StatusResult, error) {
	count := incrementStatusCounter(request.PreflightActionExecutionId)
	loadedExecution, ok := runningPreflights.Load(request.PreflightActionExecutionId)
	if !ok {
		return nil, extutil.Ptr(extension_kit.ToError("Could not find preflight", errors.New("preflight not found")))
	}
	var execution = loadedExecution.(preflight_kit_api.ExperimentExecutionAO)
	if execution.Name != nil && strings.Contains(strings.ToLower(*execution.Name), "technicalerror") {
		return nil, extutil.Ptr(extension_kit.ToError("This is a test error", errors.New("with some details")))
	} else if execution.Name != nil && strings.Contains(strings.ToLower(*execution.Name), "failed") && !strings.Contains(strings.ToLower(*execution.Name), "inflightfailed") {
		return &preflight_kit_api.StatusResult{
			Completed: true,
			Error:     &preflight_kit_api.PreflightKitError{Title: "Preflight says: NO!", Detail: extutil.Ptr("because no"), Status: extutil.Ptr(preflight_kit_api.Failed)},
		}, nil
	} else if execution.Name != nil && strings.Contains(strings.ToLower(*execution.Name), "error") && !strings.Contains(strings.ToLower(*execution.Name), "inflighterror") {
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
	} else if execution.Name != nil && strings.Contains(strings.ToLower(*execution.Name), "slow") {
		if count < 10 {
			return &preflight_kit_api.StatusResult{Completed: false, Error: nil}, nil
		}
	} else if execution.Name != nil && strings.Contains(strings.ToLower(*execution.Name), "inflighterror") {
		if execution.Started != nil && time.Since(*execution.Started).Seconds() > 10 {
			return nil, extension_kit.ToError("Simulated error for inflight checks after 10 seconds.", nil)
		}
	} else if execution.Name != nil && strings.Contains(strings.ToLower(*execution.Name), "inflightfailed") {
		if execution.Started != nil && time.Since(*execution.Started).Seconds() > 10 {
			return &preflight_kit_api.StatusResult{
				Completed: true,
				Error:     &preflight_kit_api.PreflightKitError{Title: "Inflight says: NO!", Detail: extutil.Ptr("because no"), Status: extutil.Ptr(preflight_kit_api.Failed)},
			}, nil
		}
	}
	return &preflight_kit_api.StatusResult{Completed: true}, nil
}

func incrementStatusCounter(preflightActionExecutionId uuid.UUID) int {
	increment, _ := statusCount.LoadOrStore(preflightActionExecutionId, 0)
	count := increment.(int) + 1
	statusCount.Store(preflightActionExecutionId, count)
	return count
}

func (preflight *GithubActionPreflight) Cancel(_ context.Context, request preflight_kit_api.CancelPreflightRequestBody) (*preflight_kit_api.CancelResult, error) {
	runningPreflights.Delete(request.PreflightActionExecutionId)
	statusCount.Delete(request.PreflightActionExecutionId)
	return &preflight_kit_api.CancelResult{}, nil
}
