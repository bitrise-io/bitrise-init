package toolscanner

import (
	"log"

	"github.com/bitrise-core/bitrise-init/models"
)

// ProjectTypeEnvKey is the name of the enviroment variable used to substitute the project type for
// automation tool scanner's config
const (
	ProjectTypeUserTitle = "Project type"
	ProjectTypeEnvKey    = "PROJECT_TYPE"
)

// AddProjectTypeToToolScanner is used to add a project type for automation tool scanners's option map
func AddProjectTypeToToolScanner(toolScannerOptionTree models.OptionNode, detectedProjectTypes []string) models.OptionNode {
	log.Printf("toolScannerOptionTree: %s, detectedProjectTypes: %s", toolScannerOptionTree, detectedProjectTypes)

	optionsTreeWithProjectTypeRoot := models.NewOption(ProjectTypeUserTitle, ProjectTypeEnvKey)
	for _, projectType := range detectedProjectTypes {
		optionsTreeWithProjectTypeRoot.AddOption(projectType, &toolScannerOptionTree)
	}

	log.Printf("New options root: %s", optionsTreeWithProjectTypeRoot)
	return *optionsTreeWithProjectTypeRoot
}
