package jekyll

import (
	"github.com/bitrise-core/bitrise-init/utility"
	"github.com/bitrise-core/bitrise-init/steps"
	"github.com/bitrise-core/bitrise-init/models"
	envmanModels "github.com/bitrise-io/envman/models"
)

const (
	// ScannerName ...
	ScannerName   = "jekyll"
	// ConfigName ...
	ConfigName = "jekyll-config"
	// DefaultConfigName ...
	DefaultConfigName = "default-jekyll-config"

	configYmlFile = "_config.yml"
	gemfileFile   = "Gemfile"

	jekyllInitialBuildScriptName = "Install dependencies & build"
	jekyllInitialBuildCommand =
		"#!/usr/bin/env bash\n" +
		"# fail if any commands fails\n" +
		"set -e\n" +
		"# debug log\n" +
		"set -x\n" +
		"bundle install && bundle exec jekyll build\n"
)

// GenerateConfigBuilder ...
func GenerateConfigBuilder(isIncludeCache bool) models.ConfigBuilderModel {
	configBuilder := models.NewDefaultConfigBuilder()

	configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.DefaultPrepareStepList(isIncludeCache)...)
	configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.ScriptSteplistItem(jekyllInitialBuildScriptName,
		envmanModels.EnvironmentItemModel{"content": jekyllInitialBuildCommand},
	))
	configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.DefaultDeployStepList(isIncludeCache)...)

	configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.DefaultPrepareStepList(isIncludeCache)...)
	configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.ScriptSteplistItem(jekyllInitialBuildScriptName,
		envmanModels.EnvironmentItemModel{"content": jekyllInitialBuildCommand},
	))
	configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.DefaultDeployStepList(isIncludeCache)...)


	return *configBuilder
}

func filterProjectFile(fileName string, fileList []string) (string, error) {
	allowGivenFileBaseFilter := utility.BaseFilter(fileName, true)
	filePaths, err := utility.FilterPaths(fileList, allowGivenFileBaseFilter)
	if err != nil {
		return "", err
	}

	if len(filePaths) == 0 {
		return "", nil
	}

	return filePaths[0], nil
}
