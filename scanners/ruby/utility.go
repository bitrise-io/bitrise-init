package ruby

import (
	"fmt"
	"path/filepath"

	"gopkg.in/yaml.v2"

	"github.com/bitrise-io/bitrise-init/models"
	"github.com/bitrise-io/bitrise-init/steps"
	"github.com/bitrise-io/bitrise-init/utility"
	"github.com/bitrise-io/go-utils/log"
)

const (
	rubyInstallScriptStepTitle   = "Install Ruby"
	rubyInstallScriptStepContent = `#!/usr/bin/env bash
set -euxo pipefail

pushd "${RUBY_PROJECT_DIR:-.}" > /dev/null

# Bitrise stacks come with asdf pre-installed to help auto-switch between various software versions
# asdf looks for the Ruby version in these files: .tool-versions, .ruby-version
# See: https://github.com/asdf-vm/asdf-ruby
asdf install ruby

popd > /dev/null
`

	bundlerInstallScriptStepTitle   = "Install dependencies"
	bundlerInstallScriptStepContent = `#!/usr/bin/env bash
set -euxo pipefail

pushd "${RUBY_PROJECT_DIR:-.}" > /dev/null

bundle install

popd > /dev/null
`
)

func checkBundler(searchDir string) bool {
	log.TPrintf("Checking for Bundler")
	hasGemfileLock := utility.FileExists(filepath.Join(searchDir, "Gemfile.lock"))

	if !hasGemfileLock {
		log.TPrintf("- Gemfile.lock - not found")
		return false
	}

	log.TPrintf("- Gemfile.lock - found")
	log.TPrintf("Bundler: detected")
	return true
}

func checkRakefile(searchDir string) bool {
	log.TPrintf("Checking for Rakefile")
	hasRakefile := utility.FileExists(filepath.Join(searchDir, "Rakefile"))

	if !hasRakefile {
		log.TPrintf("- Rakefile - not found")
		return false
	}

	log.TPrintf("- Rakefile - found")
	return true
}

func checkRubyVersion(searchDir string) bool {
	log.TPrintf("Checking for Ruby version file")

	versionFiles := []string{".ruby-version", ".tool-versions"}
	for _, versionFile := range versionFiles {
		if utility.FileExists(filepath.Join(searchDir, versionFile)) {
			log.TPrintf("- %s - found", versionFile)
			return true
		}
	}

	log.TPrintf("- Ruby version file - not found")
	return false
}

func detectTestFramework(searchDir string) string {
	log.TPrintf("Checking test framework")

	for _, fw := range testFrameworks {
		for _, detectionFile := range fw.detectionFiles {
			if utility.FileExists(filepath.Join(searchDir, detectionFile)) {
				log.TPrintf("- %s - found (%s)", fw.name, detectionFile)
				return fw.name
			}
		}
	}

	log.TPrintf("- test framework - not detected")
	return ""
}

// Options & Configs
type configDescriptor struct {
	workdir        string
	hasBundler     bool
	hasRakefile    bool
	testFramework  string
	hasRubyVersion bool
	isDefault      bool
}

func createConfigDescriptor(project project, isDefault bool) configDescriptor {
	descriptor := configDescriptor{
		workdir:        "$" + projectDirInputEnvKey,
		hasBundler:     project.hasBundler,
		hasRakefile:    project.hasRakefile,
		testFramework:  project.testFramework,
		hasRubyVersion: project.hasRubyVersion,
		isDefault:      isDefault,
	}

	// Gemfile placed in the search dir, no need to change-dir
	if project.projectRelDir == "." {
		descriptor.workdir = ""
	}

	return descriptor
}

func createDefaultConfigDescriptor() configDescriptor {
	return createConfigDescriptor(project{
		projectRelDir:  "$" + projectDirInputEnvKey,
		hasBundler:     true,
		hasRakefile:    true,
		testFramework:  "rspec",
		hasRubyVersion: true,
	}, true)
}

func configName(params configDescriptor) string {
	name := "ruby"

	if params.isDefault {
		return "default-" + name + "-config"
	}

	if params.workdir == "" {
		name = name + "-root"
	}

	if params.hasBundler {
		name = name + "-bundler"
	}

	if params.testFramework != "" {
		name = name + "-" + params.testFramework
	}

	return name + "-config"
}

