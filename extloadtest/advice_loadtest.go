package extloadtest

import (
	"embed"

	"github.com/rs/zerolog/log"
	"github.com/steadybit/advice-kit/go/advice_kit_api"
	"github.com/steadybit/extension-kit/extbuild"
	"github.com/steadybit/extension-kit/extutil"
)

const adviceInstruction = "Instruction \n Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum."
const adviceMotivation = "Motivation \n Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum."
const adviceSummary = "Summary \n Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum."
const adviceImplementedSummary = "Implemented \n Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum."

//go:embed *
var adviceContent embed.FS

func GetAdviceDescriptionKubernetesDeploymentLoadtest() advice_kit_api.AdviceDefinition {
	return advice_kit_api.AdviceDefinition{
		Id:                        "com.steadybit.extension-loadtest.advice.loadtest",
		Label:                     "Load Test Advice - Kubernetes Deployment",
		Version:                   extbuild.GetSemverVersionStringOrUnknown(),
		Icon:                      "data:image/svg+xml,%3Csvg%20width%3D%2224%22%20height%3D%2224%22%20viewBox%3D%220%200%2024%2024%22%20fill%3D%22none%22%20xmlns%3D%22http%3A%2F%2Fwww.w3.org%2F2000%2Fsvg%22%3E%0A%3Cpath%20d%3D%22M10.4478%202.65625C11.2739%202.24209%2012.2447%202.23174%2013.0794%202.62821L19.2871%205.57666C20.3333%206.07356%2021%207.12832%2021%208.28652V15.7134C21%2016.8717%2020.3333%2017.9264%2019.2871%2018.4233L13.0794%2021.3718C12.2447%2021.7682%2011.2739%2021.7579%2010.4478%2021.3437L4.65545%2018.4397L5.55182%2016.6518L11.3441%2019.5558C11.6195%2019.6939%2011.9431%2019.6973%2012.2214%2019.5652L18.429%2016.6167C18.7778%2016.4511%2019%2016.0995%2019%2015.7134V8.28652C19%207.90045%2018.7778%207.54887%2018.429%207.38323L12.2214%204.43479C11.9431%204.30263%2011.6195%204.30608%2011.3441%204.44413L5.55182%207.34814C5.21357%207.51773%205%207.8637%205%208.24208V15.7579C5%2016.1363%205.21357%2016.4822%205.55182%2016.6518L4.65545%2018.4397C3.6407%2017.931%203%2016.893%203%2015.7579V8.24208C3%207.10694%203.6407%206.06901%204.65545%205.56026L10.4478%202.65625Z%22%20fill%3D%22%231D2632%22%2F%3E%0A%3Cpath%20d%3D%22M11.1377%207.16465C11.5966%206.95033%2012.1359%206.94497%2012.5997%207.15014L16.0484%208.67595C16.6296%208.9331%2017%209.47893%2017%2010.0783V13.9217C17%2014.5211%2016.6296%2015.0669%2016.0484%2015.324L12.5997%2016.8499C12.1359%2017.055%2011.5966%2017.0497%2011.1377%2016.8353L7.9197%2015.3325C7.35594%2015.0693%207%2014.5321%207%2013.9447V10.0553C7%209.46787%207.35594%208.93074%207.9197%208.66747L11.1377%207.16465Z%22%20fill%3D%22%231D2632%22%2F%3E%0A%3C%2Fsvg%3E%0A",
		Tags:                      &[]string{"kubernetes", "deployment", "strategy"},
		AssessmentQueryApplicable: "target.type=\"com.steadybit.extension_kubernetes.kubernetes-deployment\"",
		Status: advice_kit_api.AdviceDefinitionStatus{
			ActionNeeded: advice_kit_api.AdviceDefinitionStatusActionNeeded{
				AssessmentQuery: "k8s.specification.has-rolling-update-strategy=\"false\"",
				Description: advice_kit_api.AdviceDefinitionStatusActionNeededDescription{
					Instruction: adviceInstruction,
					Motivation:  adviceMotivation,
					Summary:     adviceSummary,
				},
			},
			ValidationNeeded: advice_kit_api.AdviceDefinitionStatusValidationNeeded{
				Description: advice_kit_api.AdviceDefinitionStatusValidationNeededDescription{
					Summary: "Validation summary - with a value from the target: ${target.attr('k8s.deployment')}, ${target.attr('k8s.pod.name')}, ${target.attr('k8s.namespace')}",
				},
				Validation: &[]advice_kit_api.Validation{
					{
						Id:               "com.steadybit.extension-loadtest.advice.loadtest.validation.1",
						Name:             "Experiment Validation",
						Description:      extutil.Ptr("This validation creates an *experiment*. Old and boring."),
						ShortDescription: "With experiment",
						Type:             advice_kit_api.EXPERIMENT,
						Experiment:       extutil.Ptr(advice_kit_api.Experiment(LoadFile("advice_loadtest_validation_experiment.json"))),
					}, {
						Id:                 "com.steadybit.extension-loadtest.advice.loadtest.validation.2",
						Name:               "Experiment Template Validation",
						Description:        extutil.Ptr("This validation creates an *experiment* based on a **template**."),
						ShortDescription:   "With template",
						Type:               advice_kit_api.EXPERIMENT,
						ExperimentTemplate: extutil.Ptr(advice_kit_api.ExperimentTemplate(LoadFile("advice_loadtest_validation_experiment_template.json"))),
					}, {
						Id:                 "com.steadybit.extension-loadtest.advice.loadtest.validation.3",
						Name:               "Experiment Template Validation - No Variables",
						Description:        extutil.Ptr("This validation creates an *experiment* based on a **template** without any template variables."),
						ShortDescription:   "Template without template variables",
						Type:               advice_kit_api.EXPERIMENT,
						ExperimentTemplate: extutil.Ptr(advice_kit_api.ExperimentTemplate(LoadFile("advice_loadtest_validation_experiment_template_empty.json"))),
					}, {
						Id:                 "com.steadybit.extension-loadtest.advice.loadtest.validation.4",
						Name:               "Experiment Template Validation - Invalid",
						Description:        extutil.Ptr("This advice contains an invalid experiment template and invalid advice template syntax in the short description field."),
						ShortDescription:   "Broken advice template syntax ${target.attr('agent.hostname'),invalid}",
						Type:               advice_kit_api.EXPERIMENT,
						ExperimentTemplate: extutil.Ptr(advice_kit_api.ExperimentTemplate("{\"foo\":\"bar\"}")),
					}, {
						Id:                 "com.steadybit.extension-loadtest.advice.loadtest.validation.5",
						Name:               "Experiment Template Validation - Unknown Action",
						Description:        extutil.Ptr("This advice contains an unknown action."),
						ShortDescription:   "Unknown action",
						Type:               advice_kit_api.EXPERIMENT,
						ExperimentTemplate: extutil.Ptr(advice_kit_api.ExperimentTemplate(LoadFile("advice_loadtest_validation_experiment_template_unknown_action.json"))),
					}, {
						Id:                 "com.steadybit.extension-loadtest.advice.loadtest.validation.6",
						Name:               "Experiment Template Validation - Restricted Action",
						Description:        extutil.Ptr("This advice contains a restricted action."),
						ShortDescription:   "Restricted action",
						Type:               advice_kit_api.EXPERIMENT,
						ExperimentTemplate: extutil.Ptr(advice_kit_api.ExperimentTemplate(LoadFile("advice_loadtest_validation_experiment_template_restricted_action.json"))),
					}, {
						Id:                 "com.steadybit.extension-loadtest.advice.loadtest.validation.7",
						Name:               "Experiment Template Validation - Template with Advice Placeholder",
						Description:        extutil.Ptr("This advice contains a template with a placeholder."),
						ShortDescription:   "Template with advice placeholder",
						Type:               advice_kit_api.EXPERIMENT,
						ExperimentTemplate: extutil.Ptr(advice_kit_api.ExperimentTemplate(LoadFile("advice_loadtest_validation_experiment_template_placeholder.json"))),
					},
				},
			},
			Implemented: advice_kit_api.AdviceDefinitionStatusImplemented{
				Description: advice_kit_api.AdviceDefinitionStatusImplementedDescription{
					Summary: adviceImplementedSummary,
				},
			},
		},
	}
}

func LoadFile(filename string) string {
	content, err := adviceContent.ReadFile(filename)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to read file: %s", filename)
		return ""
	}
	return string(content)
}
