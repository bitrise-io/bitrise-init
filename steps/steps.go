package steps

import (
	bitriseModels "github.com/bitrise-io/bitrise/models"
	envmanModels "github.com/bitrise-io/envman/models"
	"github.com/bitrise-io/go-utils/pointers"
	stepmanModels "github.com/bitrise-io/stepman/models"
)

// PrepareListParams describes the default prepare Step options.
type PrepareListParams struct {
	ShouldIncludeCache       bool
	ShouldIncludeActivateSSH bool
}

func stepIDComposite(ID, version string) string {
	if version != "" {
		return ID + "@" + version
	}
	return ID
}

func StepListItem(stepIDComposite, title, runIf string, inputs ...envmanModels.EnvironmentItemModel) bitriseModels.StepListItemModel {
	step := stepmanModels.StepModel{}
	if title != "" {
		step.Title = pointers.NewStringPtr(title)
	}
	if runIf != "" {
		step.RunIf = pointers.NewStringPtr(runIf)
	}
	if len(inputs) > 0 {
		step.Inputs = inputs
	}

	return bitriseModels.StepListItemModel{
		stepIDComposite: step,
	}
}

// DefaultPrepareStepList ...
func DefaultPrepareStepList(isIncludeCache bool) []bitriseModels.StepListItemModel {
	runIfCondition := `{{getenv "SSH_RSA_PRIVATE_KEY" | ne ""}}`
	stepList := []bitriseModels.StepListItemModel{
		ActivateSSHKeyStepListItem(runIfCondition),
		GitCloneStepListItem(),
	}

	if isIncludeCache {
		stepList = append(stepList, CachePullStepListItem())
	}

	return append(stepList, ScriptSteplistItem(ScriptDefaultTitle))
}

// DefaultPrepareStepListV2 ...
func DefaultPrepareStepListV2(params PrepareListParams) []bitriseModels.StepListItemModel {
	stepList := []bitriseModels.StepListItemModel{}

	if params.ShouldIncludeActivateSSH {
		stepList = append(stepList, ActivateSSHKeyStepListItem(""))
	}

	stepList = append(stepList, GitCloneStepListItem())

	if params.ShouldIncludeCache {
		stepList = append(stepList, CachePullStepListItem())
	}

	return stepList
}

// DefaultDeployStepList ...
func DefaultDeployStepList(isIncludeCache bool) []bitriseModels.StepListItemModel {
	stepList := []bitriseModels.StepListItemModel{
		DeployToBitriseIoStepListItem(),
	}

	if isIncludeCache {
		stepList = append(stepList, CachePushStepListItem())
	}

	return stepList
}

// DefaultDeployStepListV2 ...
func DefaultDeployStepListV2(shouldIncludeCache bool) []bitriseModels.StepListItemModel {
	stepList := []bitriseModels.StepListItemModel{}

	if shouldIncludeCache {
		stepList = append(stepList, CachePushStepListItem())
	}

	stepList = append(stepList, DeployToBitriseIoStepListItem())

	return stepList
}

// ActivateSSHKeyStepListItem ...
func ActivateSSHKeyStepListItem(runIfCondition string) bitriseModels.StepListItemModel {
	stepIDComposite := stepIDComposite(ActivateSSHKeyID, ActivateSSHKeyVersion)

	return StepListItem(stepIDComposite, "", runIfCondition)
}

// AndroidLintStepListItem ...
func AndroidLintStepListItem(inputs ...envmanModels.EnvironmentItemModel) bitriseModels.StepListItemModel {
	stepIDComposite := stepIDComposite(AndroidLintID, AndroidLintVersion)
	return StepListItem(stepIDComposite, "", "", inputs...)
}

// AndroidUnitTestStepListItem ...
func AndroidUnitTestStepListItem(inputs ...envmanModels.EnvironmentItemModel) bitriseModels.StepListItemModel {
	stepIDComposite := stepIDComposite(AndroidUnitTestID, AndroidUnitTestVersion)
	return StepListItem(stepIDComposite, "", "", inputs...)
}

