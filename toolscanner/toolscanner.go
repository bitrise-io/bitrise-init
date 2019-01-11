package toolscanner

import (
	"log"

	"github.com/bitrise-core/bitrise-init/models"
)

// ProjectTypeTemplateKey is the name of the enviroment variable used to substitute the project type for
// automation tool scanner's config
const (
	ProjectTypeUserTitle   = "Project type"
	ProjectTypeTemplateKey = "PROJECT_TYPE"
)

// AddProjectTypeToToolScanner is used to add a project type for automation tool scanners's option map
func AddProjectTypeToToolScanner(scannerOptionTree models.OptionNode, detectedProjectTypes []string) models.OptionNode {
	log.Printf("AddProjectTypeToToolScanner old toolScannerOptionTree: %s, detectedProjectTypes: %s", scannerOptionTree, detectedProjectTypes)

	optionsTreeWithProjectTypeRoot := models.NewOption(ProjectTypeUserTitle, ProjectTypeTemplateKey)
	for _, projectType := range detectedProjectTypes {
		optionsTreeWithProjectTypeRoot.AddOption(projectType, &scannerOptionTree)
	}

	log.Printf("AddProjectTypeToToolScanner new options root: %s", optionsTreeWithProjectTypeRoot)
	return *optionsTreeWithProjectTypeRoot
}
