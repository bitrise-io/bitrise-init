package scanners

import (
	"github.com/bitrise-core/bitrise-init/models"
	bitriseModels "github.com/bitrise-io/bitrise/models"
	"github.com/bitrise-io/go-utils/pointers"
	stepmanModels "github.com/bitrise-io/stepman/models"
)

const (
	stepActivateSSHKeyIDComposite                 = "activate-ssh-key@3.1.0"
	stepGitCloneIDComposite                       = "git-clone@3.1.1"
	stepCertificateAndProfileInstallerIDComposite = "certificate-and-profile-installer@1.4.0"
	stepDeployToBitriseIoIDComposite              = "deploy-to-bitrise-io@1.2.2"
)

// ScannerInterface ...
type ScannerInterface interface {
	Name() string
	Configure(searchDir string)

	DetectPlatform() (bool, error)

	Options() (models.OptionModel, error)
	DefaultOptions() models.OptionModel

	Configs() map[string]bitriseModels.BitriseDataModel
	DefaultConfigs() map[string]bitriseModels.BitriseDataModel
}

func customConfigName() string {
	return "custom-config"
}

// CustomConfig ...
func CustomConfig() map[string]bitriseModels.BitriseDataModel {
	bitriseDataMap := map[string]bitriseModels.BitriseDataModel{}
	steps := []bitriseModels.StepListItemModel{}

	// ActivateSSHKey
	steps = append(steps, bitriseModels.StepListItemModel{
		stepActivateSSHKeyIDComposite: stepmanModels.StepModel{
			RunIf: pointers.NewStringPtr(`{{getenv "SSH_RSA_PRIVATE_KEY" | ne ""}}`),
		},
	})

	// GitClone
	steps = append(steps, bitriseModels.StepListItemModel{
		stepGitCloneIDComposite: stepmanModels.StepModel{},
	})

	bitriseData := models.BitriseDataWithPrimaryWorkflowSteps(steps)

	configName := customConfigName()
	bitriseDataMap[configName] = bitriseData

	return bitriseDataMap
}
