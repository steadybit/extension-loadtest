// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

package extloadtest

import "github.com/steadybit/action-kit/go/action_kit_api/v2"

// The groups of the log action cover what platforms with parameter groups must render: a group with a description,
// every parameter type inside a group, and a long group without a description. "Add a dummy artifact" stays advanced
// without a group, so it is shown in "Additional Settings" next to them.
const (
	groupFaultInjection = "Fault Injection"
	groupParameterTypes = "Parameter Types"
	groupManyParameters = "Many Parameters"
)

var (
	faultInjectionGroup = action_kit_api.ParameterGroup{
		Name:        groupFaultInjection,
		Description: "Makes one lifecycle step of the action fail or slow down.",
	}
	parameterTypesGroup = action_kit_api.ParameterGroup{
		Name:        groupParameterTypes,
		Description: "One parameter of every type, to check how each renders in a group.",
	}
)
