package toolscanner

import (
	"github.com/bitrise-core/bitrise-init/models"
)

// ProjectTypeEnvKey is the name of the enviroment variable used to substitute the project type for
// automation tool scanner's config
const (
	ProjectTypeUserTitle = "Project type"
	ProjectTypeEnvKey    = "PROJECT_TYPE"
)

// AddProjectTypeToToolScanner is used to add a project type for automation tool scanners's option map
func AddProjectTypeToToolScanner(toolScannerOptionModel models.OptionNode, detectedProjectTypes []string) models.OptionNode {
	projectTypeOption := models.NewOption(ProjectTypeUserTitle, ProjectTypeEnvKey)
	for _, projectType := range detectedProjectTypes {
		projectTypeOption.AddOption(projectType, &toolScannerOptionModel)
	}
	return *projectTypeOption
}
