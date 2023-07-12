package ios

import (
	"github.com/bitrise-io/bitrise-init/models"
	"github.com/bitrise-io/bitrise-init/steps"
	envmanModels "github.com/bitrise-io/envman/models"
)

const (
	TestRepetitionModeKey                 = "test_repetition_mode"
	TestRepetitionModeRetryOnFailureValue = "retry_on_failure"
	BuildForTestDestinationKey            = "destination"
	BuildForTestDestinationValue          = "platform=iOS Simulator,name=iPhone 8 Plus,OS=latest"
	AutomaticCodeSigningKey               = "automatic_code_signing"
	AutomaticCodeSigningValue             = "api-key"
	CacheLevelKey                         = "cache_level"
	CacheLevelNone                        = "none"

	testWorkflowID                        = "run_tests"
	testWorkflowSummary                   = "Run your Xcode tests and get the test report."
	testWorkflowDescription               = `The workflow executes the tests. The *retry_on_failure* test repetition mode is enabled.`
	testAndBuildWorkflowCommonDescription = `Next steps:
- Check out [Getting started with iOS apps](https://devcenter.bitrise.io/en/getting-started/getting-started-with-ios-apps.html).
`

	buildWorkflowID          = "build"
	buildWorkflowSummary     = "Build your Xcode project."
	buildWorkflowDescription = `The workflow only builds the project because the project scanner could not find any tests.`

	deployWorkflowID                  = "archive_and_export_app"
	deployWorkflowWithTestsSummary    = "Run your Xcode tests and create an IPA file to install your app on a device or share it with your team."
	deployWorkflowWithoutTestsSummary = "Create an IPA file to install your app on a device or share it with your team."
	deployWorkflowDescription         = `The workflow tests, builds and deploys the app using *Deploy to bitrise.io* step.

For testing the *retry_on_failure* test repetition mode is enabled.

Next steps:
- Set up [Connecting to an Apple service with API key](https://devcenter.bitrise.io/en/accounts/connecting-to-services/connecting-to-an-apple-service-with-api-key.html##).
- Or further customise code signing following our [iOS code signing](https://devcenter.bitrise.io/en/code-signing/ios-code-signing.html) guide.
`
)

type workflowSetupParams struct {
	projectType          XcodeProjectType
	configBuilder        *models.ConfigBuilderModel
	repoAccess           models.RepoAccess
	missingSharedSchemes bool
	hasTests             bool
	hasAppClip           bool
	hasPodfile           bool
	hasSPMDependencies   bool
	carthageCommand      string
	exportMethod         string
}

func createTestWorkflow(params workflowSetupParams) {
	var id string
	var summary string
	var description string

	if params.hasTests {
		id = testWorkflowID
		summary = testWorkflowSummary
		description = testWorkflowDescription
	} else {
		id = buildWorkflowID
		summary = buildWorkflowSummary
		description = buildWorkflowDescription
	}

	addSharedSetupSteps(models.WorkflowID(id), params, false, true)

	if params.hasTests {
		addTestStep(models.WorkflowID(id), params.configBuilder, params.projectType)
	} else {
		addBuildStep(models.WorkflowID(id), params.configBuilder, params.projectType)
	}

	addSharedTeardownSteps(models.WorkflowID(id), params, true)
	addSummary(params.projectType, models.WorkflowID(id), params.configBuilder, summary)
	addDescription(params.projectType, models.WorkflowID(id), params.configBuilder, description+"\n\n"+testAndBuildWorkflowCommonDescription)
}

func createDeployWorkflow(params workflowSetupParams) {
	id := deployWorkflowID
	includeCertificateAndProfileInstallStep := params.projectType == XcodeProjectTypeMacOS
	addSharedSetupSteps(models.WorkflowID(id), params, includeCertificateAndProfileInstallStep, false)

	var summary string
	if params.hasTests {
		summary = deployWorkflowWithTestsSummary
		addTestStep(models.WorkflowID(id), params.configBuilder, params.projectType)
	} else {
		summary = deployWorkflowWithoutTestsSummary
	}

	addArchiveStep(models.WorkflowID(id), params.configBuilder, params.projectType, params.hasAppClip, params.exportMethod)
	addSharedTeardownSteps(models.WorkflowID(id), params, false) // No cache in deploy workflows
	addSummary(params.projectType, models.WorkflowID(id), params.configBuilder, summary)
	addDescription(params.projectType, models.WorkflowID(id), params.configBuilder, deployWorkflowDescription)
}

// Add steps

func addTestStep(workflow models.WorkflowID, configBuilder *models.ConfigBuilderModel, projectType XcodeProjectType) {
	switch projectType {
	case XcodeProjectTypeIOS:
		configBuilder.AppendStepListItemsTo(workflow, steps.XcodeTestStepListItem(xcodeTestStepInputModels()...))
	case XcodeProjectTypeMacOS:
		configBuilder.AppendStepListItemsTo(workflow, steps.XcodeTestMacStepListItem(baseXcodeStepInputModels()...))
	}
}

func addBuildStep(workflow models.WorkflowID, configBuilder *models.ConfigBuilderModel, projectType XcodeProjectType) {
	if projectType != XcodeProjectTypeIOS {
		return
	}

	configBuilder.AppendStepListItemsTo(workflow, steps.XcodeBuildForTestStepListItem(xcodeBuildForTestStepInputModels()...))
}