// AndroidBuildStepListItem ...
func AndroidBuildStepListItem(inputs ...envmanModels.EnvironmentItemModel) bitriseModels.StepListItemModel {
	stepIDComposite := stepIDComposite(AndroidBuildID, AndroidBuildVersion)
	return StepListItem(stepIDComposite, "", "", inputs...)
}

// GitCloneStepListItem ...
func GitCloneStepListItem() bitriseModels.StepListItemModel {
	stepIDComposite := stepIDComposite(GitCloneID, GitCloneVersion)
	return StepListItem(stepIDComposite, "", "")
}

// CachePullStepListItem ...
func CachePullStepListItem() bitriseModels.StepListItemModel {
	stepIDComposite := stepIDComposite(CachePullID, CachePullVersion)
	return StepListItem(stepIDComposite, "", "")
}

// CachePushStepListItem ...
func CachePushStepListItem() bitriseModels.StepListItemModel {
	stepIDComposite := stepIDComposite(CachePushID, CachePushVersion)
	return StepListItem(stepIDComposite, "", "")
}

// CertificateAndProfileInstallerStepListItem ...
func CertificateAndProfileInstallerStepListItem() bitriseModels.StepListItemModel {
	stepIDComposite := stepIDComposite(CertificateAndProfileInstallerID, CertificateAndProfileInstallerVersion)
	return StepListItem(stepIDComposite, "", "")
}

// ChangeAndroidVersionCodeAndVersionNameStepListItem ...
func ChangeAndroidVersionCodeAndVersionNameStepListItem(inputs ...envmanModels.EnvironmentItemModel) bitriseModels.StepListItemModel {
	stepIDComposite := stepIDComposite(ChangeAndroidVersionCodeAndVersionNameID, ChangeAndroidVersionCodeAndVersionNameVersion)
	return StepListItem(stepIDComposite, "", "", inputs...)
}

// DeployToBitriseIoStepListItem ...
func DeployToBitriseIoStepListItem() bitriseModels.StepListItemModel {
	stepIDComposite := stepIDComposite(DeployToBitriseIoID, DeployToBitriseIoVersion)
	return StepListItem(stepIDComposite, "", "")
}

// ScriptSteplistItem ...
func ScriptSteplistItem(title string, inputs ...envmanModels.EnvironmentItemModel) bitriseModels.StepListItemModel {
	stepIDComposite := stepIDComposite(ScriptID, ScriptVersion)
	return StepListItem(stepIDComposite, title, "", inputs...)
}

// SignAPKStepListItem ...
func SignAPKStepListItem() bitriseModels.StepListItemModel {
	stepIDComposite := stepIDComposite(SignAPKID, SignAPKVersion)
	return StepListItem(stepIDComposite, "", `{{getenv "BITRISEIO_ANDROID_KEYSTORE_URL" | ne ""}}`)
}

// InstallMissingAndroidToolsStepListItem ....
func InstallMissingAndroidToolsStepListItem(inputs ...envmanModels.EnvironmentItemModel) bitriseModels.StepListItemModel {
	stepIDComposite := stepIDComposite(InstallMissingAndroidToolsID, InstallMissingAndroidToolsVersion)
	return StepListItem(stepIDComposite, "", "", inputs...)
}

// FastlaneStepListItem ...
func FastlaneStepListItem(inputs ...envmanModels.EnvironmentItemModel) bitriseModels.StepListItemModel {
	stepIDComposite := stepIDComposite(FastlaneID, FastlaneVersion)
	return StepListItem(stepIDComposite, "", "", inputs...)
}

// CocoapodsInstallStepListItem ...
func CocoapodsInstallStepListItem() bitriseModels.StepListItemModel {
	stepIDComposite := stepIDComposite(CocoapodsInstallID, CocoapodsInstallVersion)
	return StepListItem(stepIDComposite, "", "")
}