func generateOptions(projects []project) (models.OptionNode, models.Warnings, models.Icons, error) {
	if len(projects) == 0 {
		return models.OptionNode{}, nil, nil, fmt.Errorf("no Gemfile files found")
	}

	projectRootOption := models.NewOption(projectDirInputTitle, projectDirInputSummary, projectDirInputEnvKey, models.TypeSelector)
	for _, project := range projects {
		descriptor := createConfigDescriptor(project, false)
		configOption := models.NewConfigOption(configName(descriptor), nil)
		projectRootOption.AddConfig(project.projectRelDir, configOption)
	}

	return *projectRootOption, nil, nil, nil
}

func generateConfigs(projects []project, sshKeyActivation models.SSHKeyActivation) (models.BitriseConfigMap, error) {
	configs := models.BitriseConfigMap{}

	if len(projects) == 0 {
		return models.BitriseConfigMap{}, fmt.Errorf("no Gemfile files found")
	}

	for _, project := range projects {
		descriptor := createConfigDescriptor(project, false)
		config, err := generateConfigBasedOn(descriptor, sshKeyActivation)
		if err != nil {
			return nil, err
		}
		configs[configName(descriptor)] = config
	}

	return configs, nil
}

func generateConfigBasedOn(descriptor configDescriptor, sshKey models.SSHKeyActivation) (string, error) {
	configBuilder := models.NewDefaultConfigBuilder()
	prepareSteps := steps.DefaultPrepareStepList(steps.PrepareListParams{SSHKeyActivation: sshKey})
	configBuilder.AppendStepListItemsTo(runTestsWorkflowID, prepareSteps...)

	// Install Ruby
	if descriptor.hasRubyVersion {
		configBuilder.AppendStepListItemsTo(runTestsWorkflowID, steps.ScriptStepListItem(rubyInstallScriptStepTitle, rubyInstallScriptStepContent))
	}

	// Restore gem cache
	configBuilder.AppendStepListItemsTo(runTestsWorkflowID, steps.RestoreGemCache())

	// Install dependencies
	if descriptor.hasBundler {
		configBuilder.AppendStepListItemsTo(runTestsWorkflowID, steps.ScriptStepListItem(bundlerInstallScriptStepTitle, bundlerInstallScriptStepContent))
	}

	// Run tests based on detected framework
	testScript := generateTestScript(descriptor)
	if testScript != "" {
		configBuilder.AppendStepListItemsTo(runTestsWorkflowID, steps.ScriptStepListItem("Run tests", testScript))
	}

	// TODO: check if save and restore cache steps are used properly / if they are needed.
	// Save gem cache
	configBuilder.AppendStepListItemsTo(runTestsWorkflowID, steps.SaveGemCache())

	// Deploy steps
	configBuilder.AppendStepListItemsTo(runTestsWorkflowID, steps.DefaultDeployStepList()...)

	config, err := configBuilder.Generate(ScannerName)
	if err != nil {
		return "", err
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func generateTestScript(descriptor configDescriptor) string {
	workdirSetup := ""
	if descriptor.workdir != "" {
		workdirSetup = `pushd "${RUBY_PROJECT_DIR:-.}" > /dev/null

`
	}

	workdirCleanup := ""
	if descriptor.workdir != "" {
		workdirCleanup = `
popd > /dev/null`
	}

	testCommand := ""
	switch descriptor.testFramework {
	case "rspec":
		if descriptor.hasBundler {
			testCommand = "bundle exec rspec"
		} else {
			testCommand = "rspec"
		}
	case "minitest":
		if descriptor.hasRakefile {
			if descriptor.hasBundler {
				testCommand = "bundle exec rake test"
			} else {
				testCommand = "rake test"
			}
		} else {
			if descriptor.hasBundler {
				testCommand = "bundle exec ruby -Itest test/**/*_test.rb"
			} else {
				testCommand = "ruby -Itest test/**/*_test.rb"
			}
		}
	default:
		// Default to rake if Rakefile exists
		if descriptor.hasRakefile {
			if descriptor.hasBundler {
				testCommand = "bundle exec rake test"
			} else {
				testCommand = "rake test"
			}
		} else {
			return ""
		}
	}

	return fmt.Sprintf(`#!/usr/bin/env bash
set -euxo pipefail

%s%s%s`, workdirSetup, testCommand, workdirCleanup)
}