func addArchiveStep(workflow models.WorkflowID, configBuilder *models.ConfigBuilderModel, projectType XcodeProjectType, hasAppClip bool, exportMethod string) {
	inputModels := xcodeArchiveStepInputModels(projectType)

	switch projectType {
	case XcodeProjectTypeIOS:
		configBuilder.AppendStepListItemsTo(workflow, steps.XcodeArchiveStepListItem(inputModels...))

		if shouldAppendExportAppClipStep(hasAppClip, exportMethod) {
			appendExportAppClipStep(configBuilder, workflow)
		}
	case XcodeProjectTypeMacOS:
		configBuilder.AppendStepListItemsTo(workflow, steps.XcodeArchiveMacStepListItem(inputModels...))
	}
}

func addSharedSetupSteps(workflow models.WorkflowID, params workflowSetupParams, includeCertificateAndProfileInstallStep, includeCache bool) {
	params.configBuilder.AppendStepListItemsTo(workflow, steps.DefaultPrepareStepList(steps.PrepareListParams{
		RepoAccess: params.repoAccess,
	})...)

	if includeCache {
		if params.hasPodfile {
			params.configBuilder.AppendStepListItemsTo(workflow, steps.RestoreCocoapodsCache())
		}
		if params.carthageCommand != "" {
			params.configBuilder.AppendStepListItemsTo(workflow, steps.RestoreCarthageCache())
		}
		if params.hasSPMDependencies {
			params.configBuilder.AppendStepListItemsTo(workflow, steps.RestoreSPMCache())
		}
	}

	if includeCertificateAndProfileInstallStep {
		params.configBuilder.AppendStepListItemsTo(workflow, steps.CertificateAndProfileInstallerStepListItem())
	}

	if params.missingSharedSchemes {
		params.configBuilder.AppendStepListItemsTo(workflow, steps.RecreateUserSchemesStepListItem(
			envmanModels.EnvironmentItemModel{ProjectPathInputKey: "$" + ProjectPathInputEnvKey},
		))
	}

	if params.hasPodfile {
		params.configBuilder.AppendStepListItemsTo(workflow, steps.CocoapodsInstallStepListItem())
	}

	if params.carthageCommand != "" {
		params.configBuilder.AppendStepListItemsTo(workflow, steps.CarthageStepListItem(
			envmanModels.EnvironmentItemModel{CarthageCommandInputKey: params.carthageCommand},
		))
	}
}

func addSharedTeardownSteps(workflow models.WorkflowID, params workflowSetupParams, includeCache bool) {
	if includeCache {
		if params.hasPodfile {
			params.configBuilder.AppendStepListItemsTo(workflow, steps.SaveCocoapodsCache())
		}
		if params.carthageCommand != "" {
			params.configBuilder.AppendStepListItemsTo(workflow, steps.SaveCarthageCache())
		}
		if params.hasSPMDependencies {
			params.configBuilder.AppendStepListItemsTo(workflow, steps.SaveSPMCache())
		}
	}

	params.configBuilder.AppendStepListItemsTo(workflow, steps.DefaultDeployStepList()...)
}

func addDescription(projectType XcodeProjectType, workflow models.WorkflowID, configBuilder *models.ConfigBuilderModel, description string) {
	if projectType != XcodeProjectTypeIOS {
		return
	}

	configBuilder.SetWorkflowDescriptionTo(workflow, description)
}

func addSummary(projectType XcodeProjectType, workflow models.WorkflowID, configBuilder *models.ConfigBuilderModel, summary string) {
	if projectType != XcodeProjectTypeIOS {
		return
	}

	configBuilder.SetWorkflowSummaryTo(workflow, summary)
}

// Helpers

func baseXcodeStepInputModels() []envmanModels.EnvironmentItemModel {
	return []envmanModels.EnvironmentItemModel{
		{ProjectPathInputKey: "$" + ProjectPathInputEnvKey},
		{SchemeInputKey: "$" + SchemeInputEnvKey},
	}
}

func xcodeTestStepInputModels() []envmanModels.EnvironmentItemModel {
	inputModels := []envmanModels.EnvironmentItemModel{
		{TestRepetitionModeKey: TestRepetitionModeRetryOnFailureValue},
		{CacheLevelKey: CacheLevelNone},
	}

	return append(baseXcodeStepInputModels(), inputModels...)
}

func xcodeBuildForTestStepInputModels() []envmanModels.EnvironmentItemModel {
	inputModels := []envmanModels.EnvironmentItemModel{
		{BuildForTestDestinationKey: BuildForTestDestinationValue},
		{CacheLevelKey: CacheLevelNone},
	}

	return append(baseXcodeStepInputModels(), inputModels...)
}

func xcodeArchiveStepInputModels(projectType XcodeProjectType) []envmanModels.EnvironmentItemModel {
	var inputModels []envmanModels.EnvironmentItemModel

	if projectType == XcodeProjectTypeIOS {
		inputModels = append(inputModels, []envmanModels.EnvironmentItemModel{
			{DistributionMethodInputKey: "$" + DistributionMethodEnvKey},
			{AutomaticCodeSigningKey: AutomaticCodeSigningValue},
			{CacheLevelKey: CacheLevelNone},
		}...)
	} else {
		inputModels = []envmanModels.EnvironmentItemModel{
			{ExportMethodInputKey: "$" + ExportMethodEnvKey},
		}
	}

	return append(baseXcodeStepInputModels(), inputModels...)
}