// CarthageStepListItem ...
func CarthageStepListItem(inputs ...envmanModels.EnvironmentItemModel) bitriseModels.StepListItemModel {
	stepIDComposite := stepIDComposite(CarthageID, CarthageVersion)
	return StepListItem(stepIDComposite, "", "", inputs...)
}

// RecreateUserSchemesStepListItem ...
func RecreateUserSchemesStepListItem(inputs ...envmanModels.EnvironmentItemModel) bitriseModels.StepListItemModel {
	stepIDComposite := stepIDComposite(RecreateUserSchemesID, RecreateUserSchemesVersion)
	return StepListItem(stepIDComposite, "", "", inputs...)
}

// XcodeArchiveStepListItem ...
func XcodeArchiveStepListItem(inputs ...envmanModels.EnvironmentItemModel) bitriseModels.StepListItemModel {
	stepIDComposite := stepIDComposite(XcodeArchiveID, XcodeArchiveVersion)
	return StepListItem(stepIDComposite, "", "", inputs...)
}

// XcodeBuildForTestStepListItem ...
func XcodeBuildForTestStepListItem(inputs ...envmanModels.EnvironmentItemModel) bitriseModels.StepListItemModel {
	stepIDComposite := stepIDComposite(XcodeBuildForTestID, XcodeBuildForTestVersion)
	return StepListItem(stepIDComposite, "", "", inputs...)
}

// XcodeTestStepListItem ...
func XcodeTestStepListItem(inputs ...envmanModels.EnvironmentItemModel) bitriseModels.StepListItemModel {
	stepIDComposite := stepIDComposite(XcodeTestID, XcodeTestVersion)
	return StepListItem(stepIDComposite, "", "", inputs...)
}

// XcodeArchiveMacStepListItem ...
func XcodeArchiveMacStepListItem(inputs ...envmanModels.EnvironmentItemModel) bitriseModels.StepListItemModel {
	stepIDComposite := stepIDComposite(XcodeArchiveMacID, XcodeArchiveMacVersion)
	return StepListItem(stepIDComposite, "", "", inputs...)
}

// ExportXCArchiveStepListItem ...
func ExportXCArchiveStepListItem(inputs ...envmanModels.EnvironmentItemModel) bitriseModels.StepListItemModel {
	stepIDComposite := stepIDComposite(ExportXCArchiveID, ExportXCArchiveVersion)
	return StepListItem(stepIDComposite, "", "", inputs...)
}

// XcodeTestMacStepListItem ...
func XcodeTestMacStepListItem(inputs ...envmanModels.EnvironmentItemModel) bitriseModels.StepListItemModel {
	stepIDComposite := stepIDComposite(XcodeTestMacID, XcodeTestMacVersion)
	return StepListItem(stepIDComposite, "", "", inputs...)
}

// CordovaArchiveStepListItem ...
func CordovaArchiveStepListItem(inputs ...envmanModels.EnvironmentItemModel) bitriseModels.StepListItemModel {
	stepIDComposite := stepIDComposite(CordovaArchiveID, CordovaArchiveVersion)
	return StepListItem(stepIDComposite, "", "", inputs...)
}

// IonicArchiveStepListItem ...
func IonicArchiveStepListItem(inputs ...envmanModels.EnvironmentItemModel) bitriseModels.StepListItemModel {
	stepIDComposite := stepIDComposite(IonicArchiveID, IonicArchiveVersion)
	return StepListItem(stepIDComposite, "", "", inputs...)
}

// GenerateCordovaBuildConfigStepListItem ...
func GenerateCordovaBuildConfigStepListItem(inputs ...envmanModels.EnvironmentItemModel) bitriseModels.StepListItemModel {
	stepIDComposite := stepIDComposite(GenerateCordovaBuildConfigID, GenerateCordovaBuildConfigVersion)
	return StepListItem(stepIDComposite, "", "", inputs...)
}

// JasmineTestRunnerStepListItem ...
func JasmineTestRunnerStepListItem(inputs ...envmanModels.EnvironmentItemModel) bitriseModels.StepListItemModel {
	stepIDComposite := stepIDComposite(JasmineTestRunnerID, JasmineTestRunnerVersion)
	return StepListItem(stepIDComposite, "", "", inputs...)
}

// KarmaJasmineTestRunnerStepListItem ...
func KarmaJasmineTestRunnerStepListItem(inputs ...envmanModels.EnvironmentItemModel) bitriseModels.StepListItemModel {
	stepIDComposite := stepIDComposite(KarmaJasmineTestRunnerID, KarmaJasmineTestRunnerVersion)
	return StepListItem(stepIDComposite, "", "", inputs...)
}

// NpmStepListItem ...
func NpmStepListItem(command, workdir string) bitriseModels.StepListItemModel {
	var inputs []envmanModels.EnvironmentItemModel
	if workdir != "" {
		inputs = append(inputs, envmanModels.EnvironmentItemModel{"workdir": workdir})
	}
	if command != "" {
		inputs = append(inputs, envmanModels.EnvironmentItemModel{"command": command})
	}

	stepIDComposite := stepIDComposite(NpmID, NpmVersion)
	return StepListItem(stepIDComposite, "", "", inputs...)
}

// RunEASBuildStepListItem ...
func RunEASBuildStepListItem(workdir, platform string) bitriseModels.StepListItemModel {
	var inputs []envmanModels.EnvironmentItemModel
	if platform != "" {
		inputs = append(inputs, envmanModels.EnvironmentItemModel{"platform": platform})
	}
	if workdir != "" {
		inputs = append(inputs, envmanModels.EnvironmentItemModel{"work_dir": workdir})
	}
	stepIDComposite := stepIDComposite(RunEASBuildID, RunEASBuildVersion)
	return StepListItem(stepIDComposite, "", "", inputs...)
}

// YarnStepListItem ...
func YarnStepListItem(command, workdir string) bitriseModels.StepListItemModel {
	var inputs []envmanModels.EnvironmentItemModel
	if workdir != "" {
		inputs = append(inputs, envmanModels.EnvironmentItemModel{"workdir": workdir})
	}
	if command != "" {
		inputs = append(inputs, envmanModels.EnvironmentItemModel{"command": command})
	}

	stepIDComposite := stepIDComposite(YarnID, YarnVersion)
	return StepListItem(stepIDComposite, "", "", inputs...)
}

// FlutterInstallStepListItem ...
func FlutterInstallStepListItem(inputs ...envmanModels.EnvironmentItemModel) bitriseModels.StepListItemModel {
	stepIDComposite := stepIDComposite(FlutterInstallID, FlutterInstallVersion)
	return StepListItem(stepIDComposite, "", "", inputs...)
}

// FlutterTestStepListItem ...
func FlutterTestStepListItem(inputs ...envmanModels.EnvironmentItemModel) bitriseModels.StepListItemModel {
	stepIDComposite := stepIDComposite(FlutterTestID, FlutterTestVersion)
	return StepListItem(stepIDComposite, "", "", inputs...)
}

// FlutterAnalyzeStepListItem ...
func FlutterAnalyzeStepListItem(inputs ...envmanModels.EnvironmentItemModel) bitriseModels.StepListItemModel {
	stepIDComposite := stepIDComposite(FlutterAnalyzeID, FlutterAnalyzeVersion)
	return StepListItem(stepIDComposite, "", "", inputs...)
}

// FlutterBuildStepListItem ...
func FlutterBuildStepListItem(inputs ...envmanModels.EnvironmentItemModel) bitriseModels.StepListItemModel {
	stepIDComposite := stepIDComposite(FlutterBuildID, FlutterBuildVersion)
	return StepListItem(stepIDComposite, "", "", inputs...)
}
